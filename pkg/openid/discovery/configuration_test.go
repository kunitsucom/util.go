package discovery_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/kunitsuinc/util.go/pkg/must"
	"github.com/kunitsuinc/util.go/pkg/net/http/urlcache"
	"github.com/kunitsuinc/util.go/pkg/openid/discovery"
)

func TestDocumentURL(t *testing.T) {
	t.Parallel()

	t.Run("check(discovery.Google)", func(t *testing.T) {
		t.Parallel()
		r, err := http.Get(discovery.Google)
		if err != nil {
			t.Logf("🤔: http.Get: %v", err)
			return
		}
		defer r.Body.Close()
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, r.Body); err != nil {
			t.Logf("🤔: io.Copy: %v", err)
			return
		}
		t.Logf("✅: %s:\n"+buf.String(), discovery.Google)
	})
}

const testMetadata = `{
  "issuer": "https://server.example.com",
  "authorization_endpoint": "https://server.example.com/connect/authorize",
  "token_endpoint": "https://server.example.com/connect/token",
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "private_key_jwt"
  ],
  "token_endpoint_auth_signing_alg_values_supported": [
    "RS256",
    "ES256"
  ],
  "userinfo_endpoint": "https://server.example.com/connect/userinfo",
  "check_session_iframe": "https://server.example.com/connect/check_session",
  "end_session_endpoint": "https://server.example.com/connect/end_session",
  "jwks_uri": "https://server.example.com/jwks.json",
  "registration_endpoint": "https://server.example.com/connect/register",
  "scopes_supported": [
    "openid",
    "profile",
    "email",
    "address",
    "phone",
    "offline_access"
  ],
  "response_types_supported": [
    "code",
    "code id_token",
    "id_token",
    "token id_token"
  ],
  "acr_values_supported": [
    "urn:mace:incommon:iap:silver",
    "urn:mace:incommon:iap:bronze"
  ],
  "subject_types_supported": [
    "public",
    "pairwise"
  ],
  "userinfo_signing_alg_values_supported": [
    "RS256",
    "ES256",
    "HS256"
  ],
  "userinfo_encryption_alg_values_supported": [
    "RSA1_5",
    "A128KW"
  ],
  "userinfo_encryption_enc_values_supported": [
    "A128CBC-HS256",
    "A128GCM"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256",
    "ES256",
    "HS256"
  ],
  "id_token_encryption_alg_values_supported": [
    "RSA1_5",
    "A128KW"
  ],
  "id_token_encryption_enc_values_supported": [
    "A128CBC-HS256",
    "A128GCM"
  ],
  "request_object_signing_alg_values_supported": [
    "none",
    "RS256",
    "ES256"
  ],
  "display_values_supported": [
    "page",
    "popup"
  ],
  "claim_types_supported": [
    "normal",
    "distributed"
  ],
  "claims_supported": [
    "sub",
    "iss",
    "auth_time",
    "acr",
    "name",
    "given_name",
    "family_name",
    "nickname",
    "profile",
    "picture",
    "website",
    "email",
    "email_verified",
    "locale",
    "zoneinfo",
    "http://example.info/claims/groups"
  ],
  "claims_parameter_supported": true,
  "service_documentation": "http://server.example.com/connect/service_documentation.html",
  "ui_locales_supported": [
    "en-US",
    "en-GB",
    "en-CA",
    "fr-FR",
    "fr-CA"
  ]
}
`

func TestDiscovery_GetDocument(t *testing.T) {
	t.Parallel()

	testDiscovery := discovery.New(discovery.WithURLCacheClient(urlcache.NewClient[*discovery.OpenIDProviderMetadata](http.DefaultClient)))

	// prepare
	mux := http.NewServeMux()
	mux.HandleFunc("/success/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, testMetadata); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/failure/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "<!DOCTYPE html>")
	})
	s := httptest.NewServer(mux)
	t.Cleanup(func() {
		_ = s.Config.Shutdown(context.Background())
	})
	urlBase := fmt.Sprintf("http://%s", s.Listener.Addr())

	successURL := must.One(url.JoinPath(urlBase, "/success/.well-known/openid-configuration"))
	failureURL := must.One(url.JoinPath(urlBase, "/failure/.well-known/openid-configuration"))

	// success
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		document, err := discovery.GetOpenIDProviderMetadata(context.Background(), successURL)
		if err != nil {
			t.Errorf("❌: discovery.GetProviderConfiguration: err != nil: %v", err)
		}
		// use cache
		if _, err := discovery.GetOpenIDProviderMetadata(context.Background(), successURL); err != nil {
			t.Errorf("❌: discovery.GetProviderConfiguration: err != nil: %v", err)
		}
		// not use cache
		if _, err := discovery.GetOpenIDProviderMetadata(context.Background(), successURL); err != nil {
			t.Errorf("❌: discovery.GetProviderConfiguration: err != nil: %v", err)
		}
		const AuthorizationEndpoint = `https://server.example.com/connect/authorize`
		if document.AuthorizationEndpoint != AuthorizationEndpoint {
			t.Errorf("❌: document.AuthorizationEndpoint != %s", AuthorizationEndpoint)
		}
		t.Logf("✅: *Document: %v", document)
	})

	// failure
	t.Run("failure(url)", func(t *testing.T) {
		t.Parallel()
		targetURL := "http://%%"
		document, err := testDiscovery.GetOpenIDProviderMetadata(context.Background(), targetURL)
		if err == nil {
			t.Errorf("❌: testDiscovery.GetDocument: err == nil")
		}
		const expect = `invalid URL escape "%%"`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: testDiscovery.GetDocument: %s: not contains `%s`: %v", targetURL, expect, err)
		}
		t.Logf("✅: *Document: %v", document)
	})

	t.Run("failure(ctx)", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		targetURL := successURL
		document, err := testDiscovery.GetOpenIDProviderMetadata(ctx, targetURL)
		if err == nil {
			t.Errorf("❌: testDiscovery.GetDocument: err == nil")
		}
		const expect = `context canceled`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: testDiscovery.GetDocument: %s: not contains `%s`: %v", targetURL, expect, err)
		}
		t.Logf("✅: *Document: %v", document)
	})

	t.Run("failure(Decode)", func(t *testing.T) {
		t.Parallel()
		targetURL := failureURL
		document, err := testDiscovery.GetOpenIDProviderMetadata(context.Background(), failureURL)
		if err == nil {
			t.Errorf("❌: testDiscovery.GetDocument: err == nil")
		}
		const expect = `invalid character '<' looking for beginning of value`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: testDiscovery.GetDocument: %s: not contains `%s`: %v", targetURL, expect, err)
		}
		t.Logf("✅: *Document: %v", document)
	})
}
