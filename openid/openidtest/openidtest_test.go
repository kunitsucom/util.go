package openidtest_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kunitsucom/util.go/jose/jwk"
	"github.com/kunitsucom/util.go/must"
	"github.com/kunitsucom/util.go/openid/discovery"
	"github.com/kunitsucom/util.go/openid/openidtest"
)

func TestStartOpenIDProvider(t *testing.T) {
	t.Parallel()

	addr, metadata, _ := openidtest.StartOpenIDProvider()

	t.Logf("📝: start open id provider: %s", addr)

	func() {
		resp := must.One(http.Get(metadata.Issuer + discovery.ProviderMetadataURLPath)) //nolint:bodyclose
		defer resp.Body.Close()
		got := new(discovery.ProviderMetadata)
		must.Must(json.NewDecoder(resp.Body).Decode(got))
		t.Logf("📝: %s", must.One(json.Marshal(got)))
	}()

	func() {
		resp := must.One(http.Get(metadata.JwksURI)) //nolint:bodyclose
		defer resp.Body.Close()
		got := new(jwk.JWKSet)
		must.Must(json.NewDecoder(resp.Body).Decode(got))
		t.Logf("📝: %s", must.One(json.Marshal(got)))
	}()
}
