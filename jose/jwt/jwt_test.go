package jwt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kunitsucom/util.go/jose"
	"github.com/kunitsucom/util.go/jose/jwa"
	"github.com/kunitsucom/util.go/jose/jws"
	"github.com/kunitsucom/util.go/jose/jwt"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		hmacKey := []byte("test")
		token, err := jwt.New(
			jws.WithHMACKey(hmacKey),
			jose.NewHeader(
				jwa.HS256,
				jose.WithType("JWT"),
				jose.WithPrivateHeaderParameter("testPrivateHeaderParameter", "testPrivateHeaderParameter"),
			),
			jwt.NewClaimsSet(
				jwt.WithSubject("userID"),
				jwt.WithAudience("http://localhost/test/aud"),
				jwt.WithExpirationTime(time.Now().Add(1*time.Hour)),
				jwt.WithPrivateClaim("testPrivateClaim", "testPrivateClaim"),
			),
		)
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(jws.UseKey(hmacKey), token); err != nil {
			t.Fatalf("❌: jwt.Verify: err != nil: %v", err)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		hmacKey := []byte("test")
		_, err := jwt.New(
			jws.WithHMACKey(hmacKey),
			jose.NewHeader(
				"failure",
				jose.WithType("JWT"),
				jose.WithPrivateHeaderParameter("testPrivateHeaderParameter", "testPrivateHeaderParameter"),
			),
			jwt.NewClaimsSet(
				jwt.WithSubject("userID"),
				jwt.WithAudience("http://localhost/test/aud"),
				jwt.WithExpirationTime(time.Now().Add(1*time.Hour)),
				jwt.WithPrivateClaim("testPrivateClaim", "testPrivateClaim"),
			),
		)
		if err == nil || !errors.Is(err, jwa.ErrNotImplemented) {
			t.Fatalf("❌: jwt.New: err != jwa.ErrNotImplemented: %v", err)
		}
	})
}
