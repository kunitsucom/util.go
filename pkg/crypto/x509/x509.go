package x509z

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrInvalidPEMFormat      = errors.New("crypto/x509z: invalid pem format")
	ErrPublicKeyTypeMismatch = errors.New("crypto/x509z: public key type mismatch")
)

func ParseRSAPublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKIXPublicKey: %w", err)
	}

	rsaPublicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expect=%T actual=%T: %w", rsaPublicKey, key, ErrPublicKeyTypeMismatch)
	}

	return rsaPublicKey, nil
}

func ParseECDSAPublicKeyPEM(pemBytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKIXPublicKey: %w", err)
	}

	ecdsaPublicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expect=%T actual=%T: %w", ecdsaPublicKey, key, ErrPublicKeyTypeMismatch)
	}

	return ecdsaPublicKey, nil
}
