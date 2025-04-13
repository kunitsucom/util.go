package jwk

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"time"

	slicez "github.com/kunitsucom/util.go/slices"
	syncz "github.com/kunitsucom/util.go/sync"
)

var (
	ErrCurveNotSupported      = errors.New("jwk: specified curve parameter is not supported")
	ErrKeyIsNotForAlgorithm   = errors.New("jwk: key is not for algorithm")
	ErrResponseIsNotCacheable = errors.New("jwk: response is not cacheable")
	ErrInvalidKey             = errors.New("jwk: invalid key")
)

// ref. JSON Web Key (JWK) https://www.rfc-editor.org/rfc/rfc7517

type JWKSetURL = string //nolint:revive

// JWKSet: A JWK Set is a JSON object that represents a set of JWKs.
//
//   - ref. JWK Set Format https://www.rfc-editor.org/rfc/rfc7517#section-5
//   - ref. https://openid-foundation-japan.github.io/rfc7517.ja.html#JWKSet
type JWKSet struct { //nolint:revive
	// Keys: "keys" parameter is an array of JWK values.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-5.1
	Keys []*JSONWebKey `json:"keys"`
}

// JSONWebKey
//
//   - ref. JSON Web Key (JWK) Format https://www.rfc-editor.org/rfc/rfc7517#section-4
//   - ref. https://openid-foundation-japan.github.io/rfc7517.ja.html#JWKFormat
type JSONWebKey struct {
	// KeyType: "kty" parameter identifies the cryptographic algorithm family used with the key, such as "RSA" or "EC".
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.1
	KeyType string `json:"kty"`

	// PublicKeyUse: "use" parameter identifies the intended use of the public key.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.2
	PublicKeyUse string `json:"use,omitempty"`

	// KeyOperations: "key_ops" parameter identifies the operation(s) for which the key is intended to be used.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.3
	KeyOperations []string `json:"key_ops,omitempty"` //nolint:tagliatelle

	// Algorithm: "alg" parameter identifies the algorithm intended for use with the key.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.4
	Algorithm string `json:"alg,omitempty"`

	// KeyID
	//
	// The "kid" (key ID) parameter is used to match a specific key.  This
	// is used, for instance, to choose among a set of keys within a JWK Set
	// during key rollover.  The structure of the "kid" value is
	// unspecified.  When "kid" values are used within a JWK Set, different
	// keys within the JWK Set SHOULD use distinct "kid" values.  (One
	// example in which different keys might use the same "kid" value is if
	// they have different "kty" (key type) values but are considered to be
	// equivalent alternatives by the application using them.)  The "kid"
	// value is a case-sensitive string.  Use of this member is OPTIONAL.
	// When used with JWS or JWE, the "kid" value is used to match a JWS or
	// JWE "kid" Header Parameter value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.5
	KeyID string `json:"kid,omitempty"`

	// X509URL: "x5u" parameter is a URI [RFC3986] that refers to a resource for an X.509 public key certificate or certificate chain [RFC5280].
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.6
	X509URL string `json:"x5u,omitempty"`

	// X509CertificateChain: "x5c" parameter contains a chain of one or more PKIX certificates [RFC5280].
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.7
	X509CertificateChain []string `json:"x5c,omitempty"`

	// X509CertificateSHA1Thumbprint: "x5t" parameter is a base64url-encoded SHA-1 thumbprint (a.k.a. digest) of the DER encoding of an X.509 certificate [RFC5280].
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.8
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`

	// X509CertificateSHA256Thumbprint: "x5t#S256" parameter is a base64url-encoded SHA-256 thumbprint (a.k.a. digest) of the DER encoding of an X.509 certificate [RFC5280].
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7517#section-4.9
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle

	//
	// Parameters for Elliptic Curve Keys
	// ==================================
	//
	// JWKs can represent Elliptic Curve [DSS] keys.  In this case, the
	// "kty" member value is "EC".
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2
	//

	// Crv
	//
	// Parameters for Elliptic Curve Keys
	//
	// The "crv" (curve) parameter identifies the cryptographic curve used
	// with the key.  Curve values from [DSS] used by this specification
	// are:
	//
	//	o  "P-256"
	//	o  "P-384"
	//	o  "P-521"
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.1
	Crv string `json:"crv,omitempty"`

	// X
	//
	// Parameters for Elliptic Curve Keys
	//
	// The "x" (x coordinate) parameter contains the x coordinate for the
	// Elliptic Curve point.  It is represented as the base64url encoding of
	// the octet string representation of the coordinate, as defined in
	// Section 2.3.5 of SEC1 [SEC1].  The length of this octet string MUST
	// be the full size of a coordinate for the curve specified in the "crv"
	// parameter.  For example, if the value of "crv" is "P-521", the octet
	// string must be 66 octets long.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.2
	X string `json:"x,omitempty"`

	// Y
	//
	// Parameters for Elliptic Curve Keys
	//
	// The "y" (y coordinate) parameter contains the y coordinate for the
	// Elliptic Curve point.  It is represented as the base64url encoding of
	// the octet string representation of the coordinate, as defined in
	// Section 2.3.5 of SEC1 [SEC1].  The length of this octet string MUST
	// be the full size of a coordinate for the curve specified in the "crv"
	// parameter.  For example, if the value of "crv" is "P-521", the octet
	// string must be 66 octets long.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.3
	Y string `json:"y,omitempty"`

	//
	// Parameters for RSA Keys
	// =======================
	//
	// JWKs can represent RSA [RFC3447] keys.  In this case, the "kty"
	// member value is "RSA".  The semantics of the parameters defined below
	// are the same as those defined in Sections 3.1 and 3.2 of RFC 3447.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3
	//

	// N
	//
	// Parameters for RSA Keys
	//
	// The "n" (modulus) parameter contains the modulus value for the RSA
	// public key.  It is represented as a Base64urlUInt-encoded value.
	//
	// Note that implementers have found that some cryptographic libraries
	// prefix an extra zero-valued octet to the modulus representations they
	// return, for instance, returning 257 octets for a 2048-bit key, rather
	// than 256.  Implementations using such libraries will need to take
	// care to omit the extra octet from the base64url-encoded
	// representation.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.1
	N string `json:"n,omitempty"`

	// E
	//
	// Parameters for RSA Keys
	//
	// The "e" (exponent) parameter contains the exponent value for the RSA
	// public key.  It is represented as a Base64urlUInt-encoded value.
	//
	// For instance, when representing the value 65537, the octet sequence
	// to be base64url-encoded MUST consist of the three octets [1, 0, 1];
	// the resulting representation for this value is "AQAB".
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.2
	E string `json:"e,omitempty"`

	// P
	//
	// Parameters for RSA Keys
	//
	// The "p" (first prime factor) parameter contains the first prime
	// factor.  It is represented as a Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.2
	P string `json:"p,omitempty"`

	// Q
	//
	// Parameters for RSA Keys
	//
	// The "q" (second prime factor) parameter contains the second prime
	// factor.  It is represented as a Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.3
	Q string `json:"q,omitempty"`

	// DP
	//
	// Parameters for RSA Keys
	//
	// The "dp" (first factor CRT exponent) parameter contains the Chinese
	// Remainder Theorem (CRT) exponent of the first factor.  It is
	// represented as a Base64urlUInt-encoded value.
	//
	// Parameters for RSA Keys
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.4
	DP string `json:"dp,omitempty"`

	// DQ
	//
	// Parameters for RSA Keys
	//
	// The "dq" (second factor CRT exponent) parameter contains the CRT
	// exponent of the second factor.  It is represented as a Base64urlUInt-
	// encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.5
	DQ string `json:"dq,omitempty"`

	// QI
	//
	// Parameters for RSA Keys
	//
	// The "qi" (first CRT coefficient) parameter contains the CRT
	// coefficient of the second factor.  It is represented as a
	// Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.6
	QI string `json:"qi,omitempty"`

	// Oth
	//
	// Parameters for RSA Keys
	//
	// The "oth" (other primes info) parameter contains an array of
	// information about any third and subsequent primes, should they exist.
	// When only two primes have been used (the normal case), this parameter
	// MUST be omitted.  When three or more primes have been used, the
	// number of array elements MUST be the number of primes used minus two.
	// For more information on this case, see the description of the
	// OtherPrimeInfo parameters in Appendix A.1.2 of RFC 3447 [RFC3447],
	// upon which the following parameters are modeled.  If the consumer of
	// a JWK does not support private keys with more than two primes and it
	// encounters a private key that includes the "oth" parameter, then it
	// MUST NOT use the key.  Each array element MUST be an object with the
	// following members.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7
	Oth []OtherPrimesInfo `json:"oth,omitempty"`

	//
	// Parameters for Elliptic Curve Keys or RSA Keys
	// ==============================================
	//

	// D is "ECC private key" for EC, or "private exponent" for RSA
	//
	// Parameters for RSA Private Keys
	//
	// The "d" (ECC private key) parameter contains the Elliptic Curve
	// private key value.  It is represented as the base64url encoding of
	// the octet string representation of the private key value, as defined
	// in Section 2.3.7 of SEC1 [SEC1].  The length of this octet string
	// MUST be ceiling(log-base-2(n)/8) octets (where n is the order of the
	// curve).
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2.2.1
	//
	// Parameters for Elliptic Curve Private Keys
	//
	// The "d" (private exponent) parameter contains the private exponent
	// value for the RSA private key.  It is represented as a Base64urlUInt-
	// encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.1
	//
	D string `json:"d,omitempty"`

	//
	// Parameters for Symmetric Keys
	// ==================================
	//
	// When the JWK "kty" member value is "oct" (octet sequence), the member
	// "k" (see Section 6.4.1) is used to represent a symmetric key (or
	// another key whose value is a single octet sequence).  An "alg" member
	// SHOULD also be present to identify the algorithm intended to be used
	// with the key, unless the application uses another means or convention
	// to determine the algorithm used.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.2
	//

	// K
	//
	// Parameters for Symmetric Keys
	//
	// The "k" (key value) parameter contains the value of the symmetric (or
	// other single-valued) key.  It is represented as the base64url
	// encoding of the octet sequence containing the key value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.4.1
	K string `json:"k,omitempty"`
}

// OtherPrimesInfo is member struct of "oth" (other primes info).
type OtherPrimesInfo struct {
	// PrimeFactor
	//
	// The "r" (prime factor) parameter within an "oth" array member
	// represents the value of a subsequent prime factor.  It is represented
	// as a Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.1
	PrimeFactor string `json:"r,omitempty"`

	// FactorCRTExponent
	//
	// The "d" (factor CRT exponent) parameter within an "oth" array member
	// represents the CRT exponent of the corresponding prime factor.  It is
	// represented as a Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.2
	FactorCRTExponent string `json:"d,omitempty"`

	// FactorCRTCoefficient
	//
	// The "t" (factor CRT coefficient) parameter within an "oth" array
	// member represents the CRT coefficient of the corresponding prime
	// factor.  It is represented as a Base64urlUInt-encoded value.
	//
	//   - ref. https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.3
	FactorCRTCoefficient string `json:"t,omitempty"`
}

type JSONWebKeyOption func(jwk *JSONWebKey)

func WithKeyType(kty string) JSONWebKeyOption {
	return func(jwk *JSONWebKey) {
		jwk.KeyType = kty
	}
}

func WithKeyID(kid string) JSONWebKeyOption {
	return func(jwk *JSONWebKey) {
		jwk.KeyID = kid
	}
}

func WithAlgorithm(alg string) JSONWebKeyOption {
	return func(jwk *JSONWebKey) {
		jwk.Algorithm = alg
	}
}

// TODO: WithPublicKeyUse() and so on

func (jwk *JSONWebKey) EncodeRSAPublicKey(key *rsa.PublicKey, opts ...JSONWebKeyOption) *JSONWebKey {
	if jwk == nil {
		jwk = new(JSONWebKey)
	}
	for _, opt := range opts {
		opt(jwk)
	}
	jwk.KeyType = "RSA"
	jwk.N = base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	jwk.E = base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes())
	return jwk
}

func (jwk *JSONWebKey) DecodeRSAPublicKey() (*rsa.PublicKey, error) {
	nByte, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.N=%s: %w", jwk.N, err)
	}
	n := big.NewInt(0).SetBytes(nByte)

	eByte, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.E=%s: %w", jwk.E, err)
	}
	eBigInt := big.NewInt(0).SetBytes(eByte)
	if eBigInt.Uint64() > math.MaxInt {
		return nil, fmt.Errorf("e=%d: %w", eBigInt.Uint64(), ErrInvalidKey)
	}
	e := int(eBigInt.Uint64())

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

func (jwk *JSONWebKey) EncodeRSAPrivateKey(key *rsa.PrivateKey, opts ...JSONWebKeyOption) *JSONWebKey {
	jwk = jwk.EncodeRSAPublicKey(&key.PublicKey)
	for _, opt := range opts {
		opt(jwk)
	}
	jwk.D = base64.RawURLEncoding.EncodeToString(key.D.Bytes())
	jwk.P = base64.RawURLEncoding.EncodeToString(key.Primes[0].Bytes())
	jwk.Q = base64.RawURLEncoding.EncodeToString(key.Primes[1].Bytes()) // TODO: implement for case of len() less than 1
	return jwk
}

func (jwk *JSONWebKey) DecodeRSAPrivateKey() (*rsa.PrivateKey, error) {
	pub, err := jwk.DecodeRSAPublicKey()
	if err != nil {
		return nil, err
	}

	d, err := base64.RawURLEncoding.DecodeString(jwk.D)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.D: %w", err)
	}

	p, err := base64.RawURLEncoding.DecodeString(jwk.P)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.P: %w", err)
	}

	q, err := base64.RawURLEncoding.DecodeString(jwk.Q)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.Q: %w", err)
	}

	return &rsa.PrivateKey{
		PublicKey: *pub,
		D:         big.NewInt(0).SetBytes(d),
		Primes: []*big.Int{
			big.NewInt(0).SetBytes(p),
			big.NewInt(0).SetBytes(q),
		},
	}, nil
}

func (jwk *JSONWebKey) EncodeECDSAPublicKey(key *ecdsa.PublicKey, opts ...JSONWebKeyOption) *JSONWebKey {
	if jwk == nil {
		jwk = new(JSONWebKey)
	}
	for _, opt := range opts {
		opt(jwk)
	}
	jwk.Crv = key.Params().Name
	jwk.X = base64.RawURLEncoding.EncodeToString(key.X.Bytes())
	jwk.Y = base64.RawURLEncoding.EncodeToString(key.Y.Bytes())
	return jwk
}

func (jwk *JSONWebKey) DecodeECDSAPublicKey() (*ecdsa.PublicKey, error) {
	var crv elliptic.Curve
	switch jwk.Crv {
	case "P-256":
		crv = elliptic.P256()
	case "P-384":
		crv = elliptic.P384()
	case "P-521":
		crv = elliptic.P521()
	default:
		return nil, fmt.Errorf("crv=%s: %w", jwk.Crv, ErrCurveNotSupported)
	}

	x, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.X=%s: %w", jwk.X, err)
	}

	y, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.Y=%s: %w", jwk.Y, err)
	}

	return &ecdsa.PublicKey{
		Curve: crv,
		X:     big.NewInt(0).SetBytes(x),
		Y:     big.NewInt(0).SetBytes(y),
	}, nil
}

func (jwk *JSONWebKey) EncodeECDSAPrivateKey(key *ecdsa.PrivateKey, opts ...JSONWebKeyOption) *JSONWebKey {
	jwk = jwk.EncodeECDSAPublicKey(&key.PublicKey)
	for _, opt := range opts {
		opt(jwk)
	}
	jwk.D = base64.RawURLEncoding.EncodeToString(key.D.Bytes())
	return jwk
}

func (jwk *JSONWebKey) DecodeECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	pub, err := jwk.DecodeECDSAPublicKey()
	if err != nil {
		return nil, err
	}

	d, err := base64.RawURLEncoding.DecodeString(jwk.D)
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: JSONWebKey.D: %w", err)
	}

	return &ecdsa.PrivateKey{
		PublicKey: *pub,
		D:         big.NewInt(0).SetBytes(d),
	}, nil
}

func (jwk *JSONWebKey) DecodePublicKey() (crypto.PublicKey, error) {
	switch jwk.KeyType {
	case "RSA":
		key, err := jwk.DecodeRSAPublicKey()
		if err != nil {
			return nil, fmt.Errorf("jwk.DecodeRSAPublicKey: %w", err)
		}
		return key, nil
	case "EC":
		key, err := jwk.DecodeECDSAPublicKey()
		if err != nil {
			return nil, fmt.Errorf("jwk.DecodeRSAPublicKey: %w", err)
		}
		return key, nil
	}

	return nil, fmt.Errorf("kty=%s: %w", jwk.KeyType, ErrKeyIsNotForAlgorithm)
}

// TODO: implement
// func (jwk *JSONWebKey) DecodePrivateKey(alg string) (crypto.PrivateKey, error) {
// }

type Client struct { //nolint:revive
	client   *http.Client
	cacheMap syncz.Map[*JWKSet]
}

func NewClient(ctx context.Context, opts ...ClientOption) *Client {
	const defaultTTL = 10 * time.Minute
	d := &Client{
		client: &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				// do not redirect for avoiding open redirect from jku.
				return http.ErrUseLastResponse
			},
		},
		cacheMap: syncz.NewMap[*JWKSet](ctx, syncz.WithNewMapOptionTTL(defaultTTL)),
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

func WithCacheMap(cacheMap syncz.Map[*JWKSet]) ClientOption {
	return func(d *Client) {
		d.cacheMap = cacheMap
	}
}

func (d *Client) GetJWKSet(ctx context.Context, jwksURL JWKSetURL) (*JWKSet, error) {
	if jwks, ok := d.cacheMap.Load(jwksURL); ok {
		return jwks, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
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
		const cutOffSize = 100
		bodyCutOff := slicez.CutOff(body, cutOffSize)
		return nil, fmt.Errorf("code=%d body=%q: %w", resp.StatusCode, string(bodyCutOff), ErrResponseIsNotCacheable)
	}

	r := new(JWKSet)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.JWKSet): %w", err)
	}

	d.cacheMap.Store(jwksURL, r)
	return r, nil
}

//nolint:gochecknoglobals
var (
	Default = NewClient(context.Background())
)

func GetJWKSet(ctx context.Context, jwksURL JWKSetURL) (*JWKSet, error) {
	return Default.GetJWKSet(ctx, jwksURL)
}

var ErrKidNotFound = errors.New("jwk: kid not found in jwks")

func (jwks *JWKSet) GetJSONWebKey(kid string) (*JSONWebKey, error) {
	for _, jwk := range jwks.Keys {
		if jwk.KeyID == kid {
			return jwk, nil
		}
	}

	return nil, fmt.Errorf("kid=%s: %w", kid, ErrKidNotFound)
}
