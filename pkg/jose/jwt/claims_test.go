package jwt //nolint:testpackage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	errorz "github.com/kunitsuinc/util.go/pkg/errors"
	testz "github.com/kunitsuinc/util.go/pkg/test"
)

const (
	testPrivateClaim1Key   = "name"
	testPrivateClaim1Value = "value"
)

var (
	testClaimsSet = &ClaimsSet{
		Issuer:         "http://localhost/iss",
		Subject:        "userID",
		Audience:       []string{"http://localhost/test/aud"},
		ExpirationTime: 1671745431,
		NotBefore:      1671745431,
		IssuedAt:       1671745431,
		JWTID:          "jwtID",
		PrivateClaims: map[string]any{
			testPrivateClaim1Key: testPrivateClaim1Value,
		},
	}
	testClaimsString  = fmt.Sprintf(`{"iss":"http://localhost/iss","sub":"userID","aud":["http://localhost/test/aud"],"exp":1671745431,"nbf":1671745431,"iat":1671745431,"jti":"jwtID","%s":"%s"}`, testPrivateClaim1Key, testPrivateClaim1Value)
	testClaimsEncoded = "eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0L2lzcyIsInN1YiI6InVzZXJJRCIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0L3Rlc3QvYXVkIl0sImV4cCI6MTY3MTc0NTQzMSwibmJmIjoxNjcxNzQ1NDMxLCJpYXQiOjE2NzE3NDU0MzEsImp0aSI6Imp3dElEIiwibmFtZSI6InZhbHVlIn0"
	// NOTE: backup for handling "aud" claim as string
	//	testClaimsString  = fmt.Sprintf(`{"iss":"http://localhost/iss","sub":"userID","aud":"http://localhost/test/aud","exp":1671745431,"nbf":1671745431,"iat":1671745431,"jti":"jwtID","%s":"%s"}`, testPrivateClaim1Key, testPrivateClaim1Value)
	//	testClaimsEncoded = "eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0L2lzcyIsInN1YiI6InVzZXJJRCIsImF1ZCI6Imh0dHA6Ly9sb2NhbGhvc3QvdGVzdC9hdWQiLCJleHAiOjE2NzE3NDU0MzEsIm5iZiI6MTY3MTc0NTQzMSwiaWF0IjoxNjcxNzQ1NDMxLCJqdGkiOiJqd3RJRCIsIm5hbWUiOiJ2YWx1ZSJ9"
)

func TestAudience_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("success,string_to_string", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`"http://localhost/test/aud"`)); err != nil {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
		if expect, actual := "http://localhost/test/aud", aud[0]; expect != actual {
			t.Fatalf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})

	t.Run("success,array_to_string", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`["http://localhost/test/aud"]`)); err != nil {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
		if expect, actual := "http://localhost/test/aud", aud[0]; expect != actual {
			t.Fatalf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})

	t.Run("success,array_to_array", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`["http://localhost/test/aud","http://localhost/test/aud/2"]`)); err != nil {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
		if expect, actual := Audience([]string{"http://localhost/test/aud", "http://localhost/test/aud/2"}), aud; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})

	t.Run("failure,jwt.ErrAudienceIsNil", func(t *testing.T) {
		t.Parallel()
		aud := (*Audience)(nil)
		if expect, err := ErrAudienceIsNil, aud.UnmarshalJSON([]byte(`"aud"`)); !errors.Is(err, expect) {
			t.Fatalf("❌: aud.UnmarshalJSON: expect(%v) != actual(%v)", expect, err)
		}
	})

	t.Run("failure,json.Unmarshal", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`]`)); !errorz.Contains(err, "invalid character ']' looking for beginning of value") {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
	})

	t.Run("failure,not_string_array", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`[0]`)); !errors.Is(err, ErrUnsupportedType) {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
	})

	t.Run("failure,number", func(t *testing.T) {
		t.Parallel()
		aud := make(Audience, 0)
		if err := aud.UnmarshalJSON([]byte(`0`)); !errors.Is(err, ErrUnsupportedType) {
			t.Fatalf("❌: aud.UnmarshalJSON: err != nil: %v", err)
		}
	})
}

func TestClaims_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		claims := new(ClaimsSet)
		if err := json.Unmarshal([]byte(testClaimsString), claims); err != nil {
			t.Fatalf("❌: json.Unmarshal: err != nil: %v", err)
		}
		v, ok := claims.PrivateClaims[testPrivateClaim1Key]
		if !ok {
			t.Fatalf("❌: header.PrivateHeaderParameters[testPrivatePrivateHeaderParameter1Key]: want(%T) != got(%T)", v, claims.PrivateClaims[testPrivateClaim1Key])
		}
		if actual, expect := v, testPrivateClaim1Value; actual != expect {
			t.Fatalf("❌: actual != expect: %v != %v", actual, expect)
		}
		t.Logf("✅: claims: %#v", claims)
	})
}

func TestClaims_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		b, err := json.Marshal(testClaimsSet)
		if err != nil {
			t.Fatalf("❌: json.Marshal: %v", err)
		}
		if actual, expect := string(b), testClaimsString; actual != expect {
			t.Fatalf("❌: actual != expect: %v != %v", actual, expect)
		}
		t.Logf("✅: claims: %s", b)
	})

	t.Run("success(len(PrivateClaims)==0)", func(t *testing.T) {
		t.Parallel()
		expect := []byte(`{"exp":1671745371,"iat":1671745311}`)
		exp := time.Date(2022, 12, 23, 6, 42, 51, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
		h := NewClaimsSet(WithIssuedAt(exp.Add(-1*time.Minute)), WithExpirationTime(exp))
		actual, err := json.Marshal(h)
		if err != nil {
			t.Fatalf("❌: err != nil: %v", err)
		}
		if !bytes.Equal(actual, expect) {
			t.Fatalf("❌: expect != actual: %s != %s", actual, expect)
		}
	})
}

func TestClaims_marshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("failure(json_Marshal)", func(t *testing.T) {
		t.Parallel()
		_, err := testClaimsSet.marshalJSON(
			func(v any) ([]byte, error) { return nil, testz.ErrTestError },
			bytes.HasSuffix,
			bytes.HasPrefix,
		)
		if !errors.Is(err, testz.ErrTestError) {
			t.Fatalf("❌: err != testz.ErrTestError: %v", err)
		}
	})

	t.Run("failure(invalid)", func(t *testing.T) {
		t.Parallel()
		h := &ClaimsSet{
			PrivateClaims: map[string]any{
				"invalid": func() {},
			},
		}
		_, err := h.marshalJSON(
			json.Marshal,
			bytes.HasSuffix,
			bytes.HasPrefix,
		)
		if err == nil {
			t.Fatalf("❌: err == nil: %v", err)
		}
		if expect, actual := "invalid private claims", err.Error(); !strings.Contains(actual, expect) {
			t.Fatalf("❌: expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("failure(bytes_HasSuffix)", func(t *testing.T) {
		t.Parallel()
		_, err := testClaimsSet.marshalJSON(
			json.Marshal,
			func(s, suffix []byte) bool { return false },
			bytes.HasPrefix,
		)
		if !errors.Is(err, ErrInvalidJSON) {
			t.Fatalf("❌: err != ErrInvalidJSON: %v", err)
		}
	})

	t.Run("failure(bytes_HasPrefix)", func(t *testing.T) {
		t.Parallel()
		_, err := testClaimsSet.marshalJSON(
			json.Marshal,
			bytes.HasSuffix,
			func(s, suffix []byte) bool { return false },
		)
		if !errors.Is(err, ErrInvalidJSON) {
			t.Fatalf("❌: err != ErrInvalidJSON: %v", err)
		}
	})
}

func TestClaims_Encode(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		h := NewClaimsSet(
			WithIssuer("http://localhost/iss"),
			WithSubject("userID"),
			WithAudience("http://localhost/test/aud"),
			WithExpirationTime(time.Date(2022, 12, 23, 6, 43, 51, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
			WithNotBefore(time.Date(2022, 12, 23, 6, 43, 51, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
			WithIssuedAt(time.Date(2022, 12, 23, 6, 43, 51, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
			WithJWTID("jwtID"),
			WithPrivateClaim("name", "value"),
		)
		actual, err := h.Encode()
		if err != nil {
			t.Fatalf("❌: err != nil: %v", err)
		}
		if expect := testClaimsEncoded; expect != actual {
			t.Fatalf("❌: expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		h := &ClaimsSet{
			PrivateClaims: map[string]any{
				"invalid": func() {},
			},
		}
		_, err := h.Encode()
		if err == nil {
			t.Fatalf("❌: err == nil: %v", err)
		}
		if expect, actual := "invalid private claims", err.Error(); !strings.Contains(actual, expect) {
			t.Fatalf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}

func TestClaims_Decode(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		actual := new(ClaimsSet)
		if err := actual.Decode(testClaimsEncoded); err != nil {
			t.Fatalf("❌: err != nil: %v", err)
		}
	})

	t.Run("failure(base64.RawURLEncoding.DecodeString)", func(t *testing.T) {
		t.Parallel()
		err := new(ClaimsSet).Decode("inv@lid")
		if expect, actual := "illegal base64 data at input byte 3", err.Error(); !strings.Contains(actual, expect) {
			t.Fatalf("❌: expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("failure(json.Unmarshal)", func(t *testing.T) {
		t.Parallel()
		err := new(ClaimsSet).Decode("aW52QGxpZA") // invalid (base64-encoded)
		if expect, actual := "invalid character 'i' looking for beginning of value", err.Error(); !strings.Contains(actual, expect) {
			t.Fatalf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}

func TestClaimsSet_GetPrivateClaim(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		expect := "testValue"
		c := NewClaimsSet(WithPrivateClaim(testKey, expect))
		c.SetPrivateClaim(testKey, expect)
		var actual string
		if err := c.GetPrivateClaim(testKey, &actual); err != nil {
			t.Fatalf("❌: (*Header).GetPrivateClaim: err != nil: %v", err)
		}
		if expect != actual {
			t.Fatalf("❌: (*Header).GetPrivateClaim: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		type Expect struct {
			expect string
			if1    *ClaimsSet
		}
		expect := &Expect{expect: "test", if1: testClaimsSet}
		c := NewClaimsSet(WithPrivateClaim(testKey, expect))
		c.SetPrivateClaim(testKey, expect)
		var actual *Expect
		if err := c.GetPrivateClaim(testKey, &actual); err != nil {
			t.Fatalf("❌: (*Header).GetPrivateClaim: err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Fatalf("❌: (*Header).GetPrivateClaim: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(jose.ErrValueIsNotPointerOrInterface)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		c := NewClaimsSet()
		if err := c.GetPrivateClaim(testKey, nil); err == nil || !errors.Is(err, ErrVIsNotPointerOrInterface) {
			t.Fatalf("❌: (*Header).GetPrivateClaim: err: %v", err)
		}
	})

	t.Run("failure(jose.ErrPrivateHeaderParameterIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		c := NewClaimsSet()
		var v string
		if err := c.GetPrivateClaim(testKey, &v); err == nil || !errors.Is(err, ErrPrivateClaimIsNotFound) {
			t.Fatalf("❌: (*Header).GetPrivateClaim: err: %v", err)
		}
	})

	t.Run("failure(jose.ErrPrivateHeaderParameterIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		type Expect struct {
			expect string
			if1    any
		}
		expect := &Expect{expect: "test", if1: "test"}
		c := NewClaimsSet(WithPrivateClaim(testKey, expect))
		c.SetPrivateClaim(testKey, expect)
		var actual string
		if err := c.GetPrivateClaim(testKey, &actual); err == nil || !errors.Is(err, ErrPrivateClaimTypeIsNotMatch) {
			t.Fatalf("❌: (*Header).GetPrivateClaim: err: %v", err)
		}
	})
}
