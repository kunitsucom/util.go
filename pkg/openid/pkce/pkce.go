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
	// CodeChallengeMethodPlainShouldNotBeUsed
	//
	// **WARNING**: plain SHOULD NOT BE USED FOR code_challenge_method.
	//
	// ref. https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
	CodeChallengeMethodPlainShouldNotBeUsed CodeChallengeMethod = "plain"
	// CodeChallengeMethodS256
	//
	//	code_challenge = BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
	//
	// ref. https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
	CodeChallengeMethodS256 CodeChallengeMethod = "S256"
)

func (c CodeChallengeMethod) String() string { return string(c) }

func (c CodeVerifier) Encode(method CodeChallengeMethod) string {
	switch method {
	case CodeChallengeMethodS256:
		a := sha256.Sum256([]byte(c))
		return base64.RawURLEncoding.EncodeToString(a[:])
	case CodeChallengeMethodPlainShouldNotBeUsed:
		return string(c)
	default:
		return string(c)
	}
}
