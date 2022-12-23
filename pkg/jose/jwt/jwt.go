package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//   - ref. JOSE Header - JSON Web Token (JWT) https://www.rfc-editor.org/rfc/rfc7519#section-5

// Claims
//   - ref. RFC 7519 - JSON Web Token (JWT) https://www.rfc-editor.org/rfc/rfc7519#section-4.1
type Claims struct {
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.1
	Issuer string `json:"iss,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.2
	Subject string `json:"sub,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.3
	Audience string `json:"aud,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.4
	ExpirationTime int64 `json:"exp,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.5
	NotBefore int64 `json:"nbf,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.6
	IssuedAt int64 `json:"iat,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.1.7
	JWTID string `json:"jti,omitempty"`

	// ref. https://www.rfc-editor.org/rfc/rfc7519#section-4.3
	PrivateClaims PrivateClaims `json:"-"`
}

type PrivateClaims map[string]any

type ClaimsOption func(c *Claims)

func WithIssuer(iss string) ClaimsOption {
	return func(c *Claims) {
		c.Issuer = iss
	}
}

func WithSubject(sub string) ClaimsOption {
	return func(c *Claims) {
		c.Subject = sub
	}
}

func WithAudience(aud string) ClaimsOption {
	return func(c *Claims) {
		c.Audience = aud
	}
}

func WithExpirationTime(exp time.Time) ClaimsOption {
	return func(c *Claims) {
		c.ExpirationTime = exp.Unix()
	}
}

func WithNotBefore(nbf time.Time) ClaimsOption {
	return func(c *Claims) {
		c.NotBefore = nbf.Unix()
	}
}

func WithIssuedAt(iat time.Time) ClaimsOption {
	return func(c *Claims) {
		c.IssuedAt = iat.Unix()
	}
}

func WithJWTID(jti string) ClaimsOption {
	return func(c *Claims) {
		c.JWTID = jti
	}
}

func WithPrivateClaim(name string, value any) ClaimsOption {
	return func(c *Claims) {
		c.PrivateClaims[name] = value
	}
}

func NewClaims(opts ...ClaimsOption) *Claims {
	c := &Claims{
		PrivateClaims: make(PrivateClaims),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

var ErrInvalidJSON = errors.New("jwt: invalid JSON")

func (c *Claims) UnmarshalJSON(data []byte) (err error) {
	// avoid recursion
	type _Claims Claims
	_claims := _Claims{}

	err = json.Unmarshal(data, &_claims)
	if err == nil {
		*c = Claims(_claims)
	}

	privateClaims := make(map[string]any)

	err = json.Unmarshal(data, &privateClaims)
	if err == nil {
		typ := reflect.TypeOf(_claims)
		for i := 0; i < typ.NumField(); i++ {
			delete(privateClaims, strings.Split(typ.Field(i).Tag.Get("json"), ",")[0])
		}

		c.PrivateClaims = privateClaims
	}

	return err //nolint:wrapcheck
}

func (c *Claims) MarshalJSON() (data []byte, err error) {
	return c.marshalJSON(json.Marshal, bytes.HasSuffix, bytes.HasPrefix)
}

func (c *Claims) marshalJSON(
	json_Marshal func(v any) ([]byte, error), //nolint:revive,stylecheck
	bytes_HasSuffix func(s []byte, suffix []byte) bool, //nolint:revive,stylecheck
	bytes_HasPrefix func(s []byte, prefix []byte) bool, //nolint:revive,stylecheck
) (data []byte, err error) {
	// avoid recursion
	type _Claims Claims
	_claims := _Claims(*c)

	b, err := json_Marshal(&_claims)
	if err != nil {
		return nil, fmt.Errorf("invalid header: %+v: %w", _claims, err)
	}

	if len(c.PrivateClaims) == 0 {
		return b, nil
	}

	privateClaims, err := json.Marshal(c.PrivateClaims)
	if err != nil {
		return nil, fmt.Errorf("invalid private claims: %+v: %w", c.PrivateClaims, err)
	}

	if !bytes_HasSuffix(b, []byte{'}'}) {
		return nil, fmt.Errorf("%s: %w", b, ErrInvalidJSON)
	}

	if !bytes_HasPrefix(privateClaims, []byte{'{'}) {
		return nil, fmt.Errorf("%s: %w", privateClaims, ErrInvalidJSON)
	}

	b[len(b)-1] = ','
	return append(b, privateClaims[1:]...), nil
}

func (c *Claims) Encode() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
