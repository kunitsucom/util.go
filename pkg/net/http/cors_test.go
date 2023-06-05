package httpz_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
)

const (
	localhostOrigin = "http://localhost"
)

func TestNewCORSHandler(t *testing.T) {
	t.Parallel()
	t.Run("success(NoOriginHeader)", func(t *testing.T) {
		t.Parallel()
		handler := httpz.NewCORSHandler()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
	})

	t.Run("success(NoConfig)", func(t *testing.T) {
		t.Parallel()
		handler := httpz.NewCORSHandler()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
	})

	t.Run("success(OPTIONS)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = http.MethodGet + ", " + http.MethodPost
			expectAllowMethods        = "Content-Type, Content-Length"
			expectCode                = http.StatusNoContent
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = "true"
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:         localhostOrigin,
			AllowMethods:        []string{http.MethodGet, http.MethodPost},
			AllowHeaders:        []string{"Content-Type", "Content-Length"},
			AllowCredentials:    true,
			AllowPrivateNetwork: true,
			MaxAge:              86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodOptions)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})

	t.Run("success(OPTIONS,OptionsSuccessStatus)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = http.MethodGet + ", " + http.MethodPost
			expectAllowMethods        = "Content-Type, Content-Length"
			expectCode                = http.StatusNoContent
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:          localhostOrigin,
			AllowMethods:         []string{http.MethodGet, http.MethodPost},
			AllowHeaders:         []string{"Content-Type", "Content-Length"},
			AllowCredentials:     true,
			OptionsSuccessStatus: expectCode,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})

	t.Run("success(OPTIONS,OptionsPassthrough)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = http.MethodGet + ", " + http.MethodPost
			expectAllowMethods        = "Content-Type, Content-Length"
			expectCode                = http.StatusOK
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:        localhostOrigin,
			AllowMethods:       []string{http.MethodGet, http.MethodPost},
			AllowHeaders:       []string{"Content-Type", "Content-Length"},
			AllowCredentials:   true,
			MaxAge:             86400,
			OptionsPassthrough: true,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})
}

func TestCORS_methodAllowed(t *testing.T) {
	t.Parallel()

	t.Run("success(OPTIONS,!methodAllowed,PUT)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = ""
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusNotFound
			expectAllowCredentials    = ""
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:      localhostOrigin,
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodPut)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.NotFoundHandler().ServeHTTP(w, r)
				return
			}
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})

	t.Run("success(OPTIONS,!methodAllowed,AllowedMethods)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = ""
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusNotFound
			expectAllowCredentials    = ""
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:      localhostOrigin,
			AllowMethods:     nil,
			AllowHeaders:     []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodPut)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.NotFoundHandler().ServeHTTP(w, r)
				return
			}
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})
}

func TestCORS_headersAllowed(t *testing.T) {
	t.Parallel()

	t.Run("success(OPTIONS,!headersAllowed,404)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = ""
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusNotFound
			expectAllowCredentials    = ""
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:      localhostOrigin,
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "X-Test")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.NotFoundHandler().ServeHTTP(w, r)
				return
			}
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})

	t.Run("success(OPTIONS,headersAllowed,204)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = http.MethodGet + ", " + http.MethodPost
			expectAllowMethods        = "Content-Type, Content-Length"
			expectCode                = http.StatusNoContent
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:      localhostOrigin,
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodOptions)
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})

	t.Run("success(OPTIONS,headersAllowed,*)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = http.MethodGet + ", " + http.MethodPost
			expectAllowMethods        = "Content-Type, Content-Length"
			expectCode                = http.StatusNoContent
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:      localhostOrigin,
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
			MaxAge:           86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodOptions)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
	})
}

func TestCORS_handleRequest(t *testing.T) {
	t.Parallel()

	t.Run("success(GET,200,NoConfig)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = ""
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusOK
			expectAllowCredentials    = ""
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello"))
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(w.Result().Body); err != nil { //nolint:bodyclose
			t.Fatalf("❌: err != nil: %v", err)
		}
		if expectBody, actualBody := "hello", buf.String(); !reflect.DeepEqual(expectBody, actualBody) {
			t.Fatalf("❌: expectBody = %q, actualBody = %q", expectBody, actualBody)
		}
	})

	t.Run("success(GET,200)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = localhostOrigin
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusOK
			expectAllowCredentials    = "true"
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:         localhostOrigin,
			AllowMethods:        []string{http.MethodGet, http.MethodPost},
			AllowHeaders:        []string{"Content-Type", "Content-Length"},
			ExposeHeaders:       []string{"Content-Type"},
			AllowCredentials:    true,
			AllowPrivateNetwork: true,
			MaxAge:              86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodGet)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello"))
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(w.Result().Body); err != nil { //nolint:bodyclose
			t.Fatalf("❌: err != nil: %v", err)
		}
		if expectBody, actualBody := "hello", buf.String(); !reflect.DeepEqual(expectBody, actualBody) {
			t.Fatalf("❌: expectBody = %q, actualBody = %q", expectBody, actualBody)
		}
	})

	t.Run("success(GET,200)", func(t *testing.T) {
		t.Parallel()
		const (
			expectAllowOrigin         = ""
			expectMethods             = ""
			expectAllowMethods        = ""
			expectCode                = http.StatusOK
			expectAllowCredentials    = ""
			expectAllowPrivateNetwork = ""
		)
		handler := httpz.NewCORSHandler(&httpz.CORSConfig{
			AllowOrigin:         localhostOrigin,
			AllowMethods:        []string{http.MethodGet, http.MethodPost},
			AllowHeaders:        []string{"Content-Type", "Content-Length"},
			ExposeHeaders:       []string{"Content-Type"},
			AllowCredentials:    true,
			AllowPrivateNetwork: true,
			MaxAge:              86400,
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		r.Header.Add(httpz.HeaderOrigin, localhostOrigin)
		r.Header.Add(httpz.HeaderAccessControlRequestMethod, http.MethodPut)
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Type")
		r.Header.Add(httpz.HeaderAccessControlRequestHeaders, "Content-Length")
		r.Header.Add(httpz.HeaderAccessControlRequestPrivateNetwork, "true")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello"))
		})).ServeHTTP(w, r)
		t.Logf("ℹ️: response: code=%d header=%#v", w.Code, w.Header())
		if actualCode := w.Code; !reflect.DeepEqual(expectCode, actualCode) {
			t.Fatalf("❌: expectCode = %d, actualCode %d", expectCode, actualCode)
		}
		if actualOrigin := w.Header().Get(httpz.HeaderAccessControlAllowOrigin); !reflect.DeepEqual(expectAllowOrigin, actualOrigin) {
			t.Fatalf("❌: expectAllowOrigin = %q, actualOrigin = %q", expectAllowOrigin, actualOrigin)
		}
		if actualMethods := w.Header().Get(httpz.HeaderAccessControlAllowMethods); !reflect.DeepEqual(expectMethods, actualMethods) {
			t.Fatalf("❌: expectMethods = %q, actualMethod = %q", expectMethods, actualMethods)
		}
		if actualHeaders := w.Header().Get(httpz.HeaderAccessControlAllowHeaders); !reflect.DeepEqual(expectAllowMethods, actualHeaders) {
			t.Fatalf("❌: expectAllowMethods = %q, actualHeaders = %q", expectAllowMethods, actualHeaders)
		}
		if actualCredentials := w.Header().Get(httpz.HeaderAccessControlAllowCredentials); !reflect.DeepEqual(expectAllowCredentials, actualCredentials) {
			t.Fatalf("❌: expectAllowCredentials = %q, actualCredentials = %q", expectAllowCredentials, actualCredentials)
		}
		if actualPrivateNetwork := w.Header().Get(httpz.HeaderAccessControlAllowPrivateNetwork); !reflect.DeepEqual(expectAllowPrivateNetwork, actualPrivateNetwork) {
			t.Fatalf("❌: expectAllowPrivateNetwork = %q, actualPrivateNetwork = %q", expectAllowPrivateNetwork, actualPrivateNetwork)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(w.Result().Body); err != nil { //nolint:bodyclose
			t.Fatalf("❌: err != nil: %v", err)
		}
		if expectBody, actualBody := "hello", buf.String(); !reflect.DeepEqual(expectBody, actualBody) {
			t.Fatalf("❌: expectBody = %q, actualBody = %q", expectBody, actualBody)
		}
	})
}
