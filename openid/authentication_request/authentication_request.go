package authentication_request //nolint:revive,stylecheck

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kunitsuinc/util.go/openid/pkce"
)

//nolint:revive,stylecheck
type option struct {
	response_mode string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=%5BOAuth.Responses%5D%3A-,response_mode,-OPTIONAL.%20Informs%20the
	nonce         string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=following%20request%20parameters%3A-,nonce,-OPTIONAL.%20String%20value
	display       string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=Section%C2%A015.5.2.-,display,-OPTIONAL.%20ASCII%20string
	prompt        string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=an%20appropriate%20display.-,prompt,-OPTIONAL.%20Space%20delimited
	max_age       int64    // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=error%20is%20returned.-,max_age,-OPTIONAL.%20Maximum%20Authentication
	ui_locales    []string // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=auth_time%20Claim%20Value.-,ui_locales,-OPTIONAL.%20End%2DUser%27s
	id_token_hint string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OpenID%20Provider.-,id_token_hint,-OPTIONAL.%20ID%20Token
	login_hint    string   // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20id_token_hint%20value.-,login_hint,-OPTIONAL.%20Hint%20to
	acr_values    []string // ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OP%27s%20discretion.-,acr_values,-OPTIONAL.%20Requested%20Authentication

	code_challenge        string // ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=following%20additional%20parameters%3A-,code_challenge,-REQUIRED.%20%20Code%20challenge
	code_challenge_method string // ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=REQUIRED.%20%20Code%20challenge.-,code_challenge_method,-OPTIONAL

	access_type string // ref. https://developers.google.com/identity/protocols/oauth2/web-server#request-parameter-access_type
}

type Option func(*option)

var ErrParameterIsEmpty = errors.New("openid: Authentication Request: parameter is empty")

// WithResponseMode
//
//   - response_mode: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=%5BOAuth.Responses%5D%3A-,response_mode,-OPTIONAL.%20Informs%20the
func WithResponseMode(response_mode string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.response_mode = response_mode }
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
func WithMaxAge(max_age int64) Option { //nolint:revive,stylecheck
	return func(o *option) { o.max_age = max_age }
}

// WithUILocales
//
//   - ui_locales: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=auth_time%20Claim%20Value.-,ui_locales,-OPTIONAL.%20End%2DUser%27s
func WithUILocales(ui_locales []string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.ui_locales = ui_locales }
}

// WithIDTokenHint
//
//   - id_token_hint: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OpenID%20Provider.-,id_token_hint,-OPTIONAL.%20ID%20Token
func WithIDTokenHint(id_token_hint string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.id_token_hint = id_token_hint }
}

// WithLoginHint
//
// - login_hint: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20id_token_hint%20value.-,login_hint,-OPTIONAL.%20Hint%20to
func WithLoginHint(login_hint string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.login_hint = login_hint }
}

// WithACRValues
//
//   - acr_values: ref: https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20OP%27s%20discretion.-,acr_values,-OPTIONAL.%20Requested%20Authentication
func WithACRValues(acr_values []string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.acr_values = acr_values }
}

// WithCodeChallengeForPKCE
//
//   - code_challenge:
//     REQUIRED.  Code challenge.
//     ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=following%20additional%20parameters%3A-,code_challenge,-REQUIRED.%20%20Code%20challenge
//   - code_challenge_method:
//     OPTIONAL, defaults to "plain" if not present in the request. Code verifier transformation method is "S256" or "plain".
//     ref. https://www.rfc-editor.org/rfc/rfc7636#:~:text=REQUIRED.%20%20Code%20challenge.-,code_challenge_method,-OPTIONAL
func WithCodeChallengeForPKCE(code_verifier pkce.CodeVerifier, code_challenge_method pkce.CodeChallengeMethod) Option { //nolint:revive,stylecheck
	return func(o *option) {
		o.code_challenge = code_verifier.Encode(code_challenge_method)
		o.code_challenge_method = code_challenge_method.String()
	}
}

// WithAccessType
//
//   - access_type: ref. https://developers.google.com/identity/protocols/oauth2/web-server#request-parameter-access_type
func WithAccessType(access_type string) Option { //nolint:revive,stylecheck
	return func(o *option) { o.access_type = access_type }
}

// New returns URL for redirect response to User-Agent.
//
//   - authorization_endpoint: ref. https://openid.net/specs/openid-connect-discovery-1_0.html#:~:text=from%20this%20Issuer.-,authorization_endpoint,-REQUIRED.%20URL%20of
//   - scope:                  ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=Authorization%20Code%20Flow%3A-,scope,-REQUIRED.%20OpenID%20Connect
//   - response_type:          ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=by%20this%20specification.-,response_type,-REQUIRED.%20OAuth%202.0
//   - client_id:              ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=value%20is%20code.-,client_id,-REQUIRED.%20OAuth%202.0%20Client
//   - redirect_uri:           ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=the%20Authorization%20Server.-,redirect_uri,-REQUIRED.%20Redirection%20URI
//   - state:                  ref. https://openid.net/specs/openid-connect-core-1_0.html#:~:text=a%20native%20application.-,state,-RECOMMENDED.%20Opaque%20value
func New(authorization_endpoint string, scope []string, response_type, client_id, redirect_uri, state string, opts ...Option) (*url.URL, error) { //nolint:revive,stylecheck
	o := new(option)
	for _, opt := range opts {
		opt(o)
	}

	u, err := url.Parse(authorization_endpoint)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	if len(scope) == 0 {
		return nil, fmt.Errorf("scope=%v: %w", scope, ErrParameterIsEmpty)
	}

	if response_type == "" || client_id == "" || redirect_uri == "" || state == "" {
		return nil, fmt.Errorf("response_type=%s, client_id=%s, redirect_uri=%s, state=%s: %w", response_type, client_id, redirect_uri, state, ErrParameterIsEmpty)
	}

	query := make(url.Values)
	query.Add("scope", strings.Join(scope, " "))
	query.Add("response_type", response_type)
	query.Add("client_id", client_id)
	query.Add("redirect_uri", redirect_uri)
	query.Add("state", state)

	query = optionalParameters(query, o)

	u.RawQuery = strings.ReplaceAll(query.Encode(), "+", "%20")

	return u, nil
}

func optionalParameters(query url.Values, optionals *option) url.Values { //nolint:cyclop
	if optionals.response_mode != "" {
		query.Add("response_mode", optionals.response_mode)
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
	if optionals.max_age != 0 {
		query.Add("max_age", strconv.FormatInt(optionals.max_age, 10))
	}
	if len(optionals.ui_locales) != 0 {
		query.Add("ui_locales", strings.Join(optionals.ui_locales, " "))
	}
	if optionals.id_token_hint != "" {
		query.Add("id_token_hint", optionals.id_token_hint)
	}
	if optionals.login_hint != "" {
		query.Add("login_hint", optionals.login_hint)
	}
	if len(optionals.acr_values) != 0 {
		query.Add("acr_values", strings.Join(optionals.acr_values, " "))
	}
	if optionals.code_challenge != "" && optionals.code_challenge_method != "" {
		query.Add("code_challenge", optionals.code_challenge)
		query.Add("code_challenge_method", optionals.code_challenge_method)
	}
	if optionals.access_type != "" {
		query.Add("access_type", optionals.access_type)
	}

	return query
}
