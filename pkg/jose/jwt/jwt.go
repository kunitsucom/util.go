package jwt

import (
	"github.com/kunitsuinc/util.go/pkg/jose"
	"github.com/kunitsuinc/util.go/pkg/jose/jws"
)

// New is an alias of jwt.Sign.
//
// Example:
//
//	token, err := jwt.New(
//		[]byte("YOUR_HMAC_KEY"),
//		jose.NewHeader(jwa.HS256, jose.WithType("JWT")),
//		jwt.NewClaimsSet(jwt.WithSubject("userID"), jwt.WithExpirationTime(time.Now().Add(1*time.Hour))),
//	)
func New(keyOpt jws.SigningKeyOption, header *jose.Header, claimsSet *ClaimsSet) (token string, err error) {
	return Sign(keyOpt, header, claimsSet)
}
