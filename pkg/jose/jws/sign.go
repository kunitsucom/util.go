package jws

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"

	x509z "github.com/kunitsuinc/util.go/pkg/crypto/x509"
	"github.com/kunitsuinc/util.go/pkg/jose/jwa"
)

type SigningKeyOption struct {
	key any
	err error
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

func WithECDSAKeyPEM(pemBytes []byte) SigningKeyOption {
	key, err := x509z.ParseECDSAPrivateKeyPEM(pemBytes)

	return SigningKeyOption{
		key: key,
		err: fmt.Errorf("x509z.ParseECDSAPrivateKeyPEM: %w", err),
	}
}

func Sign(alg string, keyOpt SigningKeyOption, signingInput string) (signatureEncoded string, err error) {
	if keyOpt.err != nil {
		return "", fmt.Errorf("keyOpt.err: %w", keyOpt.err)
	}

	return jwa.JWS(alg).Sign(keyOpt.key, signingInput) //nolint:wrapcheck
}
