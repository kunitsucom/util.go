package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type (
	Jitter       func(duration time.Duration) (durationWithJitter time.Duration)
	jitter       struct{ rnd *rand.Rand }
	JitterOption func(j *jitter)
)

func WithJitterRand(rnd *rand.Rand) JitterOption {
	return func(j *jitter) {
		j.rnd = rnd
	}
}

func DefaultJitter(minJitter, maxJitter time.Duration, opts ...JitterOption) Jitter {
	j := &jitter{}

	for _, opt := range opts {
		opt(j)
	}

	return func(duration time.Duration) (durationWithJitter time.Duration) {
		if j.rnd == nil {
			return time.Duration(rand.Int63n(int64(minJitter)+int64(maxJitter)) - int64(minJitter)) //nolint:gosec
		}
		return time.Duration(j.rnd.Int63n(int64(minJitter)+int64(maxJitter)) - int64(minJitter))
	}
}

type Backoff func(initialInterval time.Duration, retries int) (intervalForThisRetry time.Duration)

func DefaultBackoff() Backoff {
	return func(initialInterval time.Duration, retries int) (intervalForThisRetry time.Duration) {
		return time.Duration(int64(initialInterval) << retries)
	}
}

type Config struct {
	initialInterval time.Duration
	maxInterval     time.Duration
	maxRetries      int
	backoff         Backoff
	jitter          Jitter
}

const Infinite = -1

func NewConfig(initialInterval, maxInterval time.Duration, opts ...Option) *Config {
	if initialInterval == 0 {
		initialInterval = 1 * time.Second
	}

	if maxInterval == 0 {
		maxInterval = 30 * time.Second
	}

	c := &Config{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetries:      Infinite,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(c *Config)

func WithBackoff(backoffFunc Backoff) Option {
	return func(c *Config) {
		c.backoff = backoffFunc
	}
}

func WithJitter(jitterFunc Jitter) Option {
	return func(c *Config) {
		c.jitter = jitterFunc
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.maxRetries = maxRetries
	}
}

type Retryer struct {
	config *Config
	// variables
	interval time.Duration
	retries  int
	reason   error
}

func (c *Config) Build() *Retryer {
	copied := *c
	return New(&copied)
}

// New returns *retry.Retryer. *retry.Retryer provides methods to facilitate retry execution.
//
// Is used as follows:
//
//	c := retry.NewConfig(1*time.Second, 10*time.Second, retry.WithMaxRetries(5))
//	r := retry.New(c)
//	for r.Retry() {
//		if r.Retries() == 0 {
//			log.Printf("FIRST! time=%s retries=%d retryAfter=%s", time.Now(), r.Retries(), r.RetryAfter())
//			continue
//		}
//		log.Printf("RETRYING! time=%s retries=%d retryAfter=%s", time.Now(), r.Retries(), r.RetryAfter())
//	}
//	log.Printf("COMPLETE! time=%s retries=%d error=%v", time.Now(), r.Retries(), r.Err())
//
// Then, yields the following results:
//
//	2022/11/16 19:59:17 FIRST! time=2022-11-16 19:59:17.880404 +0900 JST m=+0.002133084 retries=0 retryAfter=1s
//	2022/11/16 19:59:18 RETRYING! time=2022-11-16 19:59:18.881926 +0900 JST m=+1.003646793 retries=1 retryAfter=2s
//	2022/11/16 19:59:20 RETRYING! time=2022-11-16 19:59:20.883498 +0900 JST m=+3.005201043 retries=2 retryAfter=4s
//	2022/11/16 19:59:24 RETRYING! time=2022-11-16 19:59:24.884014 +0900 JST m=+7.005681376 retries=3 retryAfter=8s
//	2022/11/16 19:59:32 RETRYING! time=2022-11-16 19:59:32.885557 +0900 JST m=+15.007154668 retries=4 retryAfter=10s
//	2022/11/16 19:59:42 RETRYING! time=2022-11-16 19:59:42.886954 +0900 JST m=+25.008463293 retries=5 retryAfter=10s
//	2022/11/16 19:59:42 COMPLETE! time=2022-11-16 19:59:42.887136 +0900 JST m=+25.008645501 retries=5 error=retry: reached max retries
//
// If the maximum count of attempts is not given via retry.WithMaxRetries(),
// *retry.Retryer that retry.New() returned will retry infinitely many times.
func New(config *Config) *Retryer {
	copied := *config

	return &Retryer{
		config: &copied,
	}
}

func (r *Retryer) getInitialInterval() time.Duration {
	return r.config.initialInterval
}

func (r *Retryer) truncateAtMaxInterval(d time.Duration) time.Duration {
	if d > r.config.maxInterval {
		return r.config.maxInterval
	}

	return d
}

func (r *Retryer) Retries() (retries int) {
	return r.retries - 1
}

func (r *Retryer) RetryAfter() (retryAfter time.Duration) {
	return r.interval
}

func (r *Retryer) increment() {
	if r.config.backoff == nil {
		r.config.backoff = DefaultBackoff()
	}

	r.interval = r.truncateAtMaxInterval(r.config.backoff(r.getInitialInterval(), r.retries))

	if r.config.jitter != nil {
		r.interval = r.config.jitter(r.interval)
	}

	r.retries++
}

func (r *Retryer) Err() (reason error) {
	return r.reason
}

func (r *Retryer) Retry() bool {
	return r.RetryWithContext(context.Background())
}

var ErrReachedMaxRetries = errors.New("retry: reached max retries")

func (r *Retryer) RetryWithContext(ctx context.Context) bool {
	if 0 <= r.config.maxRetries && r.config.maxRetries <= r.Retries() {
		r.reason = ErrReachedMaxRetries
		return false
	}

	select {
	case <-ctx.Done():
		r.reason = ctx.Err()
		return false
	case <-time.After(r.RetryAfter()):
		r.increment()
		return true
	}
}
