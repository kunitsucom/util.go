package ecs

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/oauth2/google/externalaccount"

	httpz "github.com/kunitsucom/util.go/net/http"
	testingz "github.com/kunitsucom/util.go/testing"
	"github.com/kunitsucom/util.go/testing/assert"
)

func TestAwsEcsSecurityCredentialsSupplier_AwsRegion(t *testing.T) {
	t.Parallel()

	t.Run("success,os.Getenv", func(t *testing.T) {
		t.Parallel()

		expected := os.Getenv("AWS_REGION")
		if expected == "" {
			expected = os.Getenv("AWS_DEFAULT_REGION")
		}

		s := &AwsEcsSecurityCredentialsSupplier{
			defaultAwsRegion: expected,
		}
		actual, _ := s.AwsRegion(context.Background(), externalaccount.SupplierOptions{})
		assert.Equal(t, expected, actual)
	})

	t.Run("success,DefaultAwsRegion", func(t *testing.T) {
		t.Parallel()

		const expected = "ap-northeast-1"

		s := &AwsEcsSecurityCredentialsSupplier{
			defaultAwsRegion: expected,
		}
		region, err := s.AwsRegion(context.Background(), externalaccount.SupplierOptions{})
		assert.NoError(t, err)
		assert.Equal(t, expected, region)
	})

	t.Run("success,AWS_REGION", func(t *testing.T) {
		t.Parallel()

		const expected = "ap-northeast-1"

		s := &AwsEcsSecurityCredentialsSupplier{
			_osGetenvFunc: func(key string) string {
				if key == AWS_REGION {
					return expected
				}
				return ""
			},
		}

		region, err := s.AwsRegion(context.Background(), externalaccount.SupplierOptions{})
		assert.NoError(t, err)
		assert.Equal(t, expected, region)
	})

	t.Run("success,AWS_DEFAULT_REGION", func(t *testing.T) {
		t.Parallel()

		const expected = "ap-northeast-1"

		s := &AwsEcsSecurityCredentialsSupplier{
			_osGetenvFunc: func(key string) string {
				if key == AWS_DEFAULT_REGION {
					return expected
				}
				return ""
			},
		}

		region, err := s.AwsRegion(context.Background(), externalaccount.SupplierOptions{})
		assert.NoError(t, err)
		assert.Equal(t, expected, region)
	})

	t.Run("error,NoRegion", func(t *testing.T) {
		t.Parallel()
		s := &AwsEcsSecurityCredentialsSupplier{
			_osGetenvFunc: func(_ string) string { return "" },
		}

		region, err := s.AwsRegion(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorIs(t, err, ErrUnableToDetermineAwsRegion)
		assert.Equal(t, "", region)
	})
}

func TestAwsEcsSecurityCredentialsSupplier_AwsSecurityCredentials(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"AccessKeyId":"TestingAccessKeyId","SecretAccessKey":"TestingSecretAccessKey","Token":"TestingToken"}`))
		}))

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient:                 http.DefaultClient,
			awsEcsMetadataEndpointHost: "http://" + metadataServerMock.Listener.Addr().String(),
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.NoError(t, err)
		assert.Equal(t, "TestingAccessKeyId", creds.AccessKeyID)
		assert.Equal(t, "TestingSecretAccessKey", creds.SecretAccessKey)
		assert.Equal(t, "TestingToken", creds.SessionToken)
	})

	t.Run("failure,http.NewRequestWithContext", func(t *testing.T) {
		t.Parallel()

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient:                 http.DefaultClient,
			awsEcsMetadataEndpointHost: "\t",
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorContains(t, err, "unable to create request")
		assert.Nil(t, creds)
	})

	t.Run("failure,h.httpClient.Do", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"AccessKeyId":"TestingAccessKeyId","SecretAccessKey":"TestingSecretAccessKey","Token":"TestingToken"}`))
		}))

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient: &http.Client{
				Transport: httpz.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					return nil, testingz.ErrTestError
				}),
			},
			awsEcsMetadataEndpointHost: "http://" + metadataServerMock.Listener.Addr().String(),
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorIs(t, err, testingz.ErrTestError)
		assert.Nil(t, creds)
	})

	t.Run("failure,io.ReadAll", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"AccessKeyId":"TestingAccessKeyId","SecretAccessKey":"TestingSecretAccessKey","Token":"TestingToken"}`))
		}))

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient:                 http.DefaultClient,
			awsEcsMetadataEndpointHost: "http://" + metadataServerMock.Listener.Addr().String(),
			_ioReadAllFunc:             func(_ io.Reader) ([]byte, error) { return nil, testingz.ErrTestError },
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorIs(t, err, testingz.ErrTestError)
		assert.Nil(t, creds)
	})

	t.Run("failure,resp.StatusCode", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient:                 http.DefaultClient,
			awsEcsMetadataEndpointHost: "http://" + metadataServerMock.Listener.Addr().String(),
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorIs(t, err, ErrUnableToGetAwsCredentials)
		assert.Nil(t, creds)
	})

	t.Run("failure,json.Unmarshal", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{`))
		}))

		sipper := &AwsEcsSecurityCredentialsSupplier{
			httpClient:                 http.DefaultClient,
			awsEcsMetadataEndpointHost: "http://" + metadataServerMock.Listener.Addr().String(),
		}

		creds, err := sipper.AwsSecurityCredentials(context.Background(), externalaccount.SupplierOptions{})
		assert.ErrorContains(t, err, "unable to decode AWS credentials")
		assert.Nil(t, creds)
	})
}
