package jwa

//go:generate stringer -linecomment -type Algorithm

// Algorithm
//
// 3.1.  "alg" (Algorithm) Header Parameter Values for JWS
//
//	The table below is the set of "alg" (algorithm) Header Parameter
//	values defined by this specification for use with JWS, each of which
//	is explained in more detail in the following sections:
//
//	+--------------+-------------------------------+--------------------+
//	| "alg" Param  | Digital Signature or MAC      | Implementation     |
//	| Value        | Algorithm                     | Requirements       |
//	+--------------+-------------------------------+--------------------+
//	| HS256        | HMAC using SHA-256            | Required           |
//	| HS384        | HMAC using SHA-384            | Optional           |
//	| HS512        | HMAC using SHA-512            | Optional           |
//	| RS256        | RSASSA-PKCS1-v1_5 using       | Recommended        |
//	|              | SHA-256                       |                    |
//	| RS384        | RSASSA-PKCS1-v1_5 using       | Optional           |
//	|              | SHA-384                       |                    |
//	| RS512        | RSASSA-PKCS1-v1_5 using       | Optional           |
//	|              | SHA-512                       |                    |
//	| ES256        | ECDSA using P-256 and SHA-256 | Recommended+       |
//	| ES384        | ECDSA using P-384 and SHA-384 | Optional           |
//	| ES512        | ECDSA using P-521 and SHA-512 | Optional           |
//	| PS256        | RSASSA-PSS using SHA-256 and  | Optional           |
//	|              | MGF1 with SHA-256             |                    |
//	| PS384        | RSASSA-PSS using SHA-384 and  | Optional           |
//	|              | MGF1 with SHA-384             |                    |
//	| PS512        | RSASSA-PSS using SHA-512 and  | Optional           |
//	|              | MGF1 with SHA-512             |                    |
//	| none         | No digital signature or MAC   | Optional           |
//	|              | performed                     |                    |
//	+--------------+-------------------------------+--------------------+
//
//	The use of "+" in the Implementation Requirements column indicates
//	that the requirement strength is likely to be increased in a future
//	version of the specification.
//
//	See Appendix A.1 for a table cross-referencing the JWS digital
//	signature and MAC "alg" (algorithm) values defined in this
//	specification with the equivalent identifiers used by other standards
//	and software packages.
//
// - ref. https://www.rfc-editor.org/rfc/rfc7518#section-3.1
type Algorithm int

const (
	HS256 Algorithm = iota
	HS384
	HS512
	RS256
	RS384
	RS512
	ES256
	ES384
	ES512
	PS256
	PS384
	PS512
	None // none
)
