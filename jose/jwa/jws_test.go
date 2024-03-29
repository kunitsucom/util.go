package jwa_test

import (
	"errors"
	"testing"

	"github.com/kunitsucom/util.go/jose/jwa"
)

func TestRegisterJWSAlgorithm(t *testing.T) { //nolint:paralleltest
	t.Run("success()", func(t *testing.T) {
		const none2 = "none2"
		jwa.RegisterJWSAlgorithmFunc(none2, jwa.JWS(none2).Sign, jwa.JWS(none2).Verify)
		jwa.DeregisterJWSAlgorithm(none2)
		jwa.RegisterJWSAlgorithm(none2, jwa.JWS(none2))
		jwa.RegisterJWSAlgorithmFunc("TEST", func(key any, signingInput string) (signatureEncoded string, err error) { return "TEST", nil }, func(key any, signingInput, signatureEncoded string) (err error) { return nil })
		if _, err := jwa.JWS("TEST").Sign(0, "TEST"); err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if err := jwa.JWS("TEST").Verify(0, "TEST", "TEST"); err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		jwa.RegisterJWSAlgorithmFunc("TEST2", nil, nil)
		if _, err := jwa.JWS("TEST2").Sign(0, "TEST2"); !errors.Is(err, jwa.ErrNotImplemented) {
			t.Errorf("❌: err != jwa.ErrNotImplemented: %v", err)
		}
		if err := jwa.JWS("TEST2").Verify(0, "TEST2", "TEST2"); !errors.Is(err, jwa.ErrNotImplemented) {
			t.Errorf("❌: err != jwa.ErrNotImplemented: %v", err)
		}
		if _, err := jwa.JWS("NotImplemented").Sign(0, "NotImplemented"); !errors.Is(err, jwa.ErrNotImplemented) {
			t.Errorf("❌: err != jwa.ErrNotImplemented: %v", err)
		}
		if err := jwa.JWS("NotImplemented").Verify(0, "NotImplemented", "NotImplemented"); !errors.Is(err, jwa.ErrNotImplemented) {
			t.Errorf("❌: err != jwa.ErrNotImplemented: %v", err)
		}
	})
}
