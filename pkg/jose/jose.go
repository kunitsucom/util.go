package jose

// Header
//
//   - ref. JOSE Header - JSON Web Signature (JWS) https://datatracker.ietf.org/doc/html/rfc7515#section-4
//   - ref. JOSE Header - JSON Web Token (JWT)     https://datatracker.ietf.org/doc/html/rfc7519#section-5
type Header struct {
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.1
	Algorithm string `json:"alg"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.2
	JwksURL string `json:"jku,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.3
	JSONWebKey string `json:"jwk,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.4
	KeyID string `json:"kid,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.5
	X509URL string `json:"x5u,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.6
	X509CertificateChain string `json:"x5c,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.7
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.8
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.9
	Type string `json:"typ,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.10
	ContentType string `json:"cty,omitempty"`
	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.11
	Critical string `json:"crit,omitempty"`

	// ref. https://datatracker.ietf.org/doc/html/rfc7515#section-4.3
	PrivateHeaderParameters map[string]string `json:"-"`
}
