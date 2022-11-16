package retry_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/retry"
)

func TestRetryer_Retry(t *testing.T) {
	t.Parallel()

	t.Run("success(EXAMPLE)", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 7
			maxInterval = 100 * time.Millisecond
		)
		r := retry.NewConfig(5*time.Millisecond, 100*time.Millisecond, retry.WithMaxRetries(maxRetries)).Build()
		buf := bytes.NewBuffer(nil)
		for r.Retry(context.Background()) {
			fmt.Fprintf(buf, "[EXAMPLE] time=%s retries=%d/%d retryAfter=%s\n", time.Now(), r.Retries(), maxRetries, r.RetryAfter())
		}
		fmt.Fprintf(buf, "[EXAMPLE] time=%s If there is no difference of about %s in execution time between ↑ and ↓, it is OK.\n", time.Now(), maxInterval)
		fmt.Fprintf(buf, "[EXAMPLE] time=%s retries=%d/%d retryAfter=%s\n", time.Now(), r.Retries(), maxRetries, r.RetryAfter())
		t.Logf("\n" + buf.String())
	})

	t.Run("success(constant)", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 20
		r := retry.New(retry.NewConfig(10*time.Microsecond, 10*time.Microsecond, retry.WithMaxRetries(maxRetries)))
		buf := bytes.NewBuffer(nil)
		for r.Retry(context.Background()) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
		}
		const expect = `retries=0/20 retryAfter=10µs; retries=1/20 retryAfter=10µs; retries=2/20 retryAfter=10µs; retries=3/20 retryAfter=10µs; retries=4/20 retryAfter=10µs; retries=5/20 retryAfter=10µs; retries=6/20 retryAfter=10µs; retries=7/20 retryAfter=10µs; retries=8/20 retryAfter=10µs; retries=9/20 retryAfter=10µs; retries=10/20 retryAfter=10µs; retries=11/20 retryAfter=10µs; retries=12/20 retryAfter=10µs; retries=13/20 retryAfter=10µs; retries=14/20 retryAfter=10µs; retries=15/20 retryAfter=10µs; retries=16/20 retryAfter=10µs; retries=17/20 retryAfter=10µs; retries=18/20 retryAfter=10µs; retries=19/20 retryAfter=10µs; retries=20/20 retryAfter=10µs; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual")
		}
		t.Logf("\nactual: " + buf.String())
	})

	t.Run("success(jitter)", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 20
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		for r.Retry(context.Background()) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
		}
		const expect = `retries=0/20 retryAfter=8.165505ms; retries=1/20 retryAfter=7.393152ms; retries=2/20 retryAfter=3.995827ms; retries=3/20 retryAfter=7.197794ms; retries=4/20 retryAfter=4.376202ms; retries=5/20 retryAfter=126.063µs; retries=6/20 retryAfter=4.980153ms; retries=7/20 retryAfter=6.422456ms; retries=8/20 retryAfter=9.894929ms; retries=9/20 retryAfter=2.637646ms; retries=10/20 retryAfter=943.416µs; retries=11/20 retryAfter=6.976708ms; retries=12/20 retryAfter=9.259259ms; retries=13/20 retryAfter=885.298µs; retries=14/20 retryAfter=9.98852ms; retries=15/20 retryAfter=6.116249ms; retries=16/20 retryAfter=3.981575ms; retries=17/20 retryAfter=3.529631ms; retries=18/20 retryAfter=918.339µs; retries=19/20 retryAfter=5.164748ms; retries=20/20 retryAfter=9.441706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual")
		}
		t.Logf("\nactual: " + buf.String())
	})

	t.Run("failure(deadline)", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 1
		jitter := retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond)
		r := retry.New(retry.NewConfig(0, 0, retry.WithMaxRetries(maxRetries), retry.WithBackoff(retry.DefaultBackoff()), retry.WithJitter(jitter)))
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().AddDate(1, 0, 0))
		r.Retry(ctx) // first
		cancel()
		r.Retry(ctx) // second
		r.Retry(ctx) // third
		err := r.Err()
		if err == nil {
			t.Errorf("err == nil")
		}
		const expect = "context canceled"
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("err not contain: `%s` != `%v`", expect, err)
		}
	})
}

func TestRetryer_Do(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
			return nil
		})
		if err != nil {
			t.Errorf("err != nil")
		}
		const expect = `retries=0/20 retryAfter=8.165505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual")
		}
		t.Logf("\nactual: " + buf.String())
	})

	t.Run("failure(reachedMaxRetries)", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			return io.ErrUnexpectedEOF
		}, retry.WithRetryableErrorHandler(func(_ context.Context, r *retry.Retryer, err error) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
		}))
		if err == nil {
			t.Errorf("err == nil")
		}
		const expectErr = "retry: reached max retries: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.165505ms; retries=1/20 retryAfter=7.393152ms; retries=2/20 retryAfter=3.995827ms; retries=3/20 retryAfter=7.197794ms; retries=4/20 retryAfter=4.376202ms; retries=5/20 retryAfter=126.063µs; retries=6/20 retryAfter=4.980153ms; retries=7/20 retryAfter=6.422456ms; retries=8/20 retryAfter=9.894929ms; retries=9/20 retryAfter=2.637646ms; retries=10/20 retryAfter=943.416µs; retries=11/20 retryAfter=6.976708ms; retries=12/20 retryAfter=9.259259ms; retries=13/20 retryAfter=885.298µs; retries=14/20 retryAfter=9.98852ms; retries=15/20 retryAfter=6.116249ms; retries=16/20 retryAfter=3.981575ms; retries=17/20 retryAfter=3.529631ms; retries=18/20 retryAfter=918.339µs; retries=19/20 retryAfter=5.164748ms; retries=20/20 retryAfter=9.441706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual")
		}
		t.Logf("\nactual: " + buf.String())
	})

	t.Run("failure(unretryableErrors)", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
			return io.ErrUnexpectedEOF
		}, retry.WithUnretryableErrors([]error{io.ErrUnexpectedEOF}))
		if err == nil {
			t.Errorf("err == nil")
		}
		const expectErr = "retry: unretryable error: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.165505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual")
		}
		t.Logf("\nactual: " + buf.String())
	})
}
