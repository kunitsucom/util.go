package retry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type (
	Jitter       func(duration time.Duration) (durationWithJitter time.Duration)
	jitter       struct{ rnd *rand.Rand }
	JitterOption func(j *jitter)
)

func WithDefaultJitterRand(rnd *rand.Rand) JitterOption {
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
			return time.Duration(int64(duration) + rand.Int63n(int64(minJitter)+int64(maxJitter)) - int64(minJitter)) //nolint:gosec
		}
		return time.Duration(int64(duration) + j.rnd.Int63n(int64(minJitter)+int64(maxJitter)) - int64(minJitter))
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
// Is used as follows ( https://go.dev/play/p/nIWIWq-ib6b ):
//
//	c := retry.NewConfig(10*time.Millisecond, 500*time.Millisecond, retry.WithMaxRetries(10))
//	r := retry.New(c)
//
//	for r.Retry() {
//		if r.Retries() == 0 {
//			fmt.Printf("FIRSTTRY! time=%s retries=%d retryAfter=%s\n", time.Now(), r.Retries(), r.RetryAfter())
//			continue
//		}
//		fmt.Printf("RETRYING! time=%s retries=%d retryAfter=%s\n", time.Now(), r.Retries(), r.RetryAfter())
//	}
//	fmt.Printf("COMPLETE! time=%s retries=%d error=%v\n", time.Now(), r.Retries(), r.Err())
//
// Then, yields the following results:
//
//	FIRSTTRY! time=2009-11-10 23:00:00 +0000 UTC m=+0.000000001 retries=0 retryAfter=10ms
//	RETRYING! time=2009-11-10 23:00:00.01 +0000 UTC m=+0.010000001 retries=1 retryAfter=20ms
//	RETRYING! time=2009-11-10 23:00:00.03 +0000 UTC m=+0.030000001 retries=2 retryAfter=40ms
//	RETRYING! time=2009-11-10 23:00:00.07 +0000 UTC m=+0.070000001 retries=3 retryAfter=80ms
//	RETRYING! time=2009-11-10 23:00:00.15 +0000 UTC m=+0.150000001 retries=4 retryAfter=160ms
//	RETRYING! time=2009-11-10 23:00:00.31 +0000 UTC m=+0.310000001 retries=5 retryAfter=320ms
//	RETRYING! time=2009-11-10 23:00:00.63 +0000 UTC m=+0.630000001 retries=6 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:01.13 +0000 UTC m=+1.130000001 retries=7 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:01.63 +0000 UTC m=+1.630000001 retries=8 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:02.13 +0000 UTC m=+2.130000001 retries=9 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:02.63 +0000 UTC m=+2.630000001 retries=10 retryAfter=500ms
//	COMPLETE! time=2009-11-10 23:00:02.63 +0000 UTC m=+2.630000001 retries=10 error=retry: reached max retries
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

var (
	ErrReachedMaxRetries = errors.New("retry: reached max retries")
	ErrUnretryableError  = errors.New("retry: unretryable error")
)

func (r *Retryer) Retry(ctx context.Context) bool {
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

type doConfig struct {
	unretryableErrors     []error
	retryableErrorHandler func(ctx context.Context, r *Retryer, err error)
}

type DoOption func(c *doConfig)

func WithUnretryableErrors(errs []error) DoOption {
	return func(c *doConfig) {
		c.unretryableErrors = errs
	}
}

func WithRetryableErrorHandler(f func(ctx context.Context, r *Retryer, err error)) DoOption {
	return func(c *doConfig) {
		c.retryableErrorHandler = f
	}
}

func (r *Retryer) Do(ctx context.Context, do func(ctx context.Context) error, opts ...DoOption) (err error) {
	c := &doConfig{}

	for _, opt := range opts {
		opt(c)
	}

	for r.Retry(ctx) {
		err = do(ctx)
		if errors.Is(err, nil) {
			return nil
		}
		for _, unretryableError := range c.unretryableErrors {
			if errors.Is(err, unretryableError) {
				return fmt.Errorf("%w: %w", ErrUnretryableError, err)
			}
		}
		if c.retryableErrorHandler != nil {
			c.retryableErrorHandler(ctx, r, err)
		}
	}

	return fmt.Errorf("%w: %w", r.Err(), err)
}
