package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

// NewTokenSource creates a new token source from Google Workload Identity Federation JSON configuration.
//
// The documentation here mentions the method of Workload Identity Federation using EC2 Instance Metadata,
// but it does not mention the method of Federation using ECS Metadata.
// Additionally, golang.org/x/oauth2/google/externalaccount does not support ECS Metadata by default.
// Therefore, it is possible to enable Federation using ECS Metadata by implementing the
// golang.org/x/oauth2/google/externalaccount.AwsSecurityCredentialsSupplier interface and replacing it in the Config.
//
// example:
//
//	data, _ := os.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
//	ts, _ := NewTokenSource(ctx, data)
//	client, _ := storage.NewClient(ctx, option.WithCredentials(&google.Credentials{TokenSource: ts}))
func NewTokenSource(ctx context.Context, jsonData []byte, opts ...TokenSourceOption) (oauth2.TokenSource, error) { //nolint:ireturn
	cfg, err := TokenSourceConfigFromJSON(jsonData, opts...)
	if err != nil {
		return nil, fmt.Errorf("TokenSourceConfigFromJSON: %w", err)
	}

	ts, err := externalaccount.NewTokenSource(ctx, *cfg)
	if err != nil {
		return nil, fmt.Errorf("externalaccount.NewTokenSource: %w", err)
	}

	return ts, nil
}

// TokenSourceConfigFromJSON creates a new token source config from Google Workload Identity Federation JSON configuration.
//
// The documentation here mentions the method of Workload Identity Federation using EC2 Instance Metadata,
// but it does not mention the method of Federation using ECS Metadata.
// Additionally, golang.org/x/oauth2/google/externalaccount does not support ECS Metadata by default.
// Therefore, it is possible to enable Federation using ECS Metadata by implementing the
// golang.org/x/oauth2/google/externalaccount.AwsSecurityCredentialsSupplier interface and replacing it in the Config.
//
// example:
//
//	data, _ := os.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
//	cfg, _ := TokenSourceConfigFromJSON(data)
//	ts, _ := externalaccount.NewTokenSource(ctx, cfg)
//	client, _ := storage.NewClient(ctx, option.WithCredentials(&google.Credentials{TokenSource: ts}))
func TokenSourceConfigFromJSON(jsonData []byte, opts ...TokenSourceOption) (*externalaccount.Config, error) {
	wicfg := new(googleWorkloadIdentityFederationConfig)
	if err := json.Unmarshal(jsonData, wicfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal google workload identity federation config: json.Unmarshal: %w", err)
	}

	// set defaults **after json.Unmarshal**
	wicfg.Scopes = DefaultTokenSourceConfigScopes
	wicfg.AwsSecurityCredentialsSupplier = &AwsEcsSecurityCredentialsSupplier{
		httpClient:                         http.DefaultClient,
		defaultAwsRegion:                   os.Getenv(AWS_REGION),
		awsEcsMetadataEndpointHost:         DefaultMetadataEndpointHost,
		awsContainerCredentialsRelativeURI: os.Getenv(AWS_CONTAINER_CREDENTIALS_RELATIVE_URI),
	}

	// apply options
	for _, opt := range opts {
		opt.apply(wicfg)
	}

	// check if AWS_CONTAINER_CREDENTIALS_RELATIVE_URI is set
	if wicfg.AwsSecurityCredentialsSupplier.awsContainerCredentialsRelativeURI == "" {
		return nil, ErrEnvAwsContainerCredentialsRelativeURIIsNotSet
	}

	// create externalaccount.Config from googleWorkloadIdentityFederationConfig
	cfg := &externalaccount.Config{
		Audience:                       wicfg.Audience,
		SubjectTokenType:               wicfg.SubjectTokenType,
		ServiceAccountImpersonationURL: wicfg.ServiceAccountImpersonationURL,
		TokenURL:                       wicfg.TokenURL,
		Scopes:                         wicfg.Scopes,
		AwsSecurityCredentialsSupplier: wicfg.AwsSecurityCredentialsSupplier,
	}

	return cfg, nil
}

type TokenSourceOption interface {
	apply(cfg *googleWorkloadIdentityFederationConfig)
}

// WithTokenSourceOptionScopes sets the scopes.
func WithTokenSourceOptionScopes(scopes []string) TokenSourceOption { //nolint:ireturn
	return TokenSourceConfigOptionScopes{scopes: scopes}
}

type TokenSourceConfigOptionScopes struct{ scopes []string }

func (f TokenSourceConfigOptionScopes) apply(cfg *googleWorkloadIdentityFederationConfig) {
	cfg.Scopes = f.scopes
}

// WithTokenSourceOptionHTTPClient sets the HTTP client to be used by the AwsEcsSecurityCredentialsSupplier.
// This allows for custom configurations such as timeouts, transport settings, and other HTTP client options.
func WithTokenSourceOptionHTTPClient(httpClient *http.Client) TokenSourceOption { //nolint:ireturn
	return TokenSourceConfigOptionHTTPClient{httpClient: httpClient}
}

type TokenSourceConfigOptionHTTPClient struct{ httpClient *http.Client }

func (f TokenSourceConfigOptionHTTPClient) apply(cfg *googleWorkloadIdentityFederationConfig) {
	cfg.AwsSecurityCredentialsSupplier.httpClient = f.httpClient
}

// WithTokenSourceOptionDefaultAwsRegion sets the default AWS region.
func WithTokenSourceOptionDefaultAwsRegion(region string) TokenSourceOption { //nolint:ireturn
	return TokenSourceConfigOptionDefaultAwsRegion{defaultAwsRegion: region}
}

type TokenSourceConfigOptionDefaultAwsRegion struct{ defaultAwsRegion string }

func (f TokenSourceConfigOptionDefaultAwsRegion) apply(cfg *googleWorkloadIdentityFederationConfig) {
	cfg.AwsSecurityCredentialsSupplier.defaultAwsRegion = f.defaultAwsRegion
}

// WithTokenSourceOptionAwsEcsMetadataEndpointHost sets the AWS ECS Metadata Endpoint host.
func WithTokenSourceOptionAwsEcsMetadataEndpointHost(host string) TokenSourceOption { //nolint:ireturn
	return TokenSourceConfigOptionAwsEcsMetadataEndpointHost{awsEcsMetadataEndpointHost: host}
}

type TokenSourceConfigOptionAwsEcsMetadataEndpointHost struct{ awsEcsMetadataEndpointHost string }

func (f TokenSourceConfigOptionAwsEcsMetadataEndpointHost) apply(cfg *googleWorkloadIdentityFederationConfig) {
	cfg.AwsSecurityCredentialsSupplier.awsEcsMetadataEndpointHost = f.awsEcsMetadataEndpointHost
}

// WithTokenSourceOptionAwsContainerCredentialsRelativeURI sets the AWS container credentials relative URI.
func WithTokenSourceOptionAwsContainerCredentialsRelativeURI(uri string) TokenSourceOption { //nolint:ireturn
	return TokenSourceConfigOptionAwsContainerCredentialsRelativeURI{awsContainerCredentialsRelativeURI: uri}
}

type TokenSourceConfigOptionAwsContainerCredentialsRelativeURI struct{ awsContainerCredentialsRelativeURI string }

func (f TokenSourceConfigOptionAwsContainerCredentialsRelativeURI) apply(cfg *googleWorkloadIdentityFederationConfig) {
	cfg.AwsSecurityCredentialsSupplier.awsContainerCredentialsRelativeURI = f.awsContainerCredentialsRelativeURI
}
