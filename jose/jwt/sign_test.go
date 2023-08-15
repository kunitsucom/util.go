package jwt_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kunitsucom/util.go/jose"
	"github.com/kunitsucom/util.go/jose/jwa"
	"github.com/kunitsucom/util.go/jose/jws"
	"github.com/kunitsucom/util.go/jose/jwt"
)

func TestSign(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		hmacKey := []byte("test")
		signingInput, signatureEncoded, err := jwt.Sign(
			jws.WithHMACKey(hmacKey),
			jose.NewHeader(
				jwa.HS256,
				jose.WithType("JWT"),
				jose.WithPrivateHeaderParameter("testPrivateHeaderParameter", "testPrivateHeaderParameter"),
			),
			jwt.NewClaimsSet(
				jwt.WithSubject("userID"),
				jwt.WithAudience("http://localhost/test/aud"),
				jwt.WithIssuer("http://localhost/test/iss"),
				jwt.WithExpirationTime(time.Now().Add(1*time.Hour)),
				jwt.WithPrivateClaim("testPrivateClaim", "testPrivateClaim"),
			),
		)
		if err != nil {
			t.Errorf("❌: jwt.Sign: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(
			jws.UseKey(hmacKey),
			signingInput+"."+signatureEncoded,
			jwt.VerifyAudience("http://localhost/test/aud"),
			jwt.VerifyIssuer("http://localhost/test/iss"),
			jwt.VerifyPrivateClaims(func(pc jwt.PrivateClaims) error {
				if _, ok := pc["testPrivateClaim"]; !ok {
					return errors.New("testPrivateClaim")
				}
				return nil
			}),
		); err != nil {
			t.Errorf("❌: jwt.Verify: err != nil: %v", err)
		}
	})

	t.Run("failure(header.Encode)", func(t *testing.T) {
		t.Parallel()
		hmacKey := []byte("test")
		_, _, err := jwt.Sign(
			jws.WithHMACKey(hmacKey),
			jose.NewHeader(
				jwa.HS256,
				jose.WithType("JWT"),
				jose.WithPrivateHeaderParameter("testPrivateHeaderParameter", func() {}),
			),
			jwt.NewClaimsSet(
				jwt.WithSubject("userID"),
				jwt.WithAudience("http://localhost/test/aud"),
				jwt.WithExpirationTime(time.Now().Add(1*time.Hour)),
				jwt.WithPrivateClaim("testPrivateClaim", "testPrivateClaim"),
			),
		)
		if expect := "invalid private header parameters"; err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: jwt.Sign: err != %s: %v", expect, err)
		}
	})

	t.Run("failure(claimsSet.Encode)", func(t *testing.T) {
		t.Parallel()
		hmacKey := []byte("test")
		_, _, err := jwt.Sign(
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
				jwt.WithPrivateClaim("testPrivateClaim", func() {}),
			),
		)
		if expect := "invalid private claims"; err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: jwt.Sign: err != %s: %v", expect, err)
		}
	})

	t.Run("failure(jwa.ErrInvalidKeyReceived)", func(t *testing.T) {
		t.Parallel()
		_, _, err := jwt.Sign(
			jws.WithHMACKey(nil),
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
		if err == nil || !errors.Is(err, jwa.ErrInvalidKeyReceived) {
			t.Errorf("❌: jwt.Sign: err != jwa.ErrInvalidKeyReceived: %v", err)
		}
	})
}
