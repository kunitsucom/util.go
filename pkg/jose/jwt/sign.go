package jwt

import (
	"fmt"

	"github.com/kunitsuinc/util.go/pkg/jose"
	"github.com/kunitsuinc/util.go/pkg/jose/jws"
)

// Sign
//
// Example:
//
//	signingInput, signatureEncoded, err := jwt.Sign(
//		jws.WithHMACKey([]byte("YOUR_HMAC_KEY"),
//		jose.NewHeader(jwa.HS256, jose.WithType("JWT")),
//		jwt.NewClaimsSet(jwt.WithSubject("userID"), jwt.WithExpirationTime(time.Now().Add(1*time.Hour))),
//	)
func Sign(keyOpt jws.SigningKeyOption, header *jose.Header, claimsSet *ClaimsSet) (signingInput, signatureEncoded string, err error) {
	headerEncoded, err := header.Encode()
	if err != nil {
		return "", "", fmt.Errorf("(*jose.Header).Encode: %w", err)
	}

	claimsSetEncoded, err := claimsSet.Encode()
	if err != nil {
		return "", "", fmt.Errorf("(*jwt.ClaimsSet).Encode: %w", err)
	}

	input := headerEncoded + "." + claimsSetEncoded
	sig, err := jws.Sign(header.Algorithm, keyOpt, input)
	if err != nil {
		return "", "", fmt.Errorf("jws.Sign: %w", err)
	}

	return input, sig, nil
}
