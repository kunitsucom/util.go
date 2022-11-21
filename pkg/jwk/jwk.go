package jwk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kunitsuinc/util.go/pkg/net/http/urlcache"
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
	urlcache *urlcache.Client[*JWKSet]
}

func NewClient(opts ...ClientOption) *Client {
	d := &Client{
		urlcache: urlcache.NewClient[*JWKSet](http.DefaultClient),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

type ClientOption func(*Client)

func WithURLCacheClient(client *urlcache.Client[*JWKSet]) ClientOption {
	return func(d *Client) {
		d.urlcache = client
	}
}

func (d *Client) GetJWKSet(ctx context.Context, jwksURI JWKsURI) (*JWKSet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURI, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := d.urlcache.Do(req, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (*JWKSet, error) {
		r := new(JWKSet)
		if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
			return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.JWKSet): %w", err)
		}
		return r, nil
	})
	if err != nil {
		return nil, fmt.Errorf("(*urlcache.Client).Do: %w", err)
	}

	return resp, nil
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
