package jwk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kunitsuinc/util.go/pkg/cache"
	"github.com/kunitsuinc/util.go/pkg/slice"
)

// ref. JSON Web Key (JWK) https://www.rfc-editor.org/rfc/rfc7517

type JWKsURI = string

// ref. JWK Set Format https://www.rfc-editor.org/rfc/rfc7517#section-5
// ref. https://openid-foundation-japan.github.io/rfc7517.ja.html#JWKSet
type JWKSet struct { //nolint:revive
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-5.1
	Keys []*JSONWebKey `json:"keys"`
}

// ref. JSON Web Key (JWK) Format https://www.rfc-editor.org/rfc/rfc7517#section-4
// ref. https://openid-foundation-japan.github.io/rfc7517.ja.html#JWKFormat
type JSONWebKey struct {
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.1
	KeyType string `json:"kty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.2
	PublicKeyUse string `json:"use,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.3
	KeyOperations string `json:"key_ops,omitempty"` //nolint:tagliatelle
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.4
	Algorithm string `json:"alg,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.5
	KeyID string `json:"kid,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.6
	X509URL string `json:"x5u,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.7
	X509CertificateChain string `json:"x5c,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.8
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.9
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle

	// RSA
	// ref. https://www.rfc-editor.org/rfc/rfc7517#section-9.3

	// N is modulus
	N string `json:"n,omitempty"`
	// E is public exponent
	E string `json:"e,omitempty"`
	// D is private exponent
	D string `json:"d,omitempty"`
}

type Client struct { //nolint:revive
	client     *http.Client
	cacheStore *cache.Store[*JWKSet]
}

func NewClient(opts ...ClientOption) *Client {
	d := &Client{
		client:     http.DefaultClient,
		cacheStore: cache.NewStore[*JWKSet](),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(d *Client) {
		d.client = client
	}
}

func WithCacheStore(store *cache.Store[*JWKSet]) ClientOption {
	return func(d *Client) {
		d.cacheStore = store
	}
}

var ErrResponseIsNotCacheable = errors.New("jwk: response is not cacheable")

func (d *Client) GetJWKSet(ctx context.Context, jwksURI JWKsURI) (*JWKSet, error) {
	return d.cacheStore.GetOrSet(jwksURI, func() (*JWKSet, error) { //nolint:wrapcheck
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURI, nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		resp, err := d.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("(*http.Client).Do: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || 300 <= resp.StatusCode {
			body, _ := io.ReadAll(resp.Body)
			bodyCutOff := slice.CutOff(body, 100)
			return nil, fmt.Errorf("code=%d body=%q: %w", resp.StatusCode, string(bodyCutOff), ErrResponseIsNotCacheable)
		}

		r := new(JWKSet)
		if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
			return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.JWKSet): %w", err)
		}

		return r, nil
	})
}

//nolint:gochecknoglobals
var (
	Default = NewClient()
)

func GetJWKSet(ctx context.Context, jwksURI JWKsURI) (*JWKSet, error) {
	return Default.GetJWKSet(ctx, jwksURI)
}

var ErrKidNotFound = errors.New("jwk: kid not found in jwks_uri")

func (jwks *JWKSet) GetJSONWebKey(kid string) (*JSONWebKey, error) {
	for _, jwk := range jwks.Keys {
		if jwk.KeyID == kid {
			return jwk, nil
		}
	}

	return nil, ErrKidNotFound
}
