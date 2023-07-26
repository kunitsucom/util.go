package cipherz //nolint:testpackage

import (
	"bytes"
	"errors"
	"io"
	"testing"

	errorz "github.com/kunitsuinc/util.go/pkg/errors"
	testz "github.com/kunitsuinc/util.go/pkg/test"
)

var (
	testKey       = []byte("0123456789abcdef0123456789abcdef")
	testPlainText = []byte("testPlainText")
	testNonce     = []byte("0123456789ab")
)

func TestNewAESGCM(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		c, err := NewAESGCM(testKey)
		if err != nil {
			t.Fatalf("❌: NewAESGCM: %v", err)
		}
		if c == nil {
			t.Fatal("❌: NewAESGCM: c is nil")
		}
		cipherBytes, err := c.Encrypt(testPlainText)
		if err != nil {
			t.Fatalf("❌: c.Encrypt: %v", err)
		}
		plainBytes, err := c.Decrypt(cipherBytes)
		if err != nil {
			t.Fatalf("❌: c.Decrypt: %v", err)
		}
		if expect, actual := string(testPlainText), string(plainBytes); expect != actual {
			t.Fatalf("❌: c.Decrypt: expect(%s) != actual(%s)", expect, actual)
		}
	})
	t.Run("failure_invalid_key_size", func(t *testing.T) {
		t.Parallel()
		c, err := NewAESGCM([]byte(""))
		if expect := "aes.NewCipher: crypto/aes: invalid key size 0"; !errorz.Contains(err, expect) {
			t.Fatalf("❌: newAESGCM: expect(%s) != actual(%v)", expect, err)
		}
		if c != nil {
			t.Fatalf("❌: newAESGCM: c should be nil: %v", c)
		}
	})
	t.Run("failure,ErrCipherBytesIsTooShort", func(t *testing.T) {
		t.Parallel()
		c, err := NewAESGCM(testKey)
		if err != nil {
			t.Fatalf("❌: NewAESGCM: %v", err)
		}
		if c == nil {
			t.Fatal("❌: NewAESGCM: c is nil")
		}
		p, err := c.Decrypt([]byte(""))
		if expect := ErrCipherBytesIsTooShort; !errors.Is(err, expect) {
			t.Fatalf("❌: c.Decrypt: expect(%v) != actual(%v)", expect, err)
		}
		if p != nil {
			t.Fatalf("❌: c.Decrypt: p should be nil: %v", p)
		}
	})
	t.Run("failure,ErrCipherBytesIsTooShort", func(t *testing.T) {
		t.Parallel()
		c, err := NewAESGCM(testKey)
		if err != nil {
			t.Fatalf("❌: NewAESGCM: %v", err)
		}
		if c == nil {
			t.Fatal("❌: NewAESGCM: c is nil")
		}
		p, err := c.Decrypt(append(testNonce, 'a'))
		if expect := "gcm.aead.Open: cipher: message authentication failed"; !errorz.Contains(err, expect) {
			t.Fatalf("❌: c.Decrypt: expect(%v) != actual(%v)", expect, err)
		}
		if p != nil {
			t.Fatalf("❌: c.Decrypt: p should be nil: %v", p)
		}
	})
}

func Test_newAESGCM(t *testing.T) {
	t.Parallel()
	t.Run("failure_nonce_length", func(t *testing.T) {
		t.Parallel()
		c, err := newAESGCM(testKey, 0)
		if expect := "cipher: the nonce can't have zero length, or the security of the key will be immediately compromised"; !errorz.Contains(err, expect) {
			t.Fatalf("❌: newAESGCM: expect(%s) != actual(%v)", expect, err)
		}
		if c != nil {
			t.Fatalf("❌: newAESGCM: c should be nil: %v", c)
		}
	})

	t.Run("failure_invalid_key_size", func(t *testing.T) {
		t.Parallel()
		c, err := newAESGCM(testKey, recommendedGCMNonceSize)
		if err != nil {
			t.Fatalf("❌: NewAESGCM: %v", err)
		}
		if c == nil {
			t.Fatal("❌: NewAESGCM: c is nil")
		}
		c.rand = testz.NewReadWriter(bytes.NewBuffer(testNonce), 0, io.ErrUnexpectedEOF)
		if _, err := c.Encrypt(testPlainText); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("❌: c.Encrypt: expect(%v) != actual(%v)", io.ErrUnexpectedEOF, err)
		}
	})
}
