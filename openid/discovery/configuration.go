package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ProviderConfigurationURL string

func (pcf ProviderConfigurationURL) String() string { return string(pcf) }

const (
	Google ProviderConfigurationURL = "https://accounts.google.com/.well-known/openid-configuration"
)

const ProviderConfigurationPath = "/.well-known/openid-configuration"

// https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
//
//nolint:tagliatelle
type ProviderConfigurationResponse struct { //nolint:revive
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

type ProviderConfigurationCache struct {
	*ProviderConfigurationResponse
	ExpirationTime time.Time
}

func (d ProviderConfigurationCache) Expired() bool {
	return d._Expired(time.Now())
}

func (d ProviderConfigurationCache) _Expired(now time.Time) bool {
	return d.ExpirationTime.Before(now)
}

type Discovery struct {
	doNotUseCache bool
	cacheTTL      time.Duration
	cacheMap      map[ProviderConfigurationURL]ProviderConfigurationCache
	cacheSync     sync.Mutex
}

func New(opts ...NewDiscoveryOption) *Discovery {
	d := &Discovery{
		doNotUseCache: false,
		cacheTTL:      2 * time.Minute,
		cacheMap:      map[ProviderConfigurationURL]ProviderConfigurationCache{},
		cacheSync:     sync.Mutex{},
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

type NewDiscoveryOption func(*Discovery)

func WithCache(useCache bool) NewDiscoveryOption {
	return func(d *Discovery) {
		d.doNotUseCache = !useCache
	}
}

func WithCacheTTL(ttl time.Duration) NewDiscoveryOption {
	return func(d *Discovery) {
		d.cacheTTL = ttl
	}
}

type GetProviderConfigurationOption func(*getProviderConfigurationOption)

type getProviderConfigurationOption struct {
	DoNotUseCache bool
}

func WithUseCache(useCache bool) GetProviderConfigurationOption {
	return func(opt *getProviderConfigurationOption) {
		opt.DoNotUseCache = !useCache
	}
}

func (d *Discovery) GetProviderConfiguration(ctx context.Context, providerConfigurationURL ProviderConfigurationURL, opts ...GetProviderConfigurationOption) (*ProviderConfigurationResponse, error) {
	return d._GetProviderConfiguration(ctx, providerConfigurationURL, time.Now(), opts...)
}

func (d *Discovery) _GetProviderConfiguration(ctx context.Context, providerConfigurationURL ProviderConfigurationURL, now time.Time, opts ...GetProviderConfigurationOption) (*ProviderConfigurationResponse, error) {
	d.cacheSync.Lock()
	defer d.cacheSync.Unlock()

	opt := new(getProviderConfigurationOption)
	for _, f := range opts {
		f(opt)
	}

	if !d.doNotUseCache && !opt.DoNotUseCache && !d.cacheMap[providerConfigurationURL].Expired() {
		return d.cacheMap[providerConfigurationURL].ProviderConfigurationResponse, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, providerConfigurationURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	r := new(ProviderConfigurationResponse)
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("(*json.Decoder).Decode(*OpenIDProviderMetadata): %w", err)
	}

	if !d.doNotUseCache && !opt.DoNotUseCache {
		d.cacheMap[providerConfigurationURL] = ProviderConfigurationCache{
			ProviderConfigurationResponse: r,
			ExpirationTime:                now.Add(d.cacheTTL),
		}
	}

	return r, nil
}

//nolint:gochecknoglobals
var (
	Default = New()
)

func GetProviderConfiguration(ctx context.Context, providerConfigurationURL ProviderConfigurationURL, opts ...GetProviderConfigurationOption) (*ProviderConfigurationResponse, error) {
	return Default.GetProviderConfiguration(ctx, providerConfigurationURL, opts...)
}
