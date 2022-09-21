// nolint: paralleltest
package util_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go"
)

// nolint: revive,stylecheck
const TEST_ENV_KEY = "TEST_ENV_KEY"

func TestEnv(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual, err := util.Env(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("util.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = ""
		actual, err := util.Env(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustEnv(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual := util.MustEnv(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()
		_ = util.MustEnv(TEST_ENV_KEY)
	})
}

func TestEnvInt64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual, err := util.EnvInt64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("util.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := util.EnvInt64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := util.EnvInt64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustEnvInt64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := util.MustEnvInt64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()
		_ = util.MustEnvInt64(TEST_ENV_KEY)
	})
}

func TestEnvFloat64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000.1
		t.Setenv(TEST_ENV_KEY, strconv.FormatFloat(expect, 'f', -1, 64))
		actual, err := util.EnvFloat64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("util.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := util.EnvFloat64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := util.EnvFloat64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustEnvFloat64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000.1
		t.Setenv(TEST_ENV_KEY, strconv.FormatFloat(expect, 'f', -1, 64))
		actual := util.MustEnvFloat64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()
		_ = util.MustEnvFloat64(TEST_ENV_KEY)
	})
}

func TestEnvSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual, err := util.EnvSecond(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("util.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := util.EnvSecond(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("util.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustEnvSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual := util.MustEnvSecond(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()
		_ = util.MustEnvSecond(TEST_ENV_KEY)
	})
}
