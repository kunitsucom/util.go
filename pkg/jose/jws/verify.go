package jws

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/kunitsuinc/util.go/pkg/jose"
	"github.com/kunitsuinc/util.go/pkg/jose/jwa"
	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
)

// - ref. JSON Web Signature (JWS) https://www.rfc-editor.org/rfc/rfc7515

var (
	ErrInvalidTokenReceived = errors.New(`jws: invalid token received, token must have 3 parts`)
	ErrInvalidKeyOption     = errors.New(`jws: invalid key option`)
)

type VerificationKeyOption struct {
	key                         any
	useJSONWebKey               bool
	useJWKSetURLHeaderParameter bool
	useJWKSetURL                bool
	ctx                         context.Context //nolint:containedctx
	jwkSetURL                   string
}

func UseKey(key any) VerificationKeyOption {
	return VerificationKeyOption{
		key: key,
	}
}

func UseHMACKey(key []byte) VerificationKeyOption {
	return UseKey(key)
}

func UseRSAKey(key *rsa.PublicKey) VerificationKeyOption {
	return UseKey(key)
}

func UseECDSAKey(key *ecdsa.PublicKey) VerificationKeyOption {
	return UseKey(key)
}

func UseJSONWebKey() VerificationKeyOption {
	return VerificationKeyOption{
		useJSONWebKey: true,
	}
}

func UseJWKSetURLHeaderParameter(ctx context.Context) VerificationKeyOption {
	return VerificationKeyOption{
		ctx:                         ctx,
		useJWKSetURLHeaderParameter: true,
	}
}

// UseJWKSetURL is a VerificationKeyOption to verify with jwkSetURL.
func UseJWKSetURL(ctx context.Context, jwkSetURL string) VerificationKeyOption {
	return VerificationKeyOption{
		ctx:          ctx,
		useJWKSetURL: true,
		jwkSetURL:    jwkSetURL,
	}
}

func Verify(keyOption VerificationKeyOption, jwt string) (header *jose.Header, err error) {
	headerEncoded, payloadEncoded, signatureEncoded, err := Parse(jwt)
	if err != nil {
		return nil, fmt.Errorf("jws.ParseHeader: %w", err)
	}

	h := new(jose.Header)
	if err := h.Decode(headerEncoded); err != nil {
		return nil, fmt.Errorf("(*jose.Header).Decode: %w", err)
	}

	signingInput := headerEncoded + "." + payloadEncoded

	if keyOption.key != nil {
		return nil, verifyWithKey(h.Algorithm, keyOption.key, signingInput, signatureEncoded)
	}

	if keyOption.useJSONWebKey {
		return nil, verifyWithJSONWebKey(h, signingInput, signatureEncoded)
	}

	if keyOption.useJWKSetURLHeaderParameter {
		return nil, verifyWithJWKSetURL(keyOption.ctx, h.JWKSetURL, h, signingInput, signatureEncoded)
	}

	if keyOption.useJWKSetURL {
		return nil, verifyWithJWKSetURL(keyOption.ctx, keyOption.jwkSetURL, h, signingInput, signatureEncoded)
	}

	return nil, ErrInvalidKeyOption
}

func verifyWithKey(alg string, key any, signingInput, signatureEncoded string) error {
	return jwa.JWS(alg).Verify(key, signingInput, signatureEncoded) //nolint:wrapcheck
}

func verifyWithJSONWebKey(header *jose.Header, signingInput, signatureEncoded string) error {
	if header.JSONWebKey == nil {
		return jose.ErrJSONWebKeyIsEmpty
	}

	pub, err := header.JSONWebKey.DecodePublicKey()
	if err != nil {
		return fmt.Errorf("header.JSONWebKey.DecodePublicKey: %w", err)
	}

	return jwa.JWS(header.Algorithm).Verify(pub, signingInput, signatureEncoded) //nolint:wrapcheck
}

func verifyWithJWKSetURL(ctx context.Context, jwkSetURL string, header *jose.Header, signingInput, signatureEncoded string) error {
	jwks, err := jwk.GetJWKSet(ctx, jwkSetURL)
	if err != nil {
		return fmt.Errorf("jwk.GetJWKSet: %w", err)
	}

	var jsonWebKey *jwk.JSONWebKey

	if len(jwks.Keys) == 1 {
		jsonWebKey = jwks.Keys[0]
	} else {
		key, err := jwks.GetJSONWebKey(header.KeyID)
		if err != nil {
			return fmt.Errorf("jwks.GetJSONWebKey: %w", err)
		}
		jsonWebKey = key
	}

	pub, err := jsonWebKey.DecodePublicKey()
	if err != nil {
		return fmt.Errorf("jsonWebKey.DecodePublicKey: %w", err)
	}

	return jwa.JWS(header.Algorithm).Verify(pub, signingInput, signatureEncoded) //nolint:wrapcheck
}
