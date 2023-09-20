package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	slicez "github.com/kunitsucom/util.go/slices"
	syncz "github.com/kunitsucom/util.go/sync"
)

const (
	Apple         ProviderMetadataURL = "https://appleid.apple.com/.well-known/openid-configuration"
	GitHubActions ProviderMetadataURL = "https://token.actions.githubusercontent.com/.well-known/openid-configuration"
	Google        ProviderMetadataURL = "https://accounts.google.com/.well-known/openid-configuration"
	LINE          ProviderMetadataURL = "https://access.line.me/.well-known/openid-configuration"
	Microsoft     ProviderMetadataURL = "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration"
	Slack         ProviderMetadataURL = "https://slack.com/.well-known/openid-configuration"
	Yahoo         ProviderMetadataURL = "https://api.login.yahoo.com/.well-known/openid-configuration"
	YahooJapan    ProviderMetadataURL = "https://auth.login.yahoo.co.jp/yconnect/v2/.well-known/openid-configuration"

	ProviderMetadataURLPath = "/.well-known/openid-configuration"
)

type (
	ProviderMetadataURL = string

	// https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
	//
	//nolint:tagliatelle
	ProviderMetadata struct { //nolint:revive
		Issuer                                     string   `json:"issuer"`
		AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
		TokenEndpoint                              string   `json:"token_endpoint,omitempty"`
		UserInfoEndpoint                           string   `json:"userinfo_endpoint,omitempty"`
		JwksURI                                    string   `json:"jwks_uri"`
		RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
		ScopesSupported                            []string `json:"scopes_supported,omitempty"`
		ResponseTypesSupported                     []string `json:"response_types_supported"`
		ResponseModesSupported                     []string `json:"response_modes_supported,omitempty"`
		GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
		ACRValuesSupported                         []string `json:"acr_values_supported,omitempty"`
		SubjectTypesSupported                      []string `json:"subject_types_supported"`
		IDTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
		IDTokenEncryptionAlgValuesSupported        []string `json:"id_token_encryption_alg_values_supported,omitempty"`
		IDTokenEncryptionEncValuesSupported        []string `json:"id_token_encryption_enc_values_supported,omitempty"`
		UserinfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported,omitempty"`
		UserinfoEncryptionAlgValuesSupported       []string `json:"userinfo_encryption_alg_values_supported,omitempty"`
		UserinfoEncryptionEncValuesSupported       []string `json:"userinfo_encryption_enc_values_supported,omitempty"`
		RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported,omitempty"`
		RequestObjectEncryptionAlgValuesSupported  []string `json:"request_object_encryption_alg_values_supported,omitempty"`
		RequestObjectEncryptionEncValuesSupported  []string `json:"request_object_encryption_enc_values_supported,omitempty"`
		TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
		TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
		DisplayValuesSupported                     []string `json:"display_values_supported,omitempty"`
		ClaimTypesSupported                        []string `json:"claim_types_supported,omitempty"`
		ClaimsSupported                            []string `json:"claims_supported,omitempty"`
		ServiceDocumentation                       string   `json:"service_documentation,omitempty"`
		ClaimsLocalesSupported                     []string `json:"claims_locales_supported,omitempty"`
		UILocalesSupported                         []string `json:"ui_locales_supported,omitempty"`
		ClaimsParameterSupported                   bool     `json:"claims_parameter_supported,omitempty"`
		RequestParameterSupported                  bool     `json:"request_parameter_supported,omitempty"`
		RequestURIParameterSupported               bool     `json:"request_uri_parameter_supported,omitempty"`
		RequireRequestURIRegistration              bool     `json:"require_request_uri_registration,omitempty"`
		OPPolicyURI                                string   `json:"op_policy_uri,omitempty"`
		OPTosURI                                   string   `json:"op_tos_uri,omitempty"`
	}
)

type Client struct {
	client   *http.Client
	cacheMap syncz.Map[*ProviderMetadata]
}

func New(ctx context.Context, opts ...ClientOption) *Client {
	c := &Client{
		client:   http.DefaultClient,
		cacheMap: syncz.NewMap[*ProviderMetadata](ctx, syncz.WithNewMapOptionCleanerInterval(10*time.Minute)),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(d *Client) {
		d.client = client
	}
}

func WithCacheMap(cacheMap syncz.Map[*ProviderMetadata]) ClientOption {
	return func(d *Client) {
		d.cacheMap = cacheMap
	}
}

var ErrResponseIsNotCacheable = errors.New("discovery: response is not cacheable")

func (d *Client) GetProviderMetadata(ctx context.Context, providerMetadataURL ProviderMetadataURL) (*ProviderMetadata, error) {
	if cached, ok := d.cacheMap.Load(providerMetadataURL); ok {
		return cached, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, providerMetadataURL, nil)
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
		bodyCutOff := slicez.CutOff(body, 100)
		return nil, fmt.Errorf("code=%d body=%q: %w", resp.StatusCode, string(bodyCutOff), ErrResponseIsNotCacheable)
	}

	r := new(ProviderMetadata)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.JWKSet): %w", err)
	}

	d.cacheMap.Store(providerMetadataURL, r)
	return r, nil
}

//nolint:gochecknoglobals
var (
	Default = New(context.Background())
)

func GetProviderMetadata(ctx context.Context, providerMetadataURL ProviderMetadataURL) (*ProviderMetadata, error) {
	return Default.GetProviderMetadata(ctx, providerMetadataURL)
}
