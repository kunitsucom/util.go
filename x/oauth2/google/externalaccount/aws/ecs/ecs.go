package ecs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google/externalaccount"
)

// env keys
//
//nolint:gosec,revive,stylecheck
const (
	AWS_REGION                             = "AWS_REGION"
	AWS_DEFAULT_REGION                     = "AWS_DEFAULT_REGION"
	AWS_CONTAINER_CREDENTIALS_RELATIVE_URI = "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"
)

// defaults
const (
	DefaultMetadataEndpointHost = "http://169.254.170.2"
)

// defaults
//
//nolint:gochecknoglobals
var (
	DefaultTokenSourceConfigScopes = []string{"https://www.googleapis.com/auth/cloud-platform"}
)

// errors
var (
	ErrUnableToDetermineAwsRegion                    = errors.New("unable to determine AWS region")
	ErrUnableToGetAwsCredentials                     = errors.New("unable to get AWS credentials")
	ErrEnvAwsContainerCredentialsRelativeURIIsNotSet = errors.New(fmt.Sprintf("env %s is not set", AWS_CONTAINER_CREDENTIALS_RELATIVE_URI)) //nolint:revive,gosimple // because the return types of errors.New and fmt.Errorf are not the same, and I explicitly choose to use errors.New.
)

// googleWorkloadIdentityFederationConfig is a configuration for the Google Workload Identity Federation.
//
//nolint:tagliatelle // because this is a field in the Google Workload Identity Federation Config
type googleWorkloadIdentityFederationConfig struct {
	Audience                       string                             `json:"audience"`
	SubjectTokenType               string                             `json:"subject_token_type"`
	ServiceAccountImpersonationURL string                             `json:"service_account_impersonation_url"`
	TokenURL                       string                             `json:"token_url"`
	Scopes                         []string                           `json:"scopes"`
	AwsSecurityCredentialsSupplier *AwsEcsSecurityCredentialsSupplier `json:"-"`
}

// AwsEcsSecurityCredentialsSupplier is a supplier for AWS security credentials.
type AwsEcsSecurityCredentialsSupplier struct {
	httpClient                         *http.Client
	defaultAwsRegion                   string
	awsEcsMetadataEndpointHost         string
	awsContainerCredentialsRelativeURI string

	_osGetenvFunc  func(key string) string
	_ioReadAllFunc func(r io.Reader) ([]byte, error)
}

var _ externalaccount.AwsSecurityCredentialsSupplier = (*AwsEcsSecurityCredentialsSupplier)(nil)

func (h *AwsEcsSecurityCredentialsSupplier) AwsRegion(_ context.Context, _ externalaccount.SupplierOptions) (string, error) {
	if h.defaultAwsRegion != "" {
		return h.defaultAwsRegion, nil
	}

	if v := h._osGetenv(AWS_REGION); v != "" {
		return v, nil
	}

	if v := h._osGetenv(AWS_DEFAULT_REGION); v != "" {
		return v, nil
	}

	return "", ErrUnableToDetermineAwsRegion
}

func (h *AwsEcsSecurityCredentialsSupplier) AwsSecurityCredentials(ctx context.Context, _ externalaccount.SupplierOptions) (*externalaccount.AwsSecurityCredentials, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", h.awsEcsMetadataEndpointHost, h.awsContainerCredentialsRelativeURI), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get AWS credentials: %w", err)
	}
	defer resp.Body.Close()

	body, err := h._ioReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code = %d, body = %s: %w", resp.StatusCode, body, ErrUnableToGetAwsCredentials)
	}

	var cred externalaccount.AwsSecurityCredentials
	if err := json.Unmarshal(body, &cred); err != nil {
		return nil, fmt.Errorf("unable to decode AWS credentials: %w", err)
	}

	return &cred, nil
}

// because this is a testable function
func (h *AwsEcsSecurityCredentialsSupplier) _osGetenv(key string) string {
	if h._osGetenvFunc != nil {
		return h._osGetenvFunc(key)
	}

	return os.Getenv(key)
}

// because this is a testable function
func (h *AwsEcsSecurityCredentialsSupplier) _ioReadAll(r io.Reader) ([]byte, error) {
	if h._ioReadAllFunc != nil {
		return h._ioReadAllFunc(r)
	}

	return io.ReadAll(r) //nolint:wrapcheck
}
