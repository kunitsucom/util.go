package pkce

import (
	"crypto/sha256"
	"encoding/base64"
)

type (
	CodeVerifier        string
	CodeChallengeMethod string
)

const (
	CodeChallengeMethodPlain CodeChallengeMethod = "plain"
	CodeChallengeMethodS256  CodeChallengeMethod = "S256"
)

func (c CodeChallengeMethod) String() string { return string(c) }

func (c CodeVerifier) Encode(method CodeChallengeMethod) string {
	switch method {
	case CodeChallengeMethodS256:
		a := sha256.Sum256([]byte(c))
		return base64.RawURLEncoding.EncodeToString(a[:])
	case CodeChallengeMethodPlain:
		return string(c)
	default:
		return string(c)
	}
}
