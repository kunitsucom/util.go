package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	Google DocumentURL = "https://accounts.google.com/.well-known/openid-configuration"

	DocumentURLPath = "/.well-known/openid-configuration"
)

type (
	DocumentURL = string

	// https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
	//
	//nolint:tagliatelle
	Document struct { //nolint:revive
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

	DocumentCache struct {
		*Document
		ExpirationTime time.Time
	}
)

func (c DocumentCache) Expired() bool {
	return c.expired(time.Now())
}

func (c DocumentCache) expired(now time.Time) bool {
	return c.ExpirationTime.Before(now)
}

type Discovery struct {
	doNotUseCache bool
	cacheTTL      time.Duration
	cacheMap      map[DocumentURL]DocumentCache
	cacheSync     *sync.Mutex
}

func New(opts ...NewDiscoveryOption) *Discovery {
	d := &Discovery{
		doNotUseCache: false,
		cacheTTL:      2 * time.Minute,
		cacheMap:      make(map[DocumentURL]DocumentCache),
		cacheSync:     new(sync.Mutex),
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

type GetDocumentOption func(*getDocumentOption)

type getDocumentOption struct {
	DoNotUseCache bool
}

func WithUseCache(useCache bool) GetDocumentOption {
	return func(opt *getDocumentOption) {
		opt.DoNotUseCache = !useCache
	}
}

func (d *Discovery) GetDocument(ctx context.Context, discoveryDocumentURL DocumentURL, opts ...GetDocumentOption) (*Document, error) {
	return d.getDocument(ctx, discoveryDocumentURL, time.Now(), opts...)
}

func (d *Discovery) getDocument(ctx context.Context, discoveryDocumentURL DocumentURL, now time.Time, opts ...GetDocumentOption) (*Document, error) {
	d.cacheSync.Lock()
	defer d.cacheSync.Unlock()

	opt := new(getDocumentOption)
	for _, f := range opts {
		f(opt)
	}

	if !d.doNotUseCache && !opt.DoNotUseCache && !d.cacheMap[discoveryDocumentURL].Expired() {
		return d.cacheMap[discoveryDocumentURL].Document, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryDocumentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	r := new(Document)
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("(*json.Decoder).Decode(*discovery.Document): %w", err)
	}

	if !d.doNotUseCache && !opt.DoNotUseCache {
		d.cacheMap[discoveryDocumentURL] = DocumentCache{
			Document:       r,
			ExpirationTime: now.Add(d.cacheTTL),
		}
	}

	return r, nil
}

//nolint:gochecknoglobals
var (
	Default = New()
)

func GetDocument(ctx context.Context, discoveryDocumentURL DocumentURL, opts ...GetDocumentOption) (*Document, error) {
	return Default.GetDocument(ctx, discoveryDocumentURL, opts...)
}
