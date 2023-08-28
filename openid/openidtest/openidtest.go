package openidtest

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/kunitsucom/util.go/discard"
	"github.com/kunitsucom/util.go/jose/jwk"
	"github.com/kunitsucom/util.go/must"
	"github.com/kunitsucom/util.go/openid/discovery"
	testingz "github.com/kunitsucom/util.go/testing"
)

const (
	IDToken        = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ" //nolint:gosec
	IDTokenHeader  = `{"alg":"RS256","typ":"JWT"}`                                                                                                                                                                                                                                                                                                                                                                                                                                                             //nolint:gosec
	IDTokenPayload = `{"sub":"1234567890","name":"John Doe","admin":true,"iat":1516239022}`
)

var ErrInvalidPublicKey = errors.New("openidtest: invalid public key")

func StartOpenIDProvider() (
	addr net.Addr,
	metadata *discovery.ProviderMetadata,
	jwks *jwk.JWKSet,
) {
	mux := http.NewServeMux()
	s := httptest.NewServer(mux)
	iss := fmt.Sprintf("http://%s", s.Listener.Addr())

	// /.well-known/openid-configuration
	metadata = &discovery.ProviderMetadata{
		Issuer:                           iss,
		AuthorizationEndpoint:            must.One(url.JoinPath(iss, "/auth")),
		JwksURI:                          must.One(url.JoinPath(iss, "/certs")),
		ResponseTypesSupported:           []string{"id_token"},
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}
	mux.HandleFunc(discovery.ProviderMetadataURLPath, func(w http.ResponseWriter, r *http.Request) {
		must.Must(json.NewEncoder(w).Encode(metadata))
	})

	// /certs
	pub := must.One(x509.ParsePKIXPublicKey(discard.One(pem.Decode([]byte(testingz.TestRSAPublicKey2048BitPEM))).Bytes)).(*rsa.PublicKey) //nolint:forcetypeassert
	jwks = &jwk.JWKSet{
		Keys: []*jwk.JSONWebKey{
			{
				KeyType:      "JWT",
				PublicKeyUse: "sig",
				KeyID:        "TestKeyID",
				Algorithm:    "RS256",
				N:            base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
				E:            base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatInt(int64(pub.E), 10))),
			},
		},
	}
	mux.HandleFunc("/certs", func(w http.ResponseWriter, r *http.Request) {
		must.Must(json.NewEncoder(w).Encode(jwks))
	})

	return s.Listener.Addr(), metadata, jwks
}
