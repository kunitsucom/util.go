package jose

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// Header
//
//   - ref. JOSE Header - JSON Web Signature (JWS) https://www.rfc-editor.org/rfc/rfc7515#section-4
//   - ref. JOSE Header - JSON Web Token (JWT)     https://www.rfc-editor.org/rfc/rfc7519#section-5
type Header struct {
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.1
	Algorithm string `json:"alg"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.2
	JwksURL string `json:"jku,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.3
	JSONWebKey string `json:"jwk,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.4
	KeyID string `json:"kid,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.5
	X509URL string `json:"x5u,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.6
	X509CertificateChain string `json:"x5c,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.7
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.8
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.9
	Type string `json:"typ,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.10
	ContentType string `json:"cty,omitempty"`
	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.1.11
	Critical string `json:"crit,omitempty"`

	// ref. https://www.rfc-editor.org/rfc/rfc7515#section-4.3
	PrivateHeaderParameters map[string]interface{} `json:"-"`
}

var ErrInvalidJSON = errors.New("invalid JSON")

func (h *Header) encode() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	if len(h.PrivateHeaderParameters) == 0 {
		return base64.RawURLEncoding.EncodeToString(b), nil
	}

	privateHeaderParameters, err := json.Marshal(h.PrivateHeaderParameters)
	if err != nil {
		return "", fmt.Errorf("invalid private header parameters: %v: %w", h.PrivateHeaderParameters, err)
	}

	if !bytes.HasSuffix(b, []byte{'}'}) {
		return "", fmt.Errorf("%s: %w", b, ErrInvalidJSON)
	}
	if !bytes.HasPrefix(privateHeaderParameters, []byte{'{'}) {
		return "", fmt.Errorf("%s: %w", privateHeaderParameters, ErrInvalidJSON)
	}

	b[len(b)-1] = ','
	b = append(b, privateHeaderParameters[1:]...)
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (h *Header) Encode() (string, error) {
	return h.encode()
}
