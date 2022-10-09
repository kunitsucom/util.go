package realip_test

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kunitsuinc/util.go/net/http/realip"
	"github.com/kunitsuinc/util.go/netz"
)

const testXForwardedFor = "127.0.0.1, 33.33.33.33, 10.1.1.1, 10.10.10.10, 10.100.100.100"

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		const header = "X-Test-Real-IP"
		expect := "33.33.33.33"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		h := realip.New(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(header)
			}),
			[]*net.IPNet{netz.PrivateIPAddressClassA},
			realip.HeaderXForwardedFor,
			true,
			realip.WithClientIPAddressHeader(header),
		)

		r := httptest.NewRequest(http.MethodGet, "http://util.go/net/realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realip.HeaderXForwardedFor, testXForwardedFor)

		h.ServeHTTP(actualResponse, r)

		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("success(real_ip_header_is_not_X-Forwarded-For)", func(t *testing.T) {
		t.Parallel()

		const testHeaderKey = "Test-Header-Key"

		expect := "33.33.33.33"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		h := realip.New(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realip.HeaderXRealIP)
			}),
			[]*net.IPNet{netz.PrivateIPAddressClassA},
			testHeaderKey,
			true,
		)

		r := httptest.NewRequest(http.MethodGet, "http://util.go/net/realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(testHeaderKey, testXForwardedFor)

		h.ServeHTTP(actualResponse, r)

		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("success(X-Forwarded-For_is_empty)", func(t *testing.T) {
		t.Parallel()

		expect := "192.0.2.1"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		h := realip.New(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realip.HeaderXRealIP)
			}),
			[]*net.IPNet{netz.PrivateIPAddressClassA},
			realip.HeaderXForwardedFor,
			true,
		)

		r := httptest.NewRequest(http.MethodGet, "http://util.go/net/realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realip.HeaderXForwardedFor, "")

		h.ServeHTTP(actualResponse, r)

		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("success(real_ip_recursive_off)", func(t *testing.T) {
		t.Parallel()

		expect := "10.100.100.100"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		h := realip.New(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realip.HeaderXRealIP)
			}),
			[]*net.IPNet{netz.PrivateIPAddressClassA},
			realip.HeaderXForwardedFor,
			false,
		)

		r := httptest.NewRequest(http.MethodGet, "http://util.go/net/realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realip.HeaderXForwardedFor, testXForwardedFor)

		h.ServeHTTP(actualResponse, r)

		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})
}
