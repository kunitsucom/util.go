//nolint:paralleltest
package env_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/env"
)

//nolint:revive,stylecheck
const TEST_ENV_KEY = "TEST_ENV_KEY"

func TestString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual, err := env.String(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = ""
		actual, err := env.String(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestStringOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual := env.StringOrDefault(TEST_ENV_KEY, "default")
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = "default"
		actual := env.StringOrDefault(TEST_ENV_KEY, "default")
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual := env.MustString(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustString(TEST_ENV_KEY)
	})
}

func TestBool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual, err := env.Bool(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = false
		actual, err := env.Bool(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = false
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Bool(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestBoolOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual := env.BoolOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = true
		actual := env.BoolOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustBool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual := env.MustBool(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustBool(TEST_ENV_KEY)
	})
}

func TestInt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual, err := env.Int(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := env.Int(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Int(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestIntOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(expect, 10))
		actual := env.IntOrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 2000000000
		actual := env.IntOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustInt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := env.MustInt(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustInt(TEST_ENV_KEY)
	})
}

func TestInt64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual, err := env.Int64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := env.Int64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Int64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestInt64OrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(expect, 10))
		actual := env.Int64OrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 100000000000
		actual := env.Int64OrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustInt64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := env.MustInt64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustInt64(TEST_ENV_KEY)
	})
}

func TestUint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual, err := env.Uint(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := env.Uint(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Uint(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestUintOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := env.UintOrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 100000000000
		actual := env.UintOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustUint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := env.MustUint(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustUint(TEST_ENV_KEY)
	})
}

func TestUint64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual, err := env.Uint64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := env.Uint64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Uint64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestUint64OrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := env.Uint64OrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 100000000000
		actual := env.Uint64OrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustUint64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := env.MustUint64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustUint64(TEST_ENV_KEY)
	})
}

func TestFloat64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000.1
		t.Setenv(TEST_ENV_KEY, strconv.FormatFloat(expect, 'f', -1, 64))
		actual, err := env.Float64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := env.Float64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-atoi)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Float64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFloat64OrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000.1
		t.Setenv(TEST_ENV_KEY, strconv.FormatFloat(expect, 'f', -1, 64))
		actual := env.Float64OrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 100000000000.1
		actual := env.Float64OrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustFloat64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000.1
		t.Setenv(TEST_ENV_KEY, strconv.FormatFloat(expect, 'f', -1, 64))
		actual := env.MustFloat64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustFloat64(TEST_ENV_KEY)
	})
}

func TestSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual, err := env.Second(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: env.Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := env.Second(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: env.Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSecondOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual := env.SecondOrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 30 * time.Second
		actual := env.SecondOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual := env.MustSecond(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = env.MustSecond(TEST_ENV_KEY)
	})
}
