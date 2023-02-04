package jws

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/kunitsuinc/util.go/pkg/jose/jwa"
)

type SigningKeyOption struct {
	key any
}

func WithKey(key any) SigningKeyOption {
	return SigningKeyOption{
		key: key,
	}
}

func WithHMACKey(key []byte) SigningKeyOption {
	return WithKey(key)
}

func WithRSAKey(key *rsa.PrivateKey) SigningKeyOption {
	return WithKey(key)
}

func WithECDSAKey(key *ecdsa.PrivateKey) SigningKeyOption {
	return WithKey(key)
}

func Sign(alg string, keyOpt SigningKeyOption, signingInput string) (signatureEncoded string, err error) {
	return jwa.JWS(alg).Sign(keyOpt.key, signingInput) //nolint:wrapcheck
}
