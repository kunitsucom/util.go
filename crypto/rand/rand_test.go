package randz_test

import (
	"bytes"
	crypto_rand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"testing"

	randz "github.com/kunitsucom/util.go/crypto/rand"
)

var ExampleReader = bytes.NewBufferString("" +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production.",
)

func Example() {
	r := randz.NewReader(randz.WithRandomReader(ExampleReader))
	s, err := r.ReadString(128)
	if err != nil {
		log.Printf("(*randz.Reader).ReadString: %v", err)
		return
	}

	fmt.Printf("very secure random string: %s", s)
	// Output: very secure random string: Wqr1gr1gjg2n12gUnjmn0gox0gn6jvyunugSunj1ng31ngl07y2xv0jwmn1gUnjmn0grwgy0xm3l2rxwuWqr1gr1gjg2n12gUnjmn0gox0gn6jvyunugSunj1ng31ngl
}

func TestCreateCodeVerifier(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		actual, err := randz.ReadString(128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		t.Logf("✅: cv=%s", actual)

		backup := randz.DefaultReader
		t.Cleanup(func() { randz.DefaultReader = backup })
		randz.DefaultReader = bytes.NewBuffer(nil)
		if _, err := randz.ReadString(128); err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const expect = "wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz0123"
		r := randz.NewReader(randz.WithRandomReader(bytes.NewBufferString(
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
		)), randz.WithRandomSource(randz.DefaultRandomSource))
		actual, err := r.ReadString(128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(io.EOF)", func(t *testing.T) {
		t.Parallel()
		r := randz.NewReader(randz.WithRandomReader(bytes.NewReader(nil)))
		_, actual := r.ReadString(128)
		if actual == nil {
			t.Errorf("❌: err == nil")
		}
		expect := io.EOF
		if !errors.Is(actual, expect) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func BenchmarkGenerateRandomString(b *testing.B) {
	b.ResetTimer()

	b.Run("r.Read/mrand.Reader", func(b *testing.B) {
		r := randz.NewReader(randz.WithRandomReader(mrand.New(mrand.NewSource(0))))
		buf := make([]byte, 128)

		for i := 0; i < b.N; i++ {
			_, _ = r.Read(buf)
		}
	})

	b.Run("r.Read/crypto_rand.Reader", func(b *testing.B) {
		r := randz.NewReader(randz.WithRandomReader(crypto_rand.Reader))
		buf := make([]byte, 128)

		for i := 0; i < b.N; i++ {
			_, _ = r.Read(buf)
		}
	})
}
