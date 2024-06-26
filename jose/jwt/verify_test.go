package jwt_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/jose"
	"github.com/kunitsucom/util.go/jose/jwa"
	"github.com/kunitsucom/util.go/jose/jws"
	"github.com/kunitsucom/util.go/jose/jwt"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	testHS256Key := []byte(`your-256-bit-secret`)
	testHS256JWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), testHS256JWT); err != nil {
			t.Errorf("❌: jwt.Verify: err != nil: %v", err)
		}
	})

	t.Run("failure(jws.ErrInvalidTokenReceived)", func(t *testing.T) {
		t.Parallel()
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), "invalid.jwt"); err == nil || !errors.Is(err, jws.ErrInvalidTokenReceived) {
			t.Errorf("❌: jwt.Verify: err != jws.ErrInvalidTokenReceived: %v", err)
		}
	})

	t.Run("failure(cs.Decode)", func(t *testing.T) {
		t.Parallel()
		expect := "invalid character"
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.jwt"); err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: jwt.Verify: err != %s: %v", expect, err)
		}
	})

	t.Run("failure(exp,jwt.ErrTokenIsExpired)", func(t *testing.T) {
		t.Parallel()
		signingInput, signatureEncoded, err := jwt.Sign(jws.WithHMACKey(testHS256Key), jose.NewHeader(jwa.HS256), jwt.NewClaimsSet(jwt.WithExpirationTime(time.Now().Add(-1*time.Hour))))
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), signingInput+"."+signatureEncoded); err == nil || !errors.Is(err, jwt.ErrTokenIsExpired) {
			t.Errorf("❌: jwt.Verify: err != jwt.ErrTokenIsExpired: %v", err)
		}
	})

	t.Run("failure(nbf,jwt.ErrTokenIsExpired)", func(t *testing.T) {
		t.Parallel()
		signingInput, signatureEncoded, err := jwt.Sign(jws.WithHMACKey(testHS256Key), jose.NewHeader(jwa.HS256), jwt.NewClaimsSet(jwt.WithNotBefore(time.Now().Add(1*time.Hour))))
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), signingInput+"."+signatureEncoded); err == nil || !errors.Is(err, jwt.ErrTokenIsExpired) {
			t.Errorf("❌: jwt.Verify: err != jwt.ErrTokenIsExpired: %v", err)
		}
	})

	t.Run("failure(aud,jwt.ErrAudienceIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		signingInput, signatureEncoded, err := jwt.Sign(jws.WithHMACKey(testHS256Key), jose.NewHeader(jwa.HS256), jwt.NewClaimsSet(jwt.WithAudience("aud")))
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), signingInput+"."+signatureEncoded, jwt.VerifyAudience("notMatch")); err == nil || !errors.Is(err, jwt.ErrAudienceIsNotMatch) {
			t.Errorf("❌: jwt.Verify: err != jwt.ErrAudienceIsNotMatch: %v", err)
		}
	})

	t.Run("failure(aud,jwt.ErrIssuerIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		signingInput, signatureEncoded, err := jwt.Sign(jws.WithHMACKey(testHS256Key), jose.NewHeader(jwa.HS256), jwt.NewClaimsSet(jwt.WithIssuer("iss")))
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), signingInput+"."+signatureEncoded, jwt.VerifyIssuer("notMatch")); err == nil || !errors.Is(err, jwt.ErrIssuerIsNotMatch) {
			t.Errorf("❌: jwt.Verify: err != jwt.ErrIssuerIsNotMatch: %v", err)
		}
	})

	t.Run("failure(aud,jwt.ErrAudienceIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		signingInput, signatureEncoded, err := jwt.Sign(jws.WithHMACKey(testHS256Key), jose.NewHeader(jwa.HS256), jwt.NewClaimsSet(jwt.WithPrivateClaim("privateClaim", "privateClaim")))
		if err != nil {
			t.Fatalf("❌: jwt.New: err != nil: %v", err)
		}
		expect := "test private claim"
		if _, _, err := jwt.Verify(jws.UseKey(testHS256Key), signingInput+"."+signatureEncoded, jwt.VerifyPrivateClaims(func(privateClaims jwt.PrivateClaims) error {
			_, ok := privateClaims["privateClaimDoesNotExist"]
			if !ok {
				return errorz.Errorf("%s: %w", expect, testingz.ErrTestError)
			}
			return nil
		})); err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: jwt.Verify: err != %s: %v", expect, err)
		}
	})

	t.Run("failure(aud,jwa.ErrFailedToVerifySignature)", func(t *testing.T) {
		t.Parallel()
		if _, _, err := jwt.Verify(jws.UseKey([]byte(`your-256-bit-secret`)), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.invalidSignature"); err == nil || !errors.Is(err, jwa.ErrFailedToVerifySignature) {
			t.Errorf("❌: jwt.Verify: err != jwa.ErrFailedToVerifySignature: %v", err)
		}
	})
}
