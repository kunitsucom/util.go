//nolint:testpackage
package pkce

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCreateCodeVerifier(t *testing.T) {
	t.Parallel()

	t.Run("success(EXAMPLE)", func(t *testing.T) {
		t.Parallel()
		cv, err := CreateCodeVerifier(128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		t.Logf("\ncv=%s", cv)
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const expect = "wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz0123"
		actual, err := createCodeVerifier(bytes.NewBufferString(
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
		), 128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(ErrCodeVerifierLength)", func(t *testing.T) {
		t.Parallel()
		_, actual := createCodeVerifier(bytes.NewBufferString(
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
		), 1)
		if actual == nil {
			t.Errorf("❌: err == nil")
		}
		expect := ErrCodeVerifierLength
		if !errors.Is(actual, expect) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(io.ErrUnexpectedEOF)", func(t *testing.T) {
		t.Parallel()
		_, actual := createCodeVerifier(bytes.NewBufferString("not enough"), 128)
		if actual == nil {
			t.Errorf("❌: err == nil")
		}
		expect := io.ErrUnexpectedEOF
		if !errors.Is(actual, expect) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestCodeVerifier_Encode(t *testing.T) {
	t.Parallel()

	t.Run("success("+CodeChallengeMethodPlainShouldNotBeUsed.String()+")", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			Expect string
			Reader io.Reader
			Length int
			Method CodeChallengeMethod
		}

		cases := []testCase{
			{
				"NzPXGLSWADU6DzlHLA7GM0LX4CZOEywp50vSsDdxYMM",
				bytes.NewBufferString("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567"),
				128,
				CodeChallengeMethodS256,
			},
			{
				"wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz0123",
				bytes.NewBufferString("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567"),
				128,
				CodeChallengeMethodPlainShouldNotBeUsed,
			},
			{
				"wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz0123",
				bytes.NewBufferString("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567"),
				128,
				"no_such_method",
			},
		}

		for _, tc := range cases {
			cv, err := createCodeVerifier(tc.Reader, tc.Length)
			if err != nil {
				t.Errorf("❌: err != nil: %v", err)
			}
			actual := cv.Encode(tc.Method)
			if tc.Expect != actual {
				t.Errorf("❌: expect != actual: %v != %v", tc.Expect, actual)
			}
		}
	})
}
