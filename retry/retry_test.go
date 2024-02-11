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

	"github.com/kunitsucom/util.go/retry"
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
		t.Logf("✅: %s", buf)
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
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: " + buf.String())
	})

	t.Run("success(jitter)", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 20
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithDefaultJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		for r.Retry(context.Background()) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
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
			t.Errorf("❌: err == nil")
		}
		const expect = "context canceled"
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: err not contain: `%s` != `%v`", expect, err)
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
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithDefaultJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
			return nil
		})
		if err != nil {
			t.Errorf("❌: err != nil")
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure(reachedMaxRetries)", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithDefaultJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			return io.ErrUnexpectedEOF
		}, retry.WithRetryableErrorHandler(func(_ context.Context, r *retry.Retryer, err error) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
		}))
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		const expectErr = "retry: reached max retries: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("❌: err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure(unretryableErrors)", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		r := retry.New(retry.NewConfig(1*time.Microsecond, 999*time.Microsecond, retry.WithMaxRetries(maxRetries), retry.WithJitter(retry.DefaultJitter(1*time.Millisecond, 10*time.Millisecond, retry.WithDefaultJitterRand(rand.New(rand.NewSource(0))))))) //nolint:gosec
		buf := bytes.NewBuffer(nil)
		err := r.Do(context.Background(), func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), maxRetries, r.RetryAfter())
			return io.ErrUnexpectedEOF
		}, retry.WithUnretryableErrors([]error{io.ErrUnexpectedEOF}))
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		const expectErr = "retry: unretryable error: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("❌: err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})
}
