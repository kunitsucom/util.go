package authentication_request //nolint:revive,stylecheck

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kunitsuinc/util.go/pkg/openid/pkce"
)

//nolint:revive,stylecheck
type option struct {
	responseMode string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=%5BOAuth.Responses%5D%3A-,response_mode,-OPTIONAL.%20Informs%20the
	nonce        string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=following%20request%20parameters%3A-,nonce,-OPTIONAL.%20String%20value
	display      string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=Section%C2%A015.5.2.-,display,-OPTIONAL.%20ASCII%20string
	prompt       string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=an%20appropriate%20display.-,prompt,-OPTIONAL.%20Space%20delimited
	maxAge       int64    // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=error%20is%20returned.-,max_age,-OPTIONAL.%20Maximum%20Authentication
	uiLocales    []string // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=auth_time%20Claim%20Value.-,ui_locales,-OPTIONAL.%20End%2DUser%27s
	idTokenHint  string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OpenID%20Provider.-,id_token_hint,-OPTIONAL.%20ID%20Token
	loginHint    string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20id_token_hint%20value.-,login_hint,-OPTIONAL.%20Hint%20to
	acrValues    []string // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OP%27s%20discretion.-,acr_values,-OPTIONAL.%20Requested%20Authentication

	codeChallenge       string // ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=following%20additional%20parameters%3A-,code_challenge,-REQUIRED.%20%20Code%20challenge
	codeChallengeMethod string // ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=REQUIRED.%20%20Code%20challenge.-,code_challenge_method,-OPTIONAL

	accessType string // ref. https://developers.google.com/identity/protocols/oauth2/web-server#request-parameter-access_type
}

type Option func(*option)

var ErrParameterIsEmpty = errors.New("openid: Authentication Request: parameter is empty")

// WithResponseMode
//
//   - response_mode: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=%5BOAuth.Responses%5D%3A-,response_mode,-OPTIONAL.%20Informs%20the
func WithResponseMode(responseMode string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.responseMode = responseMode }
}

// WithNonce
//
//   - nonce: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=following%20request%20parameters%3A-,nonce,-OPTIONAL.%20String%20value
func WithNonce(nonce string) Option {
	return func(o *option) { o.nonce = nonce }
}

// WithDisplay
//
//   - display: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=Section%C2%A015.5.2.-,display,-OPTIONAL.%20ASCII%20string
func WithDisplay(display string) Option {
	return func(o *option) { o.display = display }
}

// WithPrompt
//
//   - prompt: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=an%20appropriate%20display.-,prompt,-OPTIONAL.%20Space%20delimited
func WithPrompt(prompt string) Option {
	return func(o *option) { o.prompt = prompt }
}

// WithMaxAge
//
//   - max_age: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=error%20is%20returned.-,max_age,-OPTIONAL.%20Maximum%20Authentication
func WithMaxAge(maxAge int64) Option { //nolint:revive,stylecheck
	return func(o *option) { o.maxAge = maxAge }
}

// WithUILocales
//
//   - ui_locales: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=auth_time%20Claim%20Value.-,ui_locales,-OPTIONAL.%20End%2DUser%27s
func WithUILocales(uiLocales []string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.uiLocales = uiLocales }
}

// WithIDTokenHint
//
//   - id_token_hint: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OpenID%20Provider.-,id_token_hint,-OPTIONAL.%20ID%20Token
func WithIDTokenHint(idTokenHint string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.idTokenHint = idTokenHint }
}

// WithLoginHint
//
// - login_hint: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20id_token_hint%20value.-,login_hint,-OPTIONAL.%20Hint%20to
func WithLoginHint(loginHint string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.loginHint = loginHint }
}

// WithACRValues
//
//   - acr_values: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OP%27s%20discretion.-,acr_values,-OPTIONAL.%20Requested%20Authentication
func WithACRValues(acrValues []string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.acrValues = acrValues }
}

// WithCodeChallengeForPKCE
//
//   - code_challenge:
//     REQUIRED.  Code challenge.
//     ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=following%20additional%20parameters%3A-,code_challenge,-REQUIRED.%20%20Code%20challenge
//   - code_challenge_method:
//     OPTIONAL, defaults to "plain" if not present in the request. Code verifier transformation method is "S256" or "plain".
//     ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=REQUIRED.%20%20Code%20challenge.-,code_challenge_method,-OPTIONAL
func WithCodeChallengeForPKCE(codeVerifier pkce.CodeVerifier, codeChallengeMethod pkce.CodeChallengeMethod) Option { //nolint:revive,stylecheck
	return func(o *option) {
		o.codeChallenge = codeVerifier.Encode(codeChallengeMethod)
		o.codeChallengeMethod = codeChallengeMethod.String()
	}
}

// WithAccessType
//
//   - access_type: ref. https://developers.google.com/identity/protocols/oauth2/web-server#request-parameter-access_type
func WithAccessType(accessType string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.accessType = accessType }
}

// New returns URL for redirect response to User-Agent.
//
//   - authorization_endpoint: ref. https://openid.net/specs/openid-connect-discovery-1_0.html#:~:text=from%20this%20Issuer.-,authorization_endpoint,-REQUIRED.%20URL%20of
//   - scope:                  ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=Authorization%20Code%20Flow%3A-,scope,-REQUIRED.%20OpenID%20Connect
//   - response_type:          ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=by%20this%20specification.-,response_type,-REQUIRED.%20OAuth%202.0
//   - client_id:              ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=value%20is%20code.-,client_id,-REQUIRED.%20OAuth%202.0%20Client
//   - redirect_uri:           ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20Authorization%20Server.-,redirect_uri,-REQUIRED.%20Redirection%20URI
//   - state:                  ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=a%20native%20application.-,state,-RECOMMENDED.%20Opaque%20value
func New(authorizationEndpoint string, scope []string, responseType, clientID, redirectURI, state string, opts ...Option) (*url.URL, error) { //nolint:revive,stylecheck
	o := new(option)
	for _, opt := range opts {
		opt(o)
	}

	u, err := url.Parse(authorizationEndpoint)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	if len(scope) == 0 {
		return nil, fmt.Errorf("scope: %w", ErrParameterIsEmpty)
	}

	switch "" {
	case responseType:
		return nil, fmt.Errorf("response_type: %w", ErrParameterIsEmpty)
	case clientID:
		return nil, fmt.Errorf("client_id: %w", ErrParameterIsEmpty)
	case redirectURI:
		return nil, fmt.Errorf("redirect_uri: %w", ErrParameterIsEmpty)
	case state:
		return nil, fmt.Errorf("state: %w", ErrParameterIsEmpty)
	}

	query := make(url.Values)
	query.Add("scope", strings.Join(scope, " "))
	query.Add("response_type", responseType)
	query.Add("client_id", clientID)
	query.Add("redirect_uri", redirectURI)
	query.Add("state", state)

	query = optionalParameters(query, o)

	u.RawQuery = strings.ReplaceAll(query.Encode(), "+", "%20")

	return u, nil
}

func optionalParameters(query url.Values, optionals *option) url.Values { //nolint:cyclop
	if optionals.responseMode != "" {
		query.Add("response_mode", optionals.responseMode)
	}
	if optionals.nonce != "" {
		query.Add("nonce", optionals.nonce)
	}
	if optionals.display != "" {
		query.Add("display", optionals.display)
	}
	if optionals.prompt != "" {
		query.Add("prompt", optionals.prompt)
	}
	if optionals.maxAge != 0 {
		query.Add("max_age", strconv.FormatInt(optionals.maxAge, 10))
	}
	if len(optionals.uiLocales) != 0 {
		query.Add("ui_locales", strings.Join(optionals.uiLocales, " "))
	}
	if optionals.idTokenHint != "" {
		query.Add("id_token_hint", optionals.idTokenHint)
	}
	if optionals.loginHint != "" {
		query.Add("login_hint", optionals.loginHint)
	}
	if len(optionals.acrValues) != 0 {
		query.Add("acr_values", strings.Join(optionals.acrValues, " "))
	}
	if optionals.codeChallenge != "" && optionals.codeChallengeMethod != "" {
		query.Add("code_challenge", optionals.codeChallenge)
		query.Add("code_challenge_method", optionals.codeChallengeMethod)
	}
	if optionals.accessType != "" {
		query.Add("access_type", optionals.accessType)
	}

	return query
}
