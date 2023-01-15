package jwk_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/cache"
	x509z "github.com/kunitsuinc/util.go/pkg/crypto/x509"
	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
	"github.com/kunitsuinc/util.go/pkg/must"
	testz "github.com/kunitsuinc/util.go/pkg/test"
)

func TestJSONWebKey_DecodeRSAPublicKey(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		k1 := k0.EncodeRSAPublicKey(
			must.One(x509z.ParseRSAPublicKeyPEM([]byte(testz.TestRSAPublicKey2048BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("RSA"),
		)
		k2, err := k1.DecodeRSAPublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: err != nil: %v", err)
		}
		if expect, actual := k1.E, base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(k2.E))); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: E: %v", actual)
		}
		if expect, actual := k1.N, base64.RawURLEncoding.EncodeToString(k2.N.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: N: %v", actual)
		}
	})
	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeRSAPublicKey(
			must.One(x509z.ParseRSAPublicKeyPEM([]byte(testz.TestRSAPublicKey2048BitPEM))),
		)
		k1.E = "inv@lid"
		if _, err := k1.DecodeRSAPublicKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.E: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: err != nil: %v", err)
		}
		k1.E = "invalid"
		if _, err := k1.DecodeRSAPublicKey(); err == nil || !strings.Contains(err.Error(), "invalid syntax") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: err != nil: %v", err)
		}
		k1.N = "inv@lid"
		if _, err := k1.DecodeRSAPublicKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.N: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPublicKey: err != nil: %v", err)
		}
	})
}

func TestJSONWebKey_DecodeRSAPrivateKey(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeRSAPrivateKey(
			must.One(x509z.ParseRSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKey2048BitPEM))),
			jwk.WithKeyID("testKeyID"),
		)
		k2, err := k1.DecodeRSAPrivateKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: err != nil: %v", err)
		}
		if expect, actual := k1.E, base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(k2.E))); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: E: %v", actual)
		}
		if expect, actual := k1.N, base64.RawURLEncoding.EncodeToString(k2.N.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: N: %v", actual)
		}
		if expect, actual := k1.D, base64.RawURLEncoding.EncodeToString(k2.D.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: D: %v", actual)
		}
	})
	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeRSAPrivateKey(
			must.One(x509z.ParseRSAPrivateKeyPEM([]byte(testz.TestRSAPrivateKey2048BitPEM))),
		)
		k1.Q = "inv@lid"
		if _, err := k1.DecodeRSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.Q: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: err != nil: %v", err)
		}
		k1.P = "inv@lid"
		if _, err := k1.DecodeRSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.P: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: err != nil: %v", err)
		}
		k1.D = "inv@lid"
		if _, err := k1.DecodeRSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.D: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: err != nil: %v", err)
		}
		k1.N = "inv@lid"
		if _, err := k1.DecodeRSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.N: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeRSAPrivateKey: err != nil: %v", err)
		}
	})
}

func TestJSONWebKey_DecodeECDSAPublicKey(t *testing.T) {
	t.Parallel()
	t.Run("success(P-256)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		k1 := k0.EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("RSA"),
		)
		k2, err := k1.DecodeECDSAPublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
		if expect, actual := k1.Y, base64.RawURLEncoding.EncodeToString(k2.Y.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: Y: %v", actual)
		}
		if expect, actual := k1.X, base64.RawURLEncoding.EncodeToString(k2.X.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: X: %v", actual)
		}
	})
	t.Run("success(P-384)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		k1 := k0.EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey384BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("RSA"),
		)
		k2, err := k1.DecodeECDSAPublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
		if expect, actual := k1.Y, base64.RawURLEncoding.EncodeToString(k2.Y.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: Y: %v", actual)
		}
		if expect, actual := k1.X, base64.RawURLEncoding.EncodeToString(k2.X.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: X: %v", actual)
		}
	})
	t.Run("success(P-521)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		k1 := k0.EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey521BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("RSA"),
		)
		k2, err := k1.DecodeECDSAPublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
		if expect, actual := k1.Y, base64.RawURLEncoding.EncodeToString(k2.Y.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: Y: %v", actual)
		}
		if expect, actual := k1.X, base64.RawURLEncoding.EncodeToString(k2.X.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: X: %v", actual)
		}
	})
	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))),
		)
		k1.Y = "inv@lid"
		if _, err := k1.DecodeECDSAPublicKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.Y: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
		k1.X = "inv@lid"
		if _, err := k1.DecodeECDSAPublicKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.X: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
		k1.Crv = "invalid"
		if _, err := k1.DecodeECDSAPublicKey(); err == nil || !strings.Contains(err.Error(), "crv=invalid: jwk: specified curve parameter is not supported") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPublicKey: err != nil: %v", err)
		}
	})
}

func TestJSONWebKey_DecodeECDSAPrivateKey(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeECDSAPrivateKey(
			must.One(x509z.ParseECDSAPrivateKeyPEM([]byte(testz.TestECDSAPrivateKey256BitPEM))),
			jwk.WithKeyID("testKeyID"),
		)
		k2, err := k1.DecodeECDSAPrivateKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPrivateKey: err != nil: %v", err)
		}
		if expect, actual := k1.D, base64.RawURLEncoding.EncodeToString(k2.D.Bytes()); expect != actual {
			t.Errorf("❌: (*jwk.JSONWebKey).DecodeECDSAPrivateKey: D: %v", actual)
		}
	})
	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		k1 := new(jwk.JSONWebKey).EncodeECDSAPrivateKey(
			must.One(x509z.ParseECDSAPrivateKeyPEM([]byte(testz.TestECDSAPrivateKey256BitPEM))),
		)
		k1.D = "inv@lid"
		if _, err := k1.DecodeECDSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "base64.RawURLEncoding.DecodeString: JSONWebKey.D: illegal base64 data at input byte 3") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPrivateKey: err != nil: %v", err)
		}
		k1.Crv = "invalid"
		if _, err := k1.DecodeECDSAPrivateKey(); err == nil || !strings.Contains(err.Error(), "rv=invalid: jwk: specified curve parameter is not supported") {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodeECDSAPrivateKey: err != nil: %v", err)
		}
	})
}

func TestJSONWebKey_DecodePublicKey(t *testing.T) {
	t.Parallel()
	t.Run("success(RSA)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		rsa1 := k0.EncodeRSAPublicKey(
			must.One(x509z.ParseRSAPublicKeyPEM([]byte(testz.TestRSAPublicKey2048BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("RSA"),
		)
		key1, err := rsa1.DecodePublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: err != nil: %v", err)
		}
		pub1, ok := (key1).(*rsa.PublicKey)
		if !ok {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: *rsa.PublicKey: !ok")
		}
		if rsa1.N != base64.RawURLEncoding.EncodeToString(pub1.N.Bytes()) {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: N: %s", pub1.N)
		}
	})
	t.Run("success(EC)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		rsa1 := k0.EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("EC"),
		)
		key1, err := rsa1.DecodePublicKey()
		if err != nil {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: err != nil: %v", err)
		}
		pub1, ok := (key1).(*ecdsa.PublicKey)
		if !ok {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: *rsa.PublicKey: !ok")
		}
		if rsa1.X != base64.RawURLEncoding.EncodeToString(pub1.X.Bytes()) {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: N: %s", pub1.X)
		}
	})
	t.Run("failure(oct)", func(t *testing.T) {
		t.Parallel()
		var k0 *jwk.JSONWebKey
		rsa1 := k0.EncodeECDSAPublicKey(
			must.One(x509z.ParseECDSAPublicKeyPEM([]byte(testz.TestECDSAPublicKey256BitPEM))),
			jwk.WithAlgorithm("sig"),
			jwk.WithKeyID("testKeyID"),
			jwk.WithKeyType("oct"),
		)

		const expect = "kty=oct: jwk: key is not for algorithm"
		if _, err := rsa1.DecodePublicKey(); err == nil || !strings.Contains(err.Error(), expect) {
			t.Fatalf("❌: (*jwk.JSONWebKey).DecodePublicKey: err != %s: %v", expect, err)
		}
	})
}

func TestJwksURL(t *testing.T) {
	t.Parallel()
	url := "https://www.googleapis.com/oauth2/v3/certs"

	t.Run("check("+url+")", func(t *testing.T) {
		t.Parallel()
		r, err := http.Get(url)
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
		t.Logf("✅: %s:\n"+buf.String(), url)
	})
}

func TestClient_GetJSONWebKey(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	testKid := "testKid"
	mux.HandleFunc("/success/certs", func(w http.ResponseWriter, _ *http.Request) {
		const format = `{"keys":[{"kty":"RSA","kid":"%s","n":"%s"},{"alg":"RS256","kty":"RSA","n":"z8PS6saDU3h5ZbQb3Lwl_Arwgu65ECMi79KUlzx4tqk8bgxtaaHcqyvWqVdsA9H6Q2ZtQhBZivqV4Jg0HoPHcEwv46SEziFQNR2LH86e-WIDI5pk2NKg_9cFMee9Mz7f_NSQJ3uyD1pu86bdUTYhCw57DbEVDOuubClNMUV456dWx7dx5W4kdcQe63vGg9LXQ-9PPz9AL-0ZKr8eQEHp4KRfRUfngjqjYBMTFuuo38l94KR99B04Z-FboGnqYLgNxctwZ9eXbCerb9bV5-Q9Gb3zoo0x1h90tFdgmC2ZU1xcIIjHmFqJ29mSDZHYAAYtMNAeWreK4gqWJunc9o0vpQ","use":"sig","kid":"713fd68c966e29380981edc0164a2f6c06c5702a","e":"AQAB"},{"kty":"RSA","e":"AQAB","alg":"RS256","use":"sig","kid":"27b86dc6938dc327b204333a250ebb43b32e4b3c","n":"1X7rNtYVglDjBJgsBOSv7C6MYU6Mv-yraGOp_AGs777c2UcVGj88dBe9KihGicQ3LqU8Vf5fVhPixVy0GFBS7mJt3qJryyBpmG7sChnJQBwJmZEffZUl_rLtwGli8chbZj_Fpgjd-7t74VQJmn2SYkFqHNB3vrW_I8zmwn7_Enn4N84d4dP5R9UChUSLhuPNKaKA-a4vtTKy1LNoZpbr6LG1_QaWGDKNhgPWR-6l5fmBdaXtUgDmPFwdQZuiBUDfnPQ7t1lSUD2PJMnG3M9DKG5gqpSk1L1AlWsxntideNsKWIviZ5PhCpmzEComWNtFtFrzfAWQvLkBbgb0pwWp5w"}]}`
		_, _ = io.WriteString(w, fmt.Sprintf(format, testKid, time.Now().Format(time.RFC3339Nano)))
	})
	mux.HandleFunc("/failure/400", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
	})
	mux.HandleFunc("/failure/302", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Location", "http://localhost/")
		w.WriteHeader(http.StatusFound)
		_, _ = io.WriteString(w, http.StatusText(http.StatusFound))
	})
	mux.HandleFunc("/failure/invalid", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "")
	})
	s := httptest.NewServer(mux)
	t.Cleanup(func() {
		_ = s.Config.Shutdown(context.Background())
	})
	jwksURI := fmt.Sprintf("http://%s/success/certs", s.Listener.Addr().String())
	foundURI := fmt.Sprintf("http://%s/failure/302", s.Listener.Addr().String())
	badRequestURI := fmt.Sprintf("http://%s/failure/400", s.Listener.Addr().String())
	invalidURI := fmt.Sprintf("http://%s/failure/invalid", s.Listener.Addr().String())

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background(), jwk.WithCacheStore(cache.NewStore[*jwk.JWKSet](context.Background())), jwk.WithHTTPClient(http.DefaultClient))
		jwks1, err := c.GetJWKSet(context.Background(), jwksURI)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		jwk1, err := jwks1.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		first := jwk1.N

		jwks2, err := c.GetJWKSet(context.Background(), jwksURI)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		jwk2, err := jwks2.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		cached := jwk2.N

		if first != cached {
			t.Errorf("❌: first != cached: %v != %v", first, cached)
		}
	})

	t.Run("failure(req)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background())
		_, err := c.GetJWKSet(context.Background(), "http://%%")
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		expect := `invalid URL escape "%%"`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(400)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background())
		_, err := c.GetJWKSet(context.Background(), badRequestURI)
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		expect := `code=400 body="Bad Request": jwk: response is not cacheable`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(302)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background())
		_, err := c.GetJWKSet(context.Background(), foundURI)
		expect := `code=302 body="Found": jwk: response is not cacheable`
		if err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(*json.Decoder)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background())
		_, err := c.GetJWKSet(context.Background(), invalidURI)
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		expect := io.EOF.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(ctx)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(context.Background())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := c.GetJWKSet(ctx, jwksURI)
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		expect := `context canceled`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})
}

func TestGetJSONWebKey(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	testKid := "testKid"
	mux.HandleFunc("/certs", func(w http.ResponseWriter, _ *http.Request) {
		const format = `{"keys":[{"kty":"RSA","kid":"%s","n":"%s"},{"alg":"RS256","kty":"RSA","n":"z8PS6saDU3h5ZbQb3Lwl_Arwgu65ECMi79KUlzx4tqk8bgxtaaHcqyvWqVdsA9H6Q2ZtQhBZivqV4Jg0HoPHcEwv46SEziFQNR2LH86e-WIDI5pk2NKg_9cFMee9Mz7f_NSQJ3uyD1pu86bdUTYhCw57DbEVDOuubClNMUV456dWx7dx5W4kdcQe63vGg9LXQ-9PPz9AL-0ZKr8eQEHp4KRfRUfngjqjYBMTFuuo38l94KR99B04Z-FboGnqYLgNxctwZ9eXbCerb9bV5-Q9Gb3zoo0x1h90tFdgmC2ZU1xcIIjHmFqJ29mSDZHYAAYtMNAeWreK4gqWJunc9o0vpQ","use":"sig","kid":"713fd68c966e29380981edc0164a2f6c06c5702a","e":"AQAB"},{"kty":"RSA","e":"AQAB","alg":"RS256","use":"sig","kid":"27b86dc6938dc327b204333a250ebb43b32e4b3c","n":"1X7rNtYVglDjBJgsBOSv7C6MYU6Mv-yraGOp_AGs777c2UcVGj88dBe9KihGicQ3LqU8Vf5fVhPixVy0GFBS7mJt3qJryyBpmG7sChnJQBwJmZEffZUl_rLtwGli8chbZj_Fpgjd-7t74VQJmn2SYkFqHNB3vrW_I8zmwn7_Enn4N84d4dP5R9UChUSLhuPNKaKA-a4vtTKy1LNoZpbr6LG1_QaWGDKNhgPWR-6l5fmBdaXtUgDmPFwdQZuiBUDfnPQ7t1lSUD2PJMnG3M9DKG5gqpSk1L1AlWsxntideNsKWIviZ5PhCpmzEComWNtFtFrzfAWQvLkBbgb0pwWp5w"}]}`
		_, _ = w.Write([]byte(fmt.Sprintf(format, testKid, time.Now().Format(time.RFC3339Nano))))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(""))
	})
	s := httptest.NewServer(mux)
	jwksURI := fmt.Sprintf("http://%s/certs", s.Listener.Addr().String())

	jwks, err := jwk.GetJWKSet(context.Background(), jwksURI)
	if err != nil {
		t.Errorf("❌: err != nil: %v", err)
	}

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		_, err := jwks.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		_, err := jwks.GetJSONWebKey("no_such_key_id")
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		expect := jwk.ErrKidNotFound.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: expect != err: (%v) != (%v)", expect, err)
		}
	})
}
