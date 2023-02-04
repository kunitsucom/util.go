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
//	token, err := jwt.Sign(
//		[]byte("YOUR_HMAC_KEY"),
//		jose.NewHeader(jwa.HS256, jose.WithType("JWT")),
//		jwt.NewClaimsSet(jwt.WithSubject("userID"), jwt.WithExpirationTime(time.Now().Add(1*time.Hour))),
//	)
func Sign(keyOpt jws.SigningKeyOption, header *jose.Header, claimsSet *ClaimsSet) (token string, err error) {
	headerEncoded, err := header.Encode()
	if err != nil {
		return "", fmt.Errorf("(*jose.Header).Encode: %w", err)
	}

	claimsSetEncoded, err := claimsSet.Encode()
	if err != nil {
		return "", fmt.Errorf("(*jwt.ClaimsSet).Encode: %w", err)
	}

	signingInput := headerEncoded + "." + claimsSetEncoded
	signatureEncoded, err := jws.Sign(header.Algorithm, keyOpt, signingInput)
	if err != nil {
		return "", fmt.Errorf("jws.Sign: %w", err)
	}

	return signingInput + "." + signatureEncoded, nil
}
