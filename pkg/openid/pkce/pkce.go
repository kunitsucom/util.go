package pkce

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	randz "github.com/kunitsuinc/util.go/pkg/crypto/rand"
)

type (
	CodeVerifier        string
	CodeChallengeMethod string
)

const (
	// CodeChallengeMethodPlainShouldNotBeUsed
	//
	// Clients are permitted to use "plain" only if they cannot support
	// "S256" for some technical reason and know via out-of-band
	// configuration that the server supports "plain".
	//
	// ref. https://www.rfc-editor.org/rfc/rfc7636#section-4.2
	CodeChallengeMethodPlainShouldNotBeUsed CodeChallengeMethod = "plain"
	// CodeChallengeMethodS256
	//
	// If the client is capable of using "S256", it MUST use "S256", as
	// "S256" is Mandatory To Implement (MTI) on the server.
	//
	//	code_challenge = BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
	//
	// ref. https://www.rfc-editor.org/rfc/rfc7636#section-4.2
	CodeChallengeMethodS256 CodeChallengeMethod = "S256"
)

// UnreservedCharacters
//
//	unreserved  = ALPHA / DIGIT / "-" / "." / "_" / "~".
//
// ref. https://www.rfc-editor.org/rfc/rfc3986#section-2.3
const UnreservedCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"

var ErrCodeVerifierLength = errors.New("pkce: code verifier length must be between 43 and 128")

func CreateCodeVerifier(length int) (CodeVerifier, error) {
	return createCodeVerifier(rand.Reader, length)
}

func createCodeVerifier(randReader io.Reader, length int) (CodeVerifier, error) {
	// https://www.rfc-editor.org/rfc/rfc7636#section-4.1
	// code-verifier = 43*128unreserved
	if length < 43 || 128 < length {
		return "", ErrCodeVerifierLength
	}

	r := randz.NewReader(randz.WithRandomSource(UnreservedCharacters), randz.WithRandomReader(randReader))
	random, err := r.ReadString(length)
	if err != nil {
		return "", fmt.Errorf("❌: randz.GenerateRandomString: %w", err)
	}

	return CodeVerifier(random), nil
}

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
