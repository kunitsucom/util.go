package netz_test

import (
	"net"
	"testing"

	netz "github.com/kunitsuinc/util.go/net"
)

func TestIPNetSet_Contains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})
		ip := net.ParseIP("10.10.10.10")

		s := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})

		if !s.Contains(ip) {
			t.Errorf("%s should contain %s", ipNetSet, ip)
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})
		ip := net.ParseIP("192.168.1.1")

		s := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})

		if s.Contains(ip) {
			t.Errorf("%s should contain %s", ipNetSet, ip)
		}
	})
}
