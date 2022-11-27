package openidtest_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
	"github.com/kunitsuinc/util.go/pkg/must"
	"github.com/kunitsuinc/util.go/pkg/openid/discovery"
	"github.com/kunitsuinc/util.go/pkg/openid/openidtest"
)

func TestStartOpenIDProvider(t *testing.T) {
	t.Parallel()

	addr, metadata, _ := openidtest.StartOpenIDProvider()

	t.Logf("start open id provider: %s", addr)

	func() {
		resp := must.One(http.Get(metadata.Issuer + discovery.ProviderMetadataURLPath)) //nolint:bodyclose
		defer resp.Body.Close()
		got := new(discovery.ProviderMetadata)
		must.Must(json.NewDecoder(resp.Body).Decode(got))
		t.Logf("%s", must.One(json.Marshal(got)))
	}()

	func() {
		resp := must.One(http.Get(metadata.JwksURI)) //nolint:bodyclose
		defer resp.Body.Close()
		got := new(jwk.JWKSet)
		must.Must(json.NewDecoder(resp.Body).Decode(got))
		t.Logf("%s", must.One(json.Marshal(got)))
	}()
}
