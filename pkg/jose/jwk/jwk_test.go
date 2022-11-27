package jwk_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/cache"
	"github.com/kunitsuinc/util.go/pkg/jose/jwk"
)

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
	mux.HandleFunc("/success/certs", func(w http.ResponseWriter, r *http.Request) {
		const format = `{"keys":[{"kty":"RSA","kid":"%s","n":"%s"},{"alg":"RS256","kty":"RSA","n":"z8PS6saDU3h5ZbQb3Lwl_Arwgu65ECMi79KUlzx4tqk8bgxtaaHcqyvWqVdsA9H6Q2ZtQhBZivqV4Jg0HoPHcEwv46SEziFQNR2LH86e-WIDI5pk2NKg_9cFMee9Mz7f_NSQJ3uyD1pu86bdUTYhCw57DbEVDOuubClNMUV456dWx7dx5W4kdcQe63vGg9LXQ-9PPz9AL-0ZKr8eQEHp4KRfRUfngjqjYBMTFuuo38l94KR99B04Z-FboGnqYLgNxctwZ9eXbCerb9bV5-Q9Gb3zoo0x1h90tFdgmC2ZU1xcIIjHmFqJ29mSDZHYAAYtMNAeWreK4gqWJunc9o0vpQ","use":"sig","kid":"713fd68c966e29380981edc0164a2f6c06c5702a","e":"AQAB"},{"kty":"RSA","e":"AQAB","alg":"RS256","use":"sig","kid":"27b86dc6938dc327b204333a250ebb43b32e4b3c","n":"1X7rNtYVglDjBJgsBOSv7C6MYU6Mv-yraGOp_AGs777c2UcVGj88dBe9KihGicQ3LqU8Vf5fVhPixVy0GFBS7mJt3qJryyBpmG7sChnJQBwJmZEffZUl_rLtwGli8chbZj_Fpgjd-7t74VQJmn2SYkFqHNB3vrW_I8zmwn7_Enn4N84d4dP5R9UChUSLhuPNKaKA-a4vtTKy1LNoZpbr6LG1_QaWGDKNhgPWR-6l5fmBdaXtUgDmPFwdQZuiBUDfnPQ7t1lSUD2PJMnG3M9DKG5gqpSk1L1AlWsxntideNsKWIviZ5PhCpmzEComWNtFtFrzfAWQvLkBbgb0pwWp5w"}]}`
		_, _ = w.Write([]byte(fmt.Sprintf(format, testKid, time.Now().Format(time.RFC3339Nano))))
	})
	mux.HandleFunc("/failure/400", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	})
	mux.HandleFunc("/failure/invalid", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(""))
	})
	s := httptest.NewServer(mux)
	t.Cleanup(func() {
		_ = s.Config.Shutdown(context.Background())
	})
	jwksURI := fmt.Sprintf("http://%s/success/certs", s.Listener.Addr().String())
	badRequestURI := fmt.Sprintf("http://%s/failure/400", s.Listener.Addr().String())
	invalidURI := fmt.Sprintf("http://%s/failure/invalid", s.Listener.Addr().String())

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient(jwk.WithCacheStore(cache.NewStore[*jwk.JWKSet]()), jwk.WithHTTPClient(http.DefaultClient))
		jwks1, err := c.GetJWKSet(context.Background(), jwksURI)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		jwk1, err := jwks1.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		first := jwk1.N

		jwks2, err := c.GetJWKSet(context.Background(), jwksURI)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		jwk2, err := jwks2.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		cached := jwk2.N

		if first != cached {
			t.Errorf("first != cached: %v != %v", first, cached)
		}
	})

	t.Run("failure(req)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient()
		_, err := c.GetJWKSet(context.Background(), "http://%%")
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := `invalid URL escape "%%"`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(400)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient()
		_, err := c.GetJWKSet(context.Background(), badRequestURI)
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := `code=400 body="Bad Request": jwk: response is not cacheable`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(*json.Decoder)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient()
		_, err := c.GetJWKSet(context.Background(), invalidURI)
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := io.EOF.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != err: (%v) != (%v)", expect, err)
		}
	})

	t.Run("failure(ctx)", func(t *testing.T) {
		t.Parallel()

		c := jwk.NewClient()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := c.GetJWKSet(ctx, jwksURI)
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := `context canceled`
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != err: (%v) != (%v)", expect, err)
		}
	})
}

func TestGetJSONWebKey(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	testKid := "testKid"
	mux.HandleFunc("/certs", func(w http.ResponseWriter, r *http.Request) {
		const format = `{"keys":[{"kty":"RSA","kid":"%s","n":"%s"},{"alg":"RS256","kty":"RSA","n":"z8PS6saDU3h5ZbQb3Lwl_Arwgu65ECMi79KUlzx4tqk8bgxtaaHcqyvWqVdsA9H6Q2ZtQhBZivqV4Jg0HoPHcEwv46SEziFQNR2LH86e-WIDI5pk2NKg_9cFMee9Mz7f_NSQJ3uyD1pu86bdUTYhCw57DbEVDOuubClNMUV456dWx7dx5W4kdcQe63vGg9LXQ-9PPz9AL-0ZKr8eQEHp4KRfRUfngjqjYBMTFuuo38l94KR99B04Z-FboGnqYLgNxctwZ9eXbCerb9bV5-Q9Gb3zoo0x1h90tFdgmC2ZU1xcIIjHmFqJ29mSDZHYAAYtMNAeWreK4gqWJunc9o0vpQ","use":"sig","kid":"713fd68c966e29380981edc0164a2f6c06c5702a","e":"AQAB"},{"kty":"RSA","e":"AQAB","alg":"RS256","use":"sig","kid":"27b86dc6938dc327b204333a250ebb43b32e4b3c","n":"1X7rNtYVglDjBJgsBOSv7C6MYU6Mv-yraGOp_AGs777c2UcVGj88dBe9KihGicQ3LqU8Vf5fVhPixVy0GFBS7mJt3qJryyBpmG7sChnJQBwJmZEffZUl_rLtwGli8chbZj_Fpgjd-7t74VQJmn2SYkFqHNB3vrW_I8zmwn7_Enn4N84d4dP5R9UChUSLhuPNKaKA-a4vtTKy1LNoZpbr6LG1_QaWGDKNhgPWR-6l5fmBdaXtUgDmPFwdQZuiBUDfnPQ7t1lSUD2PJMnG3M9DKG5gqpSk1L1AlWsxntideNsKWIviZ5PhCpmzEComWNtFtFrzfAWQvLkBbgb0pwWp5w"}]}`
		_, _ = w.Write([]byte(fmt.Sprintf(format, testKid, time.Now().Format(time.RFC3339Nano))))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(""))
	})
	s := httptest.NewServer(mux)
	jwksURI := fmt.Sprintf("http://%s/certs", s.Listener.Addr().String())

	jwks, err := jwk.GetJWKSet(context.Background(), jwksURI)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		_, err := jwks.GetJSONWebKey(testKid)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		_, err := jwks.GetJSONWebKey("no_such_key_id")
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := jwk.ErrKidNotFound.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != err: (%v) != (%v)", expect, err)
		}
	})
}
