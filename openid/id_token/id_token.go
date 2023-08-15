package id_token //nolint:revive,stylecheck

// NOTE: ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken
// NOTE: ref. http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#IDToken
type Claims struct {
	// Issuer: "iss"
	//
	// REQUIRED. Issuer Identifier for the Issuer of the response. The iss value is a case sensitive URL using the https scheme that contains scheme, host, and optionally, port number and path components and no query or fragment components.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=by%20OpenID%20Connect%3A-,iss,-REQUIRED.%20Issuer%20Identifier
	Issuer string `json:"iss"`
	// SubjectIdentifier: "sub"
	//
	// REQUIRED. Subject Identifier. A locally unique and never reassigned identifier within the Issuer for the End-User, which is intended to be consumed by the Client, e.g., 24400320 or AItOawmwtWwcT0k51BayewNvutrJUqsvl6qs7A4. It MUST NOT exceed 255 ASCII characters in length. The sub value is a case sensitive string.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=or%20fragment%20components.-,sub,-REQUIRED.%20Subject%20Identifier
	SubjectIdentifier string `json:"sub"` // set client identifier. len <= 255.
	// Audience: "aud"
	//
	// REQUIRED. Audience(s) that this ID Token is intended for. It MUST contain the OAuth 2.0 client_id of the Relying Party as an audience value. It MAY also contain identifiers for other audiences. In the general case, the aud value is an array of case sensitive strings. In the common special case when there is one audience, the aud value MAY be a single case sensitive string.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=case%20sensitive%20string.-,aud,-REQUIRED.%20Audience(s
	Audience string `json:"aud"`
	// ExpirationTime: "exp"
	//
	// REQUIRED. Expiration time on or after which the ID Token MUST NOT be accepted for processing. The processing of this parameter requires that the current date/time MUST be before the expiration date/time listed in the value. Implementers MAY provide for some small leeway, usually no more than a few minutes, to account for clock skew. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. See RFC 3339Klyne, G., Ed. and C. Newman, “Date and Time on the Internet: Timestamps,” July 2002. [RFC3339] for details regarding date/times in general and UTC in particular.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=case%20sensitive%20string.-,exp,-REQUIRED.%20Expiration%20time
	ExpirationTime int64 `json:"exp"`
	// IssuedAt: "iat"
	//
	// REQUIRED. Time at which the JWT was issued. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=UTC%20in%20particular.-,iat,-REQUIRED.%20Time%20at
	IssuedAt int64 `json:"iat"`
	// AuthenticationTime: "auth_time"
	//
	// Time when the End-User authentication occurred. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. When a max_age request is made or when auth_time is requested as an Essential Claim, then this Claim is REQUIRED; otherwise, its inclusion is OPTIONAL. (The auth_time Claim semantically corresponds to the OpenID 2.0 PAPE [OpenID.PAPE] auth_time response parameter.)
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=the%20date/time.-,auth_time,-Time%20when%20the
	AuthenticationTime int64 `json:"auth_time,omitempty"` //nolint:tagliatelle
	// Nonce: "nonce"
	//
	// String value used to associate a Client session with an ID Token, and to mitigate replay attacks. The value is passed through unmodified from the Authentication Request to the ID Token. If present in the ID Token, Clients MUST verify that the nonce Claim Value is equal to the value of the nonce parameter sent in the Authentication Request. If present in the Authentication Request, Authorization Servers MUST include a nonce Claim in the ID Token with the Claim Value being the nonce value sent in the Authentication Request. Authorization Servers SHOULD perform no other processing on nonce values used. The nonce value is a case sensitive string.
	//
	//   - ref. hthttps://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=response%20parameter.%29-,nonce,-String%20value%20used
	Nonce string `json:"nonce,omitempty"`
	// AuthenticationContextClassReference: "acr"
	//
	// OPTIONAL. Authentication Context Class Reference. String specifying an Authentication Context Class Reference value that identifies the Authentication Context Class that the authentication performed satisfied. The value "0" indicates the End-User authentication did not meet the requirements of ISO/IEC 29115 [ISO29115] level 1. Authentication using a long-lived browser cookie, for instance, is one example where the use of "level 0" is appropriate. Authentications with level 0 SHOULD NOT be used to authorize access to any resource of any monetary value. (This corresponds to the OpenID 2.0 PAPE [OpenID.PAPE] nist_auth_level 0.) An absolute URI or an RFC 6711 [RFC6711] registered name SHOULD be used as the acr value; registered names MUST NOT be used with a different meaning than that which is registered. Parties using this claim will need to agree upon the meanings of the values used, which may be context-specific. The acr value is a case sensitive string.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=case%20sensitive%20string.-,acr,-OPTIONAL.%20Authentication%20Context
	AuthenticationContextClassReference string `json:"acr,omitempty"`
	// AuthenticationMethodsReferences: "amr"
	//
	// OPTIONAL. Authentication Methods References. JSON array of strings that are identifiers for authentication methods used in the authentication. For instance, values might indicate that both password and OTP authentication methods were used. The definition of particular values to be used in the amr Claim is beyond the scope of this specification. Parties using this claim will need to agree upon the meanings of the values used, which may be context-specific. The amr value is an array of case sensitive strings.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=case%20sensitive%20strings.-,azp,-OPTIONAL.%20Authorized%20party
	AuthenticationMethodsReferences string `json:"amr,omitempty"`
	// AuthorizedParty: "azp"
	//
	// OPTIONAL. Authorized party - the party to which the ID Token was issued. If present, it MUST contain the OAuth 2.0 Client ID of this party. This Claim is only needed when the ID Token has a single audience value and that audience is different than the authorized party. It MAY be included even when the authorized party is the same as the sole audience. The azp value is a case sensitive string containing a StringOrURI value.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#IDToken:~:text=or%20fragment%20components.-,sub,-REQUIRED.%20Subject%20Identifier
	AuthorizedParty string `json:"azp,omitempty"`

	// AccessTokenHash: "at_hash"
	//
	// OPTIONAL. Access Token hash value. Its value is the base64url encoding of the left-most half of the hash of the octets of the ASCII representation of the access_token value, where the hash algorithm used is the hash algorithm used in the alg Header Parameter of the ID Token's JOSE Header. For instance, if the alg is RS256, hash the access_token value with SHA-256, then take the left-most 128 bits and base64url encode them. The at_hash value is a case sensitive string.
	//
	//   - ref. https://openid.net/specs/openid-connect-core-1_0.html#CodeIDToken
	AccessTokenHash string `json:"at_hash,omitempty"` //nolint:tagliatelle
}

// TODO: impl
