package id_token //nolint:revive,stylecheck

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kunitsuinc/util.go/pkg/jose"
	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
	"github.com/kunitsuinc/util.go/pkg/openid/discovery"
)

// NOTE: ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken
// NOTE: ref. http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#IDToken
type Claims struct {
	// iss
	// REQUIRED. Issuer Identifier for the Issuer of the response. The iss value is a case sensitive URL using the https scheme that contains scheme, host, and optionally, port number and path components and no query or fragment components.
	Issuer string `json:"iss"`
	// sub
	// REQUIRED. Subject Identifier. A locally unique and never reassigned identifier within the Issuer for the End-User, which is intended to be consumed by the Client, e.g., 24400320 or AItOawmwtWwcT0k51BayewNvutrJUqsvl6qs7A4. It MUST NOT exceed 255 ASCII characters in length. The sub value is a case sensitive string.
	SubjectIdentifier string `json:"sub"` // set client identifier. len <= 255.
	// aud
	// REQUIRED. Audience(s) that this ID Token is intended for. It MUST contain the OAuth 2.0 client_id of the Relying Party as an audience value. It MAY also contain identifiers for other audiences. In the general case, the aud value is an array of case sensitive strings. In the common special case when there is one audience, the aud value MAY be a single case sensitive string.
	Audience string `json:"aud"`
	// exp
	// REQUIRED. Expiration time on or after which the ID Token MUST NOT be accepted for processing. The processing of this parameter requires that the current date/time MUST be before the expiration date/time listed in the value. Implementers MAY provide for some small leeway, usually no more than a few minutes, to account for clock skew. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. See RFC 3339Klyne, G., Ed. and C. Newman, “Date and Time on the Internet: Timestamps,” July 2002. [RFC3339] for details regarding date/times in general and UTC in particular.
	ExpirationTime int64 `json:"exp"`
	// iat
	// REQUIRED. Time at which the JWT was issued. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time.
	IssuedAt int64 `json:"iat"`
	// auth_time
	// Time when the End-User authentication occurred. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. When a max_age request is made or when auth_time is requested as an Essential Claim, then this Claim is REQUIRED; otherwise, its inclusion is OPTIONAL. (The auth_time Claim semantically corresponds to the OpenID 2.0 PAPE [OpenID.PAPE] auth_time response parameter.)
	AuthenticationTime int64 `json:"auth_time"` //nolint:tagliatelle
	// nonce
	// String value used to associate a Client session with an ID Token, and to mitigate replay attacks. The value is passed through unmodified from the Authentication Request to the ID Token. If present in the ID Token, Clients MUST verify that the nonce Claim Value is equal to the value of the nonce parameter sent in the Authentication Request. If present in the Authentication Request, Authorization Servers MUST include a nonce Claim in the ID Token with the Claim Value being the nonce value sent in the Authentication Request. Authorization Servers SHOULD perform no other processing on nonce values used. The nonce value is a case sensitive string.
	Nonce string `json:"nonce,omitempty"`
	// acr
	// OPTIONAL. Authentication Context Class Reference. String specifying an Authentication Context Class Reference value that identifies the Authentication Context Class that the authentication performed satisfied. The value "0" indicates the End-User authentication did not meet the requirements of ISO/IEC 29115 [ISO29115] level 1. Authentication using a long-lived browser cookie, for instance, is one example where the use of "level 0" is appropriate. Authentications with level 0 SHOULD NOT be used to authorize access to any resource of any monetary value. (This corresponds to the OpenID 2.0 PAPE [OpenID.PAPE] nist_auth_level 0.) An absolute URI or an RFC 6711 [RFC6711] registered name SHOULD be used as the acr value; registered names MUST NOT be used with a different meaning than that which is registered. Parties using this claim will need to agree upon the meanings of the values used, which may be context-specific. The acr value is a case sensitive string.
	AuthenticationContextClassReference string `json:"acr,omitempty"`
	// amr
	// OPTIONAL. Authentication Methods References. JSON array of strings that are identifiers for authentication methods used in the authentication. For instance, values might indicate that both password and OTP authentication methods were used. The definition of particular values to be used in the amr Claim is beyond the scope of this specification. Parties using this claim will need to agree upon the meanings of the values used, which may be context-specific. The amr value is an array of case sensitive strings.
	AuthenticationMethodsReferences string `json:"amr,omitempty"`
	// azp
	// OPTIONAL. Authorized party - the party to which the ID Token was issued. If present, it MUST contain the OAuth 2.0 Client ID of this party. This Claim is only needed when the ID Token has a single audience value and that audience is different than the authorized party. It MAY be included even when the authorized party is the same as the sole audience. The azp value is a case sensitive string containing a StringOrURI value.
	AuthorizedParty string `json:"azp,omitempty"`

	// https://datatracker.ietf.org/doc/html/rfc7519#section-4.3
	PrivateClaims map[string]string `json:"-"`
}

type Client struct {
	client *jwk.Client
}

type ClientOption func(c *Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

var ErrNotMatchFormat = errors.New(`id_token: not match format`)

// NOTE: ref. http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#IDTokenValidation
func (c *Client) Verify(id_token string) error { //nolint:revive,stylecheck
	s := strings.Split(id_token, ".")
	if len(s) != 3 {
		return ErrNotMatchFormat
	}
	headerBase64, payloadBase64, signatureBase64 := s[0], s[1], s[2]
	fmt.Printf("%s: [DEBUG]: %s %s %s\n", time.Now(), headerBase64, payloadBase64, signatureBase64) //nolint:forbidigo

	headerBytes, err := base64.RawURLEncoding.DecodeString(headerBase64)
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}
	header := new(jose.Header)
	if err := json.NewDecoder(bytes.NewReader(headerBytes)).Decode(header); err != nil {
		return fmt.Errorf("(*json.Decoder).Decode: %w", err)
	}
	fmt.Printf("%s: [DEBUG]: (*json.Decoder).Decode: %#v\n", time.Now(), header) //nolint:forbidigo

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}
	payload := new(Claims)
	if err := json.NewDecoder(bytes.NewReader(payloadBytes)).Decode(payload); err != nil {
		return fmt.Errorf("(*json.Decoder).Decode: %w", err)
	}
	fmt.Printf("%s: [DEBUG]: (*json.Decoder).Decode: %#v\n", time.Now(), payload) //nolint:forbidigo

	discoveryURI := payload.Issuer + discovery.ProviderMetadataURLPath
	fmt.Printf("%s: [DEBUG]: discoveryURI: %#v\n", time.Now(), discoveryURI) //nolint:forbidigo

	doc, err := discovery.GetProviderMetadata(context.Background(), discoveryURI)
	if err != nil {
		return fmt.Errorf("%s: discovery.GetDocument: %w", time.Now(), err)
	}

	jwks, err := c.client.GetJWKSet(context.Background(), doc.JwksURI)
	if err != nil {
		return fmt.Errorf("%s: discovery.GetDocument: %w", time.Now(), err)
	}
	fmt.Printf("%s: [DEBUG]: (*json.Decoder).Decode: %#v\n", time.Now(), jwks) //nolint:forbidigo

	key, err := jwks.GetJSONWebKey(header.KeyID)
	if err != nil {
		return fmt.Errorf("jwk.GetJSONWebKey: %w", err)
	}
	fmt.Printf("%s: [DEBUG]: (*json.Decoder).Decode: %#v\n", time.Now(), key) //nolint:forbidigo

	return nil
}

//nolint:gochecknoglobals
var (
	Default = NewClient()
)

func Verify(id_token string) error { //nolint:revive,stylecheck
	return Default.Verify(id_token)
}
