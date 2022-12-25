package jws

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
)

// Header
//
//   - ref. JOSE Header - JSON Web Signature (JWS) https://www.rfc-editor.org/rfc/rfc7515#section-4
type Header struct {
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.1
	Algorithm string `json:"alg"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.2
	JwksURL string `json:"jku,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.3
	JSONWebKey *jwk.JSONWebKey `json:"jwk,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.4
	KeyID string `json:"kid,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.5
	X509URL string `json:"x5u,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.6
	X509CertificateChain []string `json:"x5c,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.7
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.8
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.9
	Type string `json:"typ,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.10
	ContentType string `json:"cty,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.11
	Critical []string `json:"crit,omitempty"`

	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.3
	PrivateHeaderParameters PrivateHeaderParameters `json:"-"`
}

type PrivateHeaderParameters map[string]any

func NewHeader(alg string, opts ...HeaderOption) *Header {
	h := &Header{
		Algorithm:               alg,
		PrivateHeaderParameters: make(PrivateHeaderParameters),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

type HeaderOption func(h *Header)

func WithJwksURL(jku string) HeaderOption {
	return func(h *Header) {
		h.JwksURL = jku
	}
}

func WithJSONWebKey(jwk *jwk.JSONWebKey) HeaderOption {
	return func(h *Header) {
		h.JSONWebKey = jwk
	}
}

func WithKeyID(kid string) HeaderOption {
	return func(h *Header) {
		h.KeyID = kid
	}
}

func WithX509URL(x5u string) HeaderOption {
	return func(h *Header) {
		h.X509URL = x5u
	}
}

func WithX509CertificateChain(x5c []string) HeaderOption {
	return func(h *Header) {
		h.X509CertificateChain = x5c
	}
}

func WithX509CertificateSHA1Thumbprint(x5t string) HeaderOption {
	return func(h *Header) {
		h.X509CertificateSHA1Thumbprint = x5t
	}
}

func WithX509CertificateSHA256Thumbprint(x5tS256 string) HeaderOption {
	return func(h *Header) {
		h.X509CertificateSHA256Thumbprint = x5tS256
	}
}

func WithType(typ string) HeaderOption {
	return func(h *Header) {
		h.Type = typ
	}
}

func WithContentType(cty string) HeaderOption {
	return func(h *Header) {
		h.ContentType = cty
	}
}

func WithCritical(crit []string) HeaderOption {
	return func(h *Header) {
		h.Critical = crit
	}
}

func WithPrivateHeaderParameter(name string, value any) HeaderOption {
	return func(h *Header) {
		h.PrivateHeaderParameters[name] = value
	}
}

var ErrInvalidJSON = errors.New("jws: invalid JSON")

func (h *Header) UnmarshalJSON(data []byte) (err error) {
	// avoid recursion
	type _Header Header
	_header := _Header{}

	err = json.Unmarshal(data, &_header)
	if err == nil {
		*h = Header(_header)
	}

	privateHeaderParameters := make(map[string]any)

	err = json.Unmarshal(data, &privateHeaderParameters)
	if err == nil {
		typ := reflect.TypeOf(_header)
		for i := 0; i < typ.NumField(); i++ {
			delete(privateHeaderParameters, strings.Split(typ.Field(i).Tag.Get("json"), ",")[0])
		}

		h.PrivateHeaderParameters = privateHeaderParameters
	}

	return err //nolint:wrapcheck
}

func (h *Header) MarshalJSON() (data []byte, err error) {
	return h.marshalJSON(json.Marshal, bytes.HasSuffix, bytes.HasPrefix)
}

func (h *Header) marshalJSON(
	json_Marshal func(v any) ([]byte, error), //nolint:revive,stylecheck
	bytes_HasSuffix func(s []byte, suffix []byte) bool, //nolint:revive,stylecheck
	bytes_HasPrefix func(s []byte, prefix []byte) bool, //nolint:revive,stylecheck
) (data []byte, err error) {
	// avoid recursion
	type _Header Header
	_header := _Header(*h)

	b, err := json_Marshal(&_header)
	if err != nil {
		return nil, fmt.Errorf("invalid header: %+v: %w", _header, err)
	}

	if len(h.PrivateHeaderParameters) == 0 {
		return b, nil
	}

	privateHeaderParameters, err := json.Marshal(h.PrivateHeaderParameters)
	if err != nil {
		return nil, fmt.Errorf("invalid private header parameters: %+v: %w", h.PrivateHeaderParameters, err)
	}

	if !bytes_HasSuffix(b, []byte{'}'}) {
		return nil, fmt.Errorf("%s: %w", b, ErrInvalidJSON)
	}

	if !bytes_HasPrefix(privateHeaderParameters, []byte{'{'}) {
		return nil, fmt.Errorf("%s: %w", privateHeaderParameters, ErrInvalidJSON)
	}

	b[len(b)-1] = ','
	return append(b, privateHeaderParameters[1:]...), nil
}

func (h *Header) Encode() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Encode(h *Header) (string, error) {
	return h.Encode()
}

func (h *Header) Decode(header string) error {
	decoded, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}

	if err := json.Unmarshal(decoded, h); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func Decode(header string) (*Header, error) {
	h := new(Header)
	return h, h.Decode(header)
}
