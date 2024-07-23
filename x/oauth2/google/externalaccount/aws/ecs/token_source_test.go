package ecs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/x/oauth2/google/externalaccount/aws/ecs"
)

func TestTokenSourceConfigFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("success,AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", func(t *testing.T) {
		t.Parallel()

		metadataServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"AccessKeyId":"TestingAccessKeyId","SecretAccessKey":"TestingSecretAccessKey","Token":"TestingToken"}`))
		}))

		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecs.NewTokenSource(
			context.Background(),
			jsonData,
			ecs.WithTokenSourceOptionScopes(ecs.DefaultTokenSourceConfigScopes),
			ecs.WithTokenSourceOptionHTTPClient(http.DefaultClient),
			ecs.WithTokenSourceOptionDefaultAwsRegion("ap-northeast-1"),
			ecs.WithTokenSourceOptionAwsEcsMetadataEndpointHost("http://"+metadataServerMock.Listener.Addr().String()),
			ecs.WithTokenSourceOptionAwsContainerCredentialsRelativeURI("/v2/credentials/00000000-0000-0000-0000-000000000000"),
		)
		assert.NotNil(t, ts)
		assert.NoError(t, err)
		tok, err := ts.Token()
		if err != nil {
			t.Logf("📝: ts.Token: %v", err)
		}

		// error
		assert.Error(t, err) // can't get token because credentials and audience is invalid.
		assert.Nil(t, tok)
	})

	t.Run("failure,ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet,FederationID", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecs.NewTokenSource(context.Background(), jsonData)
		assert.ErrorIs(t, err, ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet)
		assert.Nil(t, ts)
	})

	t.Run("failure,ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet,ServiceAccountImpersonation", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecs.NewTokenSource(context.Background(), jsonData)
		assert.ErrorIs(t, err, ecs.ErrEnvAwsContainerCredentialsRelativeURIIsNotSet)
		assert.Nil(t, ts)
	})

	t.Run("failure,json.Unmarshal", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/0000000000000/locations/global/workloadIdentityPools/testing-pool/providers/testing-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
`)
		ts, err := ecs.NewTokenSource(context.Background(), jsonData)
		assert.ErrorContains(t, err, `failed to unmarshal google workload identity federation config: json.Unmarshal: `)
		assert.Nil(t, ts)
	})

	t.Run("failure,audience", func(t *testing.T) {
		t.Parallel()
		jsonData := []byte(`{
  "type": "external_account",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/testing-service-account@testing-google-project.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token"
}`)
		ts, err := ecs.NewTokenSource(context.Background(), jsonData, ecs.WithTokenSourceOptionAwsContainerCredentialsRelativeURI("/v2/credentials/00000000-0000-0000-0000-000000000000"))
		assert.ErrorContains(t, err, `oauth2/google/externalaccount: Audience must be set`)
		assert.Nil(t, ts)
	})
}
