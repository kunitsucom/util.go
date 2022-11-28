package jws

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"strings"

	"github.com/kunitsuinc/util.go/pkg/jose"
	"github.com/kunitsuinc/util.go/pkg/jose/jwa"
)

var (
	ErrInvalidTokenReceived        = errors.New(`jws: invalid token received, token must have 3 parts`)
	ErrAlgorithmNoneIsNotSupported = errors.New(`jws: algorithm "none" is not supported`)
	ErrInvalidKeyReceived          = errors.New(`jws: invalid key received`)
	ErrFailedToVerifySignature     = errors.New(`jws: failed to verify signature`)
	ErrInvalidAlgorithm            = errors.New(`jws: invalid algorithm`)
)

// memo: https://cs.opensource.google/go/x/oauth2/+/refs/tags/v0.2.0:jws/jws.go;l=167
func Verify(token string, key crypto.PublicKey) error { //nolint:funlen,cyclop
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrInvalidTokenReceived
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: header: %w", err)
	}

	header := new(jose.Header)
	if err := json.Unmarshal(headerJSON, header); err != nil {
		return fmt.Errorf("json.Unmarshal: header: %w", err)
	}

	jwsSigningInput := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: signature: %w", err)
	}

	switch header.Algorithm {
	case jwa.HS256.String():
		if err := verifyHS(key, sha256.New, jwsSigningInput, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.HS384.String():
		if err := verifyHS(key, sha512.New384, jwsSigningInput, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.HS512.String():
		if err := verifyHS(key, sha512.New, jwsSigningInput, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.RS256.String():
		if err := verifyRS(key, sha256.New, jwsSigningInput, crypto.SHA256, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.RS384.String():
		if err := verifyRS(key, sha512.New384, jwsSigningInput, crypto.SHA384, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.RS512.String():
		if err := verifyRS(key, sha512.New, jwsSigningInput, crypto.SHA512, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.ES256.String():
		if err := verifyES(key, crypto.SHA256, jwsSigningInput, 32, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.ES384.String():
		if err := verifyES(key, crypto.SHA384, jwsSigningInput, 48, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.ES512.String():
		if err := verifyES(key, crypto.SHA512, jwsSigningInput, 66, signature); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.PS256.String():
		if err := verifyPS(key, sha256.New, jwsSigningInput, crypto.SHA256, signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.PS384.String():
		if err := verifyPS(key, sha512.New384, jwsSigningInput, crypto.SHA384, signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.PS512.String():
		if err := verifyPS(key, sha512.New, jwsSigningInput, crypto.SHA512, signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
		return nil
	case jwa.None.String():
		return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, ErrAlgorithmNoneIsNotSupported)
	}

	return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, ErrInvalidAlgorithm)
}

func verifyHS(key crypto.PublicKey, hashNewFunc func() hash.Hash, jwsSigningInput string, signature []byte) error {
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hmac.New(hashNewFunc, keyBytes)
	h.Write([]byte(jwsSigningInput))
	summed := h.Sum(nil)
	if !bytes.Equal(signature, summed) {
		return ErrFailedToVerifySignature
	}
	return nil
}

func verifyRS(key crypto.PublicKey, hashNewFunc func() hash.Hash, jwsSigningInput string, cryptoHash crypto.Hash, signature []byte) error {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hashNewFunc()
	h.Write([]byte(jwsSigningInput))
	if err := rsa.VerifyPKCS1v15(pub, cryptoHash, h.Sum(nil), signature); err != nil {
		return fmt.Errorf("rsa.VerifyPKCS1v15: %w", err)
	}
	return nil
}

func verifyES(key crypto.PublicKey, cryptoHash crypto.Hash, jwsSigningInput string, keySize int, signature []byte) error {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	if len(signature) != 2*keySize {
		return fmt.Errorf("len(signature) != 2*keySize: %d != %d: %w", len(signature), 2*keySize, ErrInvalidKeyReceived)
	}
	h := cryptoHash.New()
	h.Write([]byte(jwsSigningInput))
	r := big.NewInt(0).SetBytes(signature[:keySize])
	s := big.NewInt(0).SetBytes(signature[keySize:])
	if !ecdsa.Verify(pub, h.Sum(nil), r, s) {
		return ErrFailedToVerifySignature
	}
	return nil
}

func verifyPS(key crypto.PublicKey, hashNewFunc func() hash.Hash, jwsSigningInput string, cryptoHash crypto.Hash, signature []byte, opts *rsa.PSSOptions) error {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hashNewFunc()
	h.Write([]byte(jwsSigningInput))
	if err := rsa.VerifyPSS(pub, cryptoHash, h.Sum(nil), signature, opts); err != nil {
		return fmt.Errorf("rsa.VerifyPKCS1v15: %w", err)
	}
	return nil
}
