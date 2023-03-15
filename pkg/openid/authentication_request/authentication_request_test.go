package authentication_request //nolint:revive,stylecheck,testpackage

import (
	"strings"
	"testing"

	"github.com/kunitsuinc/util.go/pkg/openid/pkce"
)

func TestNewAuthenticationRequest(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		u, err := New(
			"https://server.example.com/connect/authorize",
			[]string{"openid", "email"},
			"code",
			"client_id@server.example.com",
			"http://localhost:8022/auth/oidc",
			"state_very_very_secure_random_string",
			WithResponseMode("query"),
			WithNonce("nonce_very_very_secure_random_string"),
			WithDisplay("page"),
			WithPrompt("select_account"),
			WithMaxAge(3600),
			WithUILocales([]string{"ja", "en"}),
			WithIDTokenHint("id_token_hint"),
			WithLoginHint("login_hint"),
			WithACRValues([]string{"acr_values"}),
			WithCodeChallengeForPKCE("code_challenge_very_very_secure_random_string", pkce.CodeChallengeMethodS256),
			WithAccessType("offline"),
		)
		if err != nil {
			t.Errorf("❌: NewAuthenticationRequest: %v", err)
		}
		expect := `https://server.example.com/connect/authorize?access_type=offline&acr_values=acr_values&client_id=client_id%40server.example.com&code_challenge=mYRMWSyOX6HNnn9xOgpHHSCKwK9fO7l1ZrXQT1aGsEk&code_challenge_method=S256&display=page&id_token_hint=id_token_hint&login_hint=login_hint&max_age=3600&nonce=nonce_very_very_secure_random_string&prompt=select_account&redirect_uri=http%3A%2F%2Flocalhost%3A8022%2Fauth%2Foidc&response_mode=query&response_type=code&scope=openid%20email&state=state_very_very_secure_random_string&ui_locales=ja%20en`
		actual := u.String()
		if actual != expect {
			t.Errorf("❌: actual != expect: %v != %v", actual, expect)
		}
	})

	t.Run("failure(url.Parse)", func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"http://%%",
			[]string{"openid", "email"},
			"code",
			"client_id@server.example.com",
			"http://localhost:8022/auth/oidc",
			"state_very_very_secure_random_string",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `invalid URL escape "%%"`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})

	t.Run("failure(len(scope)==0)", func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"https://server.example.com/connect/authorize",
			nil,
			"code",
			"client_id@server.example.com",
			"http://localhost:8022/auth/oidc",
			"state_very_very_secure_random_string",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `openid: Authentication Request: parameter is empty`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})

	t.Run(`failure(responseType=="")`, func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"https://server.example.com/connect/authorize",
			[]string{"openid", "email"},
			"",
			"client_id@server.example.com",
			"http://localhost:8022/auth/oidc",
			"state_very_very_secure_random_string",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `openid: Authentication Request: parameter is empty`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})

	t.Run(`failure(clientID=="")`, func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"https://server.example.com/connect/authorize",
			[]string{"openid", "email"},
			"code",
			"",
			"http://localhost:8022/auth/oidc",
			"state_very_very_secure_random_string",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `openid: Authentication Request: parameter is empty`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})

	t.Run(`failure(redirectURI=="")`, func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"https://server.example.com/connect/authorize",
			[]string{"openid", "email"},
			"code",
			"client_id@server.example.com",
			"",
			"state_very_very_secure_random_string",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `openid: Authentication Request: parameter is empty`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})

	t.Run(`failure(state=="")`, func(t *testing.T) {
		t.Parallel()

		_, err := New(
			"https://server.example.com/connect/authorize",
			[]string{"openid", "email"},
			"code",
			"client_id@server.example.com",
			"http://localhost:8022/auth/oidc",
			"",
		)
		if err == nil {
			t.Errorf("❌: NewAuthenticationRequest: err == nil")
		}
		const expect = `openid: Authentication Request: parameter is empty`
		if err != nil && !strings.Contains(err.Error(), expect) {
			t.Errorf("❌: NewAuthenticationRequest: not contains `%s`: %v", expect, err)
		}
	})
}
