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
	ErrInvalidPEMFormat = errors.New("crypto/x509z: invalid pem format")
	ErrKeyTypeMismatch  = errors.New("crypto/x509z: key type mismatch")
)

func ParseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKCS8PrivateKey: %w", err)
	}

	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expect=%T actual=%T: %w", rsaPrivateKey, key, ErrKeyTypeMismatch)
	}

	return rsaPrivateKey, nil
}

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
		return nil, fmt.Errorf("expect=%T actual=%T: %w", rsaPublicKey, key, ErrKeyTypeMismatch)
	}

	return rsaPublicKey, nil
}

func ParseECDSAPrivateKeyPEM(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKCS8PrivateKey: %w", err)
	}

	ecdsaPrivateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expect=%T actual=%T: %w", ecdsaPrivateKey, key, ErrKeyTypeMismatch)
	}

	return ecdsaPrivateKey, nil
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
		return nil, fmt.Errorf("expect=%T actual=%T: %w", ecdsaPublicKey, key, ErrKeyTypeMismatch)
	}

	return ecdsaPublicKey, nil
}
