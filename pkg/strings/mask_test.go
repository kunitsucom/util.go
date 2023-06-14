package stringz_test

import (
	"testing"

	stringz "github.com/kunitsuinc/util.go/pkg/strings"
)

func TestMaskPrefix(t *testing.T) {
	t.Parallel()
	type args struct {
		s            string
		mask         string
		unmaskSuffix int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", "*", 0}, ""},
		{"mask 0", args{"abcd", "*", 0}, "abcd"},
		{"mask 1", args{"abcd", "*", 1}, "*bcd"},
		{"mask 2", args{"abcd", "*", 2}, "**cd"},
		{"mask 3", args{"abcd", "*", 3}, "***d"},
		{"mask 4", args{"abcd", "*", 4}, "****"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := stringz.MaskPrefix(tt.args.s, tt.args.mask, tt.args.unmaskSuffix); got != tt.want {
				t.Errorf("%s: MaskPrefix() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestMaskSuffix(t *testing.T) {
	t.Parallel()
	type args struct {
		s            string
		mask         string
		unmaskPrefix int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", "*", 0}, ""},
		{"mask 0", args{"abcd", "*", 0}, "****"},
		{"mask 1", args{"abcd", "*", 1}, "a***"},
		{"mask 2", args{"abcd", "*", 2}, "ab**"},
		{"mask 3", args{"abcd", "*", 3}, "abc*"},
		{"mask 4", args{"abcd", "*", 4}, "abcd"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := stringz.MaskSuffix(tt.args.s, tt.args.mask, tt.args.unmaskPrefix); got != tt.want {
				t.Errorf("MaskSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
