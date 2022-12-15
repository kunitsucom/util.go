package jws

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
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

func VerifySignature(token string, key crypto.PublicKey) error { //nolint:funlen,cyclop
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

	signingInput := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: signature: %w", err)
	}

	// "alg" (Algorithm) Header Parameter Values for JWS - JSON Web Algorithms (JWA) ref. https://www.rfc-editor.org/rfc/rfc7518#section-3.1
	switch header.Algorithm {
	case jwa.HS256:
		if err := verifyHS(signature, signingInput, key, sha256.New); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.HS384:
		if err := verifyHS(signature, signingInput, key, sha512.New384); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.HS512:
		if err := verifyHS(signature, signingInput, key, sha512.New); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.RS256:
		if err := verifyRS(signature, signingInput, key, sha256.New, crypto.SHA256); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.RS384:
		if err := verifyRS(signature, signingInput, key, sha512.New384, crypto.SHA384); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.RS512:
		if err := verifyRS(signature, signingInput, key, sha512.New, crypto.SHA512); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.ES256:
		if err := verifyES(signature, signingInput, key, crypto.SHA256, 32); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.ES384:
		if err := verifyES(signature, signingInput, key, crypto.SHA384, 48); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.ES512:
		if err := verifyES(signature, signingInput, key, crypto.SHA512, 66); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.PS256:
		if err := verifyPS(signature, signingInput, key, sha256.New, crypto.SHA256, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.PS384:
		if err := verifyPS(signature, signingInput, key, sha512.New384, crypto.SHA384, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.PS512:
		if err := verifyPS(signature, signingInput, key, sha512.New, crypto.SHA512, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
			return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, err)
		}
	case jwa.None:
		return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, ErrAlgorithmNoneIsNotSupported)
	default:
		return fmt.Errorf("alg=%s: key=%T: %w", header.Algorithm, key, ErrInvalidAlgorithm)
	}

	return nil
}

func verifyHS(signature []byte, signingInput string, key crypto.PublicKey, hashNewFunc func() hash.Hash) error {
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hmac.New(hashNewFunc, keyBytes)
	h.Write([]byte(signingInput))
	if !bytes.Equal(signature, h.Sum(nil)) {
		return ErrFailedToVerifySignature
	}
	return nil
}

func verifyRS(signature []byte, signingInput string, key crypto.PublicKey, hashNewFunc func() hash.Hash, cryptoHash crypto.Hash) error {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hashNewFunc()
	h.Write([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(pub, cryptoHash, h.Sum(nil), signature); err != nil {
		return fmt.Errorf("rsa.VerifyPKCS1v15: %w", err)
	}
	return nil
}

func verifyES(signature []byte, signingInput string, key crypto.PublicKey, cryptoHash crypto.Hash, keySize int) error {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	if len(signature) != keySize*2 {
		return fmt.Errorf("len(signature)=(%d) != keySize*2=(%d): %w", len(signature), 2*keySize, ErrInvalidKeyReceived)
	}
	h := cryptoHash.New()
	h.Write([]byte(signingInput))
	r := big.NewInt(0).SetBytes(signature[:keySize])
	s := big.NewInt(0).SetBytes(signature[keySize:])
	if !ecdsa.Verify(pub, h.Sum(nil), r, s) {
		return ErrFailedToVerifySignature
	}
	return nil
}

func verifyPS(signature []byte, signingInput string, key crypto.PublicKey, hashNewFunc func() hash.Hash, cryptoHash crypto.Hash, opts *rsa.PSSOptions) error {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyReceived
	}
	h := hashNewFunc()
	h.Write([]byte(signingInput))
	if err := rsa.VerifyPSS(pub, cryptoHash, h.Sum(nil), signature, opts); err != nil {
		return fmt.Errorf("rsa.VerifyPKCS1v15: %w", err)
	}
	return nil
}

func Sign(alg jwa.Algorithm, signingInput string, key crypto.PrivateKey) (signature []byte, err error) { //nolint:funlen,cyclop
	// "alg" (Algorithm) Header Parameter Values for JWS - JSON Web Algorithms (JWA) ref. https://www.rfc-editor.org/rfc/rfc7518#section-3.1
	switch alg {
	case jwa.HS256:
		signature, err = signHS(signingInput, key, sha256.New)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.HS384:
		signature, err = signHS(signingInput, key, sha512.New384)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.HS512:
		signature, err = signHS(signingInput, key, sha512.New)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.RS256:
		signature, err = signRS(signingInput, key, sha256.New, crypto.SHA256)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.RS384:
		signature, err = signRS(signingInput, key, sha512.New384, crypto.SHA384)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.RS512:
		signature, err = signRS(signingInput, key, sha512.New, crypto.SHA512)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.ES256:
		signature, err = signES(signingInput, key, crypto.SHA256, 32)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.ES384:
		signature, err = signES(signingInput, key, crypto.SHA384, 48)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.ES512:
		signature, err = signES(signingInput, key, crypto.SHA512, 66)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.PS256:
		signature, err = signPS(signingInput, key, crypto.SHA256)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.PS384:
		signature, err = signPS(signingInput, key, crypto.SHA384)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.PS512:
		signature, err = signPS(signingInput, key, crypto.SHA512)
		if err != nil {
			return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, err)
		}
	case jwa.None:
		return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, ErrAlgorithmNoneIsNotSupported)
	default:
		return nil, fmt.Errorf("alg=%s: key=%T: %w", alg, key, ErrInvalidAlgorithm)
	}

	return signature, nil
}

func signHS(signingInput string, key crypto.PrivateKey, hashNewFunc func() hash.Hash) (signature []byte, err error) {
	keyBytes, ok := key.([]byte)
	if !ok {
		return nil, ErrInvalidKeyReceived
	}
	h := hmac.New(hashNewFunc, keyBytes)
	h.Write([]byte(signingInput))
	return []byte(base64.RawURLEncoding.EncodeToString(h.Sum(nil))), nil
}

func signRS(signingInput string, key crypto.PrivateKey, hashNewFunc func() hash.Hash, cryptoHash crypto.Hash) (signature []byte, err error) {
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyReceived
	}

	h := hashNewFunc()
	h.Write([]byte(signingInput))
	signature, err = rsa.SignPKCS1v15(rand.Reader, priv, cryptoHash, h.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("rsa.SignPKCS1v15: %w", err)
	}
	return []byte(base64.RawURLEncoding.EncodeToString(signature)), nil
}

func signES(signingInput string, key crypto.PrivateKey, cryptoHash crypto.Hash, keySize int) (signature []byte, err error) {
	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyReceived
	}

	h := cryptoHash.New()
	h.Write([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, priv, h.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("ecdsa.Sign: %w", err)
	}

	rBytes := r.Bytes()
	rBytesPadded := make([]byte, keySize)
	copy(rBytesPadded[keySize-len(rBytes):], rBytes)

	sBytes := s.Bytes()
	sBytesPadded := make([]byte, keySize)
	copy(sBytesPadded[keySize-len(sBytes):], sBytes)

	var rawSignature []byte
	rawSignature = append(rawSignature, rBytesPadded...)
	rawSignature = append(rawSignature, sBytesPadded...)

	return []byte(base64.RawURLEncoding.EncodeToString(rawSignature)), nil
}

func signPS(signingInput string, key crypto.PrivateKey, cryptoHash crypto.Hash) (signature []byte, err error) {
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyReceived
	}

	h := cryptoHash.New()
	h.Write([]byte(signingInput))
	signature, err = rsa.SignPSS(rand.Reader, priv, cryptoHash, h.Sum(nil), &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		return nil, fmt.Errorf("rsa.SignPSS: %w", err)
	}
	return []byte(base64.RawURLEncoding.EncodeToString(signature)), nil
}
