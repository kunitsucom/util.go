package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kunitsuinc/util.go/pkg/net/http/urlcache"
)

const (
	Google ProviderMetadataURL = "https://accounts.google.com/.well-known/openid-configuration"

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
	urlcache *urlcache.Client[*ProviderMetadata]
}

func New(opts ...ClientOption) *Client {
	c := &Client{
		urlcache: urlcache.NewClient[*ProviderMetadata](http.DefaultClient),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type ClientOption func(*Client)

func WithURLCacheClient(client *urlcache.Client[*ProviderMetadata]) ClientOption {
	return func(d *Client) {
		d.urlcache = client
	}
}

func (d *Client) GetProviderMetadata(ctx context.Context, providerMetadataURL ProviderMetadataURL) (*ProviderMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, providerMetadataURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := d.urlcache.Do(req, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (*ProviderMetadata, error) {
		r := new(ProviderMetadata)
		if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
			return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.JWKSet): %w", err)
		}
		return r, nil
	})
	if err != nil {
		return nil, fmt.Errorf("(*urlcache.Client).Do: %w", err)
	}

	return resp, nil
}

//nolint:gochecknoglobals
var (
	Default = New()
)

func GetProviderMetadata(ctx context.Context, providerMetadataURL ProviderMetadataURL) (*ProviderMetadata, error) {
	return Default.GetProviderMetadata(ctx, providerMetadataURL)
}
