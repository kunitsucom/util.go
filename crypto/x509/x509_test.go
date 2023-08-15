package x509z_test

import (
	"errors"
	"strings"
	"testing"

	x509z "github.com/kunitsucom/util.go/crypto/x509"
	testz "github.com/kunitsucom/util.go/test"
)

func TestParseRSAPrivateKeyPEM(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKey2048BitPEM))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure(ErrInvalidPEMFormat)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPrivateKeyPEM([]byte("Invalid"))
		if !errors.Is(err, x509z.ErrInvalidPEMFormat) {
			t.Errorf("❌: err != x509z.ErrInvalidPEMFormat: %v", err)
		}
	})

	t.Run("failure(x509.ParsePKCS8PrivateKey)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKeyInvalidPEM))
		if err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
		const expect = "asn1: syntax error: data truncated"
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: err != expect(%s): %v", expect, err)
		}
	})

	t.Run("failure(x509z.ErrPublicKeyTypeMismatch)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPrivateKeyPEM([]byte(testz.TestECDSAPrivateKey256BitPEM))
		if !errors.Is(err, x509z.ErrKeyTypeMismatch) {
			t.Errorf("❌: err != x509z.ErrPublicKeyTypeMismatch: %v", err)
		}
	})
}

func TestParseRSAPublicKeyPEM(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPublicKeyPEM([]byte(testz.TestRSAPublicKey2048BitPEM))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure(ErrInvalidPEMFormat)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPublicKeyPEM([]byte("Invalid"))
		if !errors.Is(err, x509z.ErrInvalidPEMFormat) {
			t.Errorf("❌: err != x509z.ErrInvalidPEMFormat: %v", err)
		}
	})

	t.Run("failure(x509.ParsePKIXPublicKey)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPublicKeyPEM([]byte(testz.TestRSAPrivateKeyInvalidPEM))
		if err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
		const expect = "asn1: syntax error: data truncated"
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: err != expect(%s): %v", expect, err)
		}
	})

	t.Run("failure(x509z.ErrPublicKeyTypeMismatch)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseRSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))
		if !errors.Is(err, x509z.ErrKeyTypeMismatch) {
			t.Errorf("❌: err != x509z.ErrPublicKeyTypeMismatch: %v", err)
		}
	})
}

func TestParseECDSAPrivateKeyPEM(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPrivateKeyPEM([]byte(testz.TestECDSAPrivateKey256BitPEM))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure(ErrInvalidPEMFormat)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPrivateKeyPEM([]byte("Invalid"))
		if !errors.Is(err, x509z.ErrInvalidPEMFormat) {
			t.Errorf("❌: err != x509z.ErrInvalidPEMFormat: %v", err)
		}
	})

	t.Run("failure(x509.ParsePKCS8PrivateKey)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKeyInvalidPEM))
		if err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
		const expect = "asn1: syntax error: data truncated"
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: err != expect(%s): %v", expect, err)
		}
	})

	t.Run("failure(x509z.ErrPublicKeyTypeMismatch)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKey2048BitPEM))
		if !errors.Is(err, x509z.ErrKeyTypeMismatch) {
			t.Errorf("❌: err != x509z.ErrPublicKeyTypeMismatch: %v", err)
		}
	})
}

func TestParseECDSAPublicKeyPEM(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure(ErrInvalidPEMFormat)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPublicKeyPEM([]byte("Invalid"))
		if !errors.Is(err, x509z.ErrInvalidPEMFormat) {
			t.Errorf("❌: err != x509z.ErrInvalidPEMFormat: %v", err)
		}
	})

	t.Run("failure(x509.ParsePKIXPublicKey)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestRSAPrivateKeyInvalidPEM))
		if err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
		const expect = "asn1: syntax error: data truncated"
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: err != expect(%s): %v", expect, err)
		}
	})

	t.Run("failure(x509z.ErrPublicKeyTypeMismatch)", func(t *testing.T) {
		t.Parallel()
		_, err := x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestRSAPublicKey2048BitPEM))
		if !errors.Is(err, x509z.ErrKeyTypeMismatch) {
			t.Errorf("❌: err != x509z.ErrPublicKeyTypeMismatch: %v", err)
		}
	})
}
