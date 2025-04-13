package jose

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kunitsucom/util.go/jose/jwk"
)

var (
	ErrJSONWebKeyIsEmpty                    = errors.New(`jose: jwk is empty`)
	ErrPrivateHeaderParameterIsNotFound     = errors.New(`jose: private header parameter is not found`)
	ErrVIsNotPointerOrInterface             = errors.New(`jose: v is not pointer or interface`)
	ErrPrivateHeaderParameterTypeIsNotMatch = errors.New(`jose: private header parameter type is not match`)
)

// Header
//
// JOSE Header
//
//   - ref. JOSE Header - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4
//   - ref. JOSE Header - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4
//   - ref. JOSE Header - JSON Web Key (JWK)        https://www.rfc-editor.org/rfc/rfc7517#appendix-C.2
//   - ref. JOSE Header - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.6.1
//   - ref. JOSE Header - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.7.1
//   - ref. JOSE Header - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.8.1
//   - ref. JOSE Header - JSON Web Token (JWT)      https://www.rfc-editor.org/rfc/rfc7519#section-5
type Header struct {
	// Algorithm
	//
	// # JSON Web Signature (JWS)
	//
	// The "alg" (algorithm) Header Parameter identifies the cryptographic
	// algorithm used to secure the JWS.  The JWS Signature value is not
	// valid if the "alg" value does not represent a supported algorithm or
	// if there is not a key for use with that algorithm associated with the
	// party that digitally signed or MACed the content.  "alg" values
	// should either be registered in the IANA "JSON Web Signature and
	// Encryption Algorithms" registry established by [JWA] or be a value
	// that contains a Collision-Resistant Name.  The "alg" value is a case-
	// sensitive ASCII string containing a StringOrURI value.  This Header
	// Parameter MUST be present and MUST be understood and processed by
	// implementations.
	//
	// A list of defined "alg" values for this use can be found in the IANA
	// "JSON Web Signature and Encryption Algorithms" registry established
	// by [JWA]; the initial contents of this registry are the values
	// defined in Section 3.1 of [JWA].
	//
	//   - ref. "alg" (Algorithm) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.1
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "alg" Header Parameter defined in Section 4.1.1 of [JWS], except
	// that the Header Parameter identifies the cryptographic algorithm used
	// to encrypt or determine the value of the CEK.  The encrypted content
	// is not usable if the "alg" value does not represent a supported
	// algorithm, or if the recipient does not have a key that can be used
	// with that algorithm.
	//
	// A list of defined "alg" values for this use can be found in the IANA
	// "JSON Web Signature and Encryption Algorithms" registry established
	// by [JWA]; the initial contents of this registry are the values
	// defined in Section 4.1 of [JWA].
	//
	//   - ref. "alg" (Algorithm) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.1
	Algorithm string `json:"alg,omitempty"`

	// EncryptionAlgorithm
	//
	// # JSON Web Encryption (JWE)
	//
	// The "enc" (encryption algorithm) Header Parameter identifies the
	// content encryption algorithm used to perform authenticated encryption
	// on the plaintext to produce the ciphertext and the Authentication
	// Tag.  This algorithm MUST be an AEAD algorithm with a specified key
	// length.  The encrypted content is not usable if the "enc" value does
	// not represent a supported algorithm.  "enc" values should either be
	// registered in the IANA "JSON Web Signature and Encryption Algorithms"
	// registry established by [JWA] or be a value that contains a
	// Collision-Resistant Name.  The "enc" value is a case-sensitive ASCII
	// string containing a StringOrURI value.  This Header Parameter MUST be
	// present and MUST be understood and processed by implementations.
	//
	// A list of defined "enc" values for this use can be found in the IANA
	// "JSON Web Signature and Encryption Algorithms" registry established
	// by [JWA]; the initial contents of this registry are the values
	// defined in Section 5.1 of [JWA].
	//
	//   - ref. "enc" (Encryption Algorithm) Header Parameter                - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.2
	//   - ref. "enc" (Encryption Algorithm) Header Parameter Values for JWE - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-5.1
	EncryptionAlgorithm string `json:"enc,omitempty"`

	// EphemeralPublicKey
	//
	// # JSON Web Encryption (JWE)
	//
	// The "epk" (ephemeral public key) value created by the originator for
	// the use in key agreement algorithms.  This key is represented as a
	// JSON Web Key [JWK] public key value.  It MUST contain only public key
	// parameters and SHOULD contain only the minimum JWK parameters
	// necessary to represent the key; other JWK parameters included can be
	// checked for consistency and honored, or they can be ignored.  This
	// Header Parameter MUST be present and MUST be understood and processed
	// by implementations when these algorithms are used.
	//
	//   - ref. "epk" (Ephemeral Public Key) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.6.1.1
	EphemeralPublicKey string `json:"epk,omitempty"`

	// AgreementPartyUInfo
	//
	// # JSON Web Encryption (JWE)
	//
	// The "apu" (agreement PartyUInfo) value for key agreement algorithms
	// using it (such as "ECDH-ES"), represented as a base64url-encoded
	// string.  When used, the PartyUInfo value contains information about
	// the producer.  Use of this Header Parameter is OPTIONAL.  This Header
	// Parameter MUST be understood and processed by implementations when
	// these algorithms are used.
	//
	//   - ref. "apu" (Agreement PartyUInfo) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.6.1.2
	AgreementPartyUInfo string `json:"apu,omitempty"`

	// AgreementPartyVInfo
	//
	// # JSON Web Encryption (JWE)
	//
	// The "apv" (agreement PartyVInfo) value for key agreement algorithms
	// using it (such as "ECDH-ES"), represented as a base64url encoded
	// string.  When used, the PartyVInfo value contains information about
	// the recipient.  Use of this Header Parameter is OPTIONAL.  This
	// Header Parameter MUST be understood and processed by implementations
	// when these algorithms are used.
	//
	//   - ref. "apv" (Agreement PartyVInfo) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.6.1.3
	AgreementPartyVInfo string `json:"apv,omitempty"`

	// InitializationVector
	//
	// # JSON Web Encryption (JWE)
	//
	// The "iv" (initialization vector) Header Parameter value is the
	// base64url-encoded representation of the 96-bit IV value used for the
	// key encryption operation.  This Header Parameter MUST be present and
	// MUST be understood and processed by implementations when these
	// algorithms are used.
	//
	//   - ref. "iv" (Initialization Vector) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.7.1.1
	InitializationVector string `json:"iv,omitempty"`

	// AuthenticationTag
	//
	// # JSON Web Encryption (JWE)
	//
	// The "tag" (authentication tag) Header Parameter value is the
	// base64url-encoded representation of the 128-bit Authentication Tag
	// value resulting from the key encryption operation.  This Header
	// Parameter MUST be present and MUST be understood and processed by
	// implementations when these algorithms are used.
	//
	//   - ref. "tag" (Authentication Tag) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.7.1.2
	AuthenticationTag string `json:"tag,omitempty"`

	// PBES2SaltInput
	//
	// # JSON Web Encryption (JWE)
	//
	// The "p2s" (PBES2 salt input) Header Parameter encodes a Salt Input
	// value, which is used as part of the PBKDF2 salt value.  The "p2s"
	// value is BASE64URL(Salt Input).  This Header Parameter MUST be
	// present and MUST be understood and processed by implementations when
	// these algorithms are used.
	//
	// The salt expands the possible keys that can be derived from a given
	// password.  A Salt Input value containing 8 or more octets MUST be
	// used.  A new Salt Input value MUST be generated randomly for every
	// encryption operation; see RFC 4086 [RFC4086] for considerations on
	// generating random values.  The salt value used is (UTF8(Alg) || 0x00
	// || Salt Input), where Alg is the "alg" (algorithm) Header Parameter
	// value.
	//
	//   - ref. "p2s" (PBES2 Salt Input) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.8.1.1
	PBES2SaltInput string `json:"p2s,omitempty"`

	// PBES2Count
	//
	// # JSON Web Encryption (JWE)
	//
	// The "p2c" (PBES2 count) Header Parameter contains the PBKDF2
	// iteration count, represented as a positive JSON integer.  This Header
	// Parameter MUST be present and MUST be understood and processed by
	// implementations when these algorithms are used.
	//
	// The iteration count adds computational expense, ideally compounded by
	// the possible range of keys introduced by the salt.  A minimum
	// iteration count of 1000 is RECOMMENDED.
	//
	//   - ref. "p2c" (PBES2 Count) Header Parameter - JSON Web Algorithms (JWA) https://www.rfc-editor.org/rfc/rfc7518#section-4.8.1.2
	PBES2Count string `json:"p2c,omitempty"`

	// CompressionAlgorithm
	//
	// # JSON Web Encryption (JWE)
	//
	// The "zip" (compression algorithm) applied to the plaintext before
	// encryption, if any.  The "zip" value defined by this specification
	// is:
	//
	// o  "DEF" - Compression with the DEFLATE [RFC1951] algorithm
	//
	// Other values MAY be used.  Compression algorithm values can be
	// registered in the IANA "JSON Web Encryption Compression Algorithms"
	// registry established by [JWA].  The "zip" value is a case-sensitive
	// string.  If no "zip" parameter is present, no compression is applied
	// to the plaintext before encryption.  When used, this Header Parameter
	// MUST be integrity protected; therefore, it MUST occur only within the
	// JWE Protected Header.  Use of this Header Parameter is OPTIONAL.
	// This Header Parameter MUST be understood and processed by
	// implementations.
	//
	//   - ref. "zip" (Compression Algorithm) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.3
	CompressionAlgorithm string `json:"zip,omitempty"`

	// JWKSetURL
	//
	// # JSON Web Signature (JWS)
	//
	// The "jku" (JWK Set URL) Header Parameter is a URI [RFC3986] that
	// refers to a resource for a set of JSON-encoded public keys, one of
	// which corresponds to the key used to digitally sign the JWS.  The
	// keys MUST be encoded as a JWK Set [JWK].  The protocol used to
	// acquire the resource MUST provide integrity protection; an HTTP GET
	// request to retrieve the JWK Set MUST use Transport Layer Security
	//
	//   - ref. "jku" (JWK Set URL) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.2
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "jku" Header Parameter defined in Section 4.1.2 of [JWS], except
	// that the JWK Set resource contains the public key to which the JWE
	// was encrypted; this can be used to determine the private key needed
	// to decrypt the JWE.
	//
	//   - ref. "jku" (JWK Set URL) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.4
	JWKSetURL string `json:"jku,omitempty"`

	// JSONWebKey
	//
	// # JSON Web Signature (JWS)
	//
	// The "jwk" (JSON Web Key) Header Parameter is the public key that
	// corresponds to the key used to digitally sign the JWS.  This key is
	// represented as a JSON Web Key [JWK].  Use of this Header Parameter is
	// OPTIONAL.
	//
	//   - ref. "jwk" (JSON Web Key) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.3
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "jwk" Header Parameter defined in Section 4.1.3 of [JWS], except
	// that the key is the public key to which the JWE was encrypted; this
	// can be used to determine the private key needed to decrypt the JWE.
	//
	//   - ref. "jwk" (JSON Web Key) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.5
	JSONWebKey *jwk.JSONWebKey `json:"jwk,omitempty"`

	// KeyID
	//
	// # JSON Web Signature (JWS)
	//
	// The "kid" (key ID) Header Parameter is a hint indicating which key
	// was used to secure the JWS.  This parameter allows originators to
	// explicitly signal a change of key to recipients.  The structure of
	// the "kid" value is unspecified.  Its value MUST be a case-sensitive
	// string.  Use of this Header Parameter is OPTIONAL.
	//
	// When used with a JWK, the "kid" value is used to match a JWK "kid"
	// parameter value.
	//
	//   - ref. "kid" (Key ID) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.4
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "kid" Header Parameter defined in Section 4.1.4 of [JWS], except
	// that the key hint references the public key to which the JWE was
	// encrypted; this can be used to determine the private key needed to
	// decrypt the JWE.  This parameter allows originators to explicitly
	// signal a change of key to JWE recipients.
	//
	//   - ref. "kid" (Key ID) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.6
	KeyID string `json:"kid,omitempty"`

	// X509URL
	//
	// # JSON Web Signature (JWS)
	//
	// The "x5u" (X.509 URL) Header Parameter is a URI [RFC3986] that refers
	// to a resource for the X.509 public key certificate or certificate
	// chain [RFC5280] corresponding to the key used to digitally sign the
	// JWS.  The identified resource MUST provide a representation of the
	// certificate or certificate chain that conforms to RFC 5280 [RFC5280]
	// in PEM-encoded form, with each certificate delimited as specified in
	// Section 6.1 of RFC 4945 [RFC4945].  The certificate containing the
	// public key corresponding to the key used to digitally sign the JWS
	// MUST be the first certificate.  This MAY be followed by additional
	// certificates, with each subsequent certificate being the one used to
	// certify the previous one.  The protocol used to acquire the resource
	// MUST provide integrity protection; an HTTP GET request to retrieve
	// the certificate MUST use TLS [RFC2818] [RFC5246]; and the identity of
	// the server MUST be validated, as per Section 6 of RFC 6125 [RFC6125].
	// Also, see Section 8 on TLS requirements.  Use of this Header
	// Parameter is OPTIONAL.
	//
	//   - ref. "x5u" (X.509 URL) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.5
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "x5u" Header Parameter defined in Section 4.1.5 of [JWS], except
	// that the X.509 public key certificate or certificate chain [RFC5280]
	// contains the public key to which the JWE was encrypted; this can be
	// used to determine the private key needed to decrypt the JWE.
	//
	//   - ref. "x5u" (X.509 URL) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.7
	X509URL string `json:"x5u,omitempty"`

	// X509CertificateChain
	//
	// # JSON Web Signature (JWS)
	//
	// The "x5c" (X.509 certificate chain) Header Parameter contains the
	// X.509 public key certificate or certificate chain [RFC5280]
	// corresponding to the key used to digitally sign the JWS.  The
	// certificate or certificate chain is represented as a JSON array of
	//
	//   - ref. "x5c" (X.509 Certificate Chain) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.6
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "x5c" Header Parameter defined in Section 4.1.6 of [JWS], except
	// that the X.509 public key certificate or certificate chain [RFC5280]
	// contains the public key to which the JWE was encrypted; this can be
	// used to determine the private key needed to decrypt the JWE.
	//
	// See Appendix B of [JWS] for an example "x5c" value.
	//
	//   - ref. "x5c" (X.509 Certificate Chain) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.8
	X509CertificateChain []string `json:"x5c,omitempty"`

	// X509CertificateSHA1Thumbprint
	//
	// # JSON Web Signature (JWS)
	//
	// The "x5t" (X.509 certificate SHA-1 thumbprint) Header Parameter is a
	// base64url-encoded SHA-1 thumbprint (a.k.a. digest) of the DER
	// encoding of the X.509 certificate [RFC5280] corresponding to the key
	// used to digitally sign the JWS.  Note that certificate thumbprints
	// are also sometimes known as certificate fingerprints.  Use of this
	// Header Parameter is OPTIONAL.
	//
	//   - ref. "x5t" (X.509 Certificate SHA-1 Thumbprint) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.7
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "x5t" Header Parameter defined in Section 4.1.7 of [JWS], except
	// that the certificate referenced by the thumbprint contains the public
	// key to which the JWE was encrypted; this can be used to determine the
	// private key needed to decrypt the JWE.  Note that certificate
	// thumbprints are also sometimes known as certificate fingerprints.
	//
	//   - ref. "x5t" (X.509 Certificate SHA-1 Thumbprint) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.9
	X509CertificateSHA1Thumbprint string `json:"x5t,omitempty"`

	// X509CertificateSHA256Thumbprint
	//
	// # JSON Web Signature (JWS)
	//
	// The "x5t#S256" (X.509 certificate SHA-256 thumbprint) Header
	// Parameter is a base64url-encoded SHA-256 thumbprint (a.k.a. digest)
	// of the DER encoding of the X.509 certificate [RFC5280] corresponding
	// to the key used to digitally sign the JWS.  Note that certificate
	// thumbprints are also sometimes known as certificate fingerprints.
	// Use of this Header Parameter is OPTIONAL.
	//
	//   - ref. "x5t#S256" (X.509 Certificate SHA-256 Thumbprint) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.8
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "x5t#S256" Header Parameter defined in Section 4.1.8 of [JWS],
	// except that the certificate referenced by the thumbprint contains the
	// public key to which the JWE was encrypted; this can be used to
	// determine the private key needed to decrypt the JWE.  Note that
	// certificate thumbprints are also sometimes known as certificate
	// fingerprints.
	//
	//   - ref. "x5t#S256" (X.509 Certificate SHA-256 Thumbprint) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.10
	X509CertificateSHA256Thumbprint string `json:"x5t#S256,omitempty"` //nolint:tagliatelle

	// Type
	//
	// # JSON Web Signature (JWS)
	//
	// The "typ" (type) Header Parameter is used by JWS applications to
	// declare the media type [IANA.MediaTypes] of this complete JWS.  This
	// is intended for use by the application when more than one kind of
	// object could be present in an application data structure that can
	// contain a JWS; the application can use this value to disambiguate
	// among the different kinds of objects that might be present.  It will
	// typically not be used by applications when the kind of object is
	// already known.  This parameter is ignored by JWS implementations; any
	// processing of this parameter is performed by the JWS application.
	// Use of this Header Parameter is OPTIONAL.
	//
	// Per RFC 2045 [RFC2045], all media type values, subtype values, and
	// parameter names are case insensitive.  However, parameter values are
	// case sensitive unless otherwise specified for the specific parameter.
	//
	// To keep messages compact in common situations, it is RECOMMENDED that
	// producers omit an "application/" prefix of a media type value in a
	// "typ" Header Parameter when no other '/' appears in the media type
	// value.  A recipient using the media type value MUST treat it as if
	// "application/" were prepended to any "typ" value not containing a
	// '/'.  For instance, a "typ" value of "example" SHOULD be used to
	// represent the "application/example" media type, whereas the media
	// type "application/example;part="1/2"" cannot be shortened to
	// "example;part="1/2"".
	//
	// The "typ" value "JOSE" can be used by applications to indicate that
	// this object is a JWS or JWE using the JWS Compact Serialization or
	// the JWE Compact Serialization.  The "typ" value "JOSE+JSON" can be
	// used by applications to indicate that this object is a JWS or JWE
	// using the JWS JSON Serialization or the JWE JSON Serialization.
	// Other type values can also be used by applications.
	//
	//   - ref. "typ" (Type) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.9
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "typ" Header Parameter defined in Section 4.1.9 of [JWS], except
	// that the type is that of this complete JWE.
	//
	//   - ref. "typ" (Type) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.11
	Type string `json:"typ,omitempty"`

	// ContentType
	//
	// # JSON Web Signature (JWS)
	//
	// The "cty" (content type) Header Parameter is used by JWS applications
	// to declare the media type [IANA.MediaTypes] of the secured content
	// (the payload).  This is intended for use by the application when more
	// than one kind of object could be present in the JWS Payload; the
	// application can use this value to disambiguate among the different
	// kinds of objects that might be present.  It will typically not be
	// used by applications when the kind of object is already known.  This
	// parameter is ignored by JWS implementations; any processing of this
	// parameter is performed by the JWS application.  Use of this Header
	// Parameter is OPTIONAL.
	//
	// Per RFC 2045 [RFC2045], all media type values, subtype values, and
	// parameter names are case insensitive.  However, parameter values are
	// case sensitive unless otherwise specified for the specific parameter.
	//
	// To keep messages compact in common situations, it is RECOMMENDED that
	// producers omit an "application/" prefix of a media type value in a
	// "cty" Header Parameter when no other '/' appears in the media type
	// value.  A recipient using the media type value MUST treat it as if
	// "application/" were prepended to any "cty" value not containing a
	// '/'.  For instance, a "cty" value of "example" SHOULD be used to
	// represent the "application/example" media type, whereas the media
	// type "application/example;part="1/2"" cannot be shortened to
	// "example;part="1/2"".
	//
	//   - ref. "cty" (Content Type) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.10
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "cty" Header Parameter defined in Section 4.1.10 of [JWS], except
	// that the type is that of the secured content (the plaintext).
	//
	//   - ref. "cty" (Content Type) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.12
	ContentType string `json:"cty,omitempty"`

	// Critical
	//
	// # JSON Web Signature (JWS)
	//
	// The "crit" (critical) Header Parameter indicates that extensions to
	// this specification and/or [JWA] are being used that MUST be
	// understood and processed.  Its value is an array listing the Header
	// Parameter names present in the JOSE Header that use those extensions.
	// If any of the listed extension Header Parameters are not understood
	// and supported by the recipient, then the JWS is invalid.  Producers
	// MUST NOT include Header Parameter names defined by this specification
	// or [JWA] for use with JWS, duplicate names, or names that do not
	// occur as Header Parameter names within the JOSE Header in the "crit"
	// list.  Producers MUST NOT use the empty list "[]" as the "crit"
	// value.  Recipients MAY consider the JWS to be invalid if the critical
	// list contains any Header Parameter names defined by this
	// specification or [JWA] for use with JWS or if any other constraints
	// on its use are violated.  When used, this Header Parameter MUST be
	// integrity protected; therefore, it MUST occur only within the JWS
	// Protected Header.  Use of this Header Parameter is OPTIONAL.  This
	// Header Parameter MUST be understood and processed by implementations.
	//
	// An example use, along with a hypothetical "exp" (expiration time)
	// field is:
	//
	//   {"alg":"ES256",
	//    "crit":["exp"],
	//    "exp":1363284000
	//   }
	//
	//   - ref. "crit" (Critical) Header Parameter - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.1.11
	//
	// # JSON Web Encryption (JWE)
	//
	// This parameter has the same meaning, syntax, and processing rules as
	// the "crit" Header Parameter defined in Section 4.1.11 of [JWS],
	// except that Header Parameters for a JWE are being referred to, rather
	// than Header Parameters for a JWS.
	//
	//   - ref. "crit" (Critical) Header Parameter - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.1.13
	Critical []string `json:"crit,omitempty"`

	// PrivateHeaderParameters
	//
	// # JSON Web Signature (JWS)
	//
	// A producer and consumer of a JWS may agree to use Header Parameter
	// names that are Private Names (names that are not Registered Header
	// Parameter names (Section 4.1)) or Public Header Parameter names
	// (Section 4.2).  Unlike Public Header Parameter names, Private Header
	// Parameter names are subject to collision and should be used with
	// caution.
	//
	//   - ref. Private Header Parameter Names - JSON Web Signature (JWS)  https://www.rfc-editor.org/rfc/rfc7515#section-4.3
	//
	// # JSON Web Encryption (JWE)
	//
	// A producer and consumer of a JWE may agree to use Header Parameter
	// names that are Private Names: names that are not Registered Header
	// Parameter names (Section 4.1) or Public Header Parameter names
	// (Section 4.2).  Unlike Public Header Parameter names, Private Header
	// Parameter names are subject to collision and should be used with
	// caution.
	//
	//   - ref. Private Header Parameter Names - JSON Web Encryption (JWE) https://www.rfc-editor.org/rfc/rfc7516#section-4.3
	PrivateHeaderParameters PrivateHeaderParameters `json:"-"`
}

type PrivateHeaderParameters map[string]any

type HeaderParameter func(h *Header)

func WithAlgorithm(alg string) HeaderParameter {
	return func(h *Header) {
		h.Algorithm = alg
	}
}

func WithEncryptionAlgorithm(enc string) HeaderParameter {
	return func(h *Header) {
		h.EncryptionAlgorithm = enc
	}
}

func WithCompressionAlgorithm(zip string) HeaderParameter {
	return func(h *Header) {
		h.CompressionAlgorithm = zip
	}
}

func WithJWKSetURL(jku string) HeaderParameter {
	return func(h *Header) {
		h.JWKSetURL = jku
	}
}

func WithJSONWebKey(jwk *jwk.JSONWebKey) HeaderParameter {
	return func(h *Header) {
		h.JSONWebKey = jwk
	}
}

func WithKeyID(kid string) HeaderParameter {
	return func(h *Header) {
		h.KeyID = kid
	}
}

func WithX509URL(x5u string) HeaderParameter {
	return func(h *Header) {
		h.X509URL = x5u
	}
}

func WithX509CertificateChain(x5c []string) HeaderParameter {
	return func(h *Header) {
		h.X509CertificateChain = x5c
	}
}

func WithX509CertificateSHA1Thumbprint(x5t string) HeaderParameter {
	return func(h *Header) {
		h.X509CertificateSHA1Thumbprint = x5t
	}
}

func WithX509CertificateSHA256Thumbprint(x5tS256 string) HeaderParameter {
	return func(h *Header) {
		h.X509CertificateSHA256Thumbprint = x5tS256
	}
}

func WithType(typ string) HeaderParameter {
	return func(h *Header) {
		h.Type = typ
	}
}

func WithContentType(cty string) HeaderParameter {
	return func(h *Header) {
		h.ContentType = cty
	}
}

func WithCritical(crit []string) HeaderParameter {
	return func(h *Header) {
		h.Critical = crit
	}
}

func WithPrivateHeaderParameter(name string, value any) HeaderParameter {
	return func(h *Header) {
		h.PrivateHeaderParameters[name] = value
	}
}

// NewHeader
//
// Example:
//
//	header := jose.NewHeader(
//		jwa.HS256,
//		jose.WithType("JWT"),
//	)
func NewHeader(alg string, parameters ...HeaderParameter) *Header {
	h := &Header{
		Algorithm:               alg,
		PrivateHeaderParameters: make(PrivateHeaderParameters),
	}

	for _, parameter := range parameters {
		parameter(h)
	}

	return h
}

var ErrInvalidJSON = errors.New("jose: invalid JSON")

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
		for i := range typ.NumField() {
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
		return nil, fmt.Errorf("invalid jose header: %+v: %w", _header, err)
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

func (h *Header) Decode(headerEncoded string) error {
	decoded, err := base64.RawURLEncoding.DecodeString(headerEncoded)
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}

	if err := json.Unmarshal(decoded, h); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

// GetPrivateHeaderParameter
//
//   - ref. https://pkg.go.dev/github.com/kunitsucom/util.go@v0.0.51/maps#Get
func (h *Header) GetPrivateHeaderParameter(parameterName string, v any) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Pointer && reflectValue.Kind() != reflect.Interface {
		return fmt.Errorf("v.(type)==%T: %w", v, ErrVIsNotPointerOrInterface)
	}
	reflectValueElem := reflectValue.Elem()
	param, ok := h.PrivateHeaderParameters[parameterName]
	if !ok {
		return fmt.Errorf("(*jose.Header).PrivateHeaderParameters[%q]: %w", parameterName, ErrPrivateHeaderParameterIsNotFound)
	}
	paramReflectValue := reflect.ValueOf(param)
	if reflectValueElem.Type() != paramReflectValue.Type() {
		return fmt.Errorf("(*jose.Header).PrivateHeaderParameters[%q].(type)==%T, v.(type)==%T: %w", parameterName, param, v, ErrPrivateHeaderParameterTypeIsNotMatch)
	}
	reflectValueElem.Set(paramReflectValue)
	return nil
}

func (h *Header) SetPrivateHeaderParameter(parameterName string, v any) {
	h.PrivateHeaderParameters[parameterName] = v
}
