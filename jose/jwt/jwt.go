package jwt

import (
	"fmt"

	"github.com/kunitsucom/util.go/jose"
	"github.com/kunitsucom/util.go/jose/jws"
)

// New
//
// Example:
//
//	token, err := jwt.New(
//		jws.WithHMACKey([]byte("YOUR_HMAC_KEY"),
//		jose.NewHeader(jwa.HS256, jose.WithType("JWT")),
//		jwt.NewClaimsSet(jwt.WithSubject("userID"), jwt.WithExpirationTime(time.Now().Add(1*time.Hour))),
//	)
func New(keyOpt jws.SigningKeyOption, header *jose.Header, claimsSet *ClaimsSet) (token string, err error) {
	signingInput, signatureEncoded, err := Sign(keyOpt, header, claimsSet)
	if err != nil {
		return "", fmt.Errorf("jwt.Sign: %w", err)
	}

	return signingInput + "." + signatureEncoded, nil
}
