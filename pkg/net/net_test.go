package netz_test

import (
	"net"
	"testing"

	netz "github.com/kunitsuinc/util.go/pkg/net"
)

func TestCIDRToIPNet(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidr := "10.0.0.0/8"
		ipNet, err := netz.CIDRToIPNet(cidr)
		if err != nil {
			t.Errorf("❌: CIDRToIPNet(%s) returned an error: %s", cidr, err)
		}
		if actual := ipNet.String(); actual != cidr {
			t.Errorf("❌: ipNet.String(%s) != cidr: %s != %s", cidr, actual, cidr)
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		cidr := "FAILURE"
		if _, err := netz.CIDRToIPNet(cidr); err == nil {
			t.Errorf("❌: CIDRToIPNet(%s) returned an error: %s", cidr, err)
		}
	})
}

func TestCIDRsToIPNets(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{"10.0.0.0/8"}
		ipNets, err := netz.CIDRsToIPNets(cidrs)
		if err != nil {
			t.Errorf("❌: CIDRsToIPNets(%s) returned an error: %s", cidrs, err)
		}
		if actual := ipNets[0].String(); actual != cidrs[0] {
			t.Errorf("❌: ipNets[0].String(%s) != cidrs[0]: %s != %s", cidrs[0], actual, cidrs[0])
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{"FAILURE"}
		if _, err := netz.CIDRsToIPNets(cidrs); err == nil {
			t.Errorf("❌: CIDRsToIPNets(%s) returned an error: %s", cidrs, err)
		}
	})
}

func TestIPNetSet_Contains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})
		ip := net.ParseIP("10.10.10.10")

		s := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})

		if !s.Contains(ip) {
			t.Errorf("❌: %s should contain %s", ipNetSet, ip)
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})
		ip := net.ParseIP("192.168.1.1")

		s := netz.IPNetSet([]*net.IPNet{netz.PrivateIPAddressClassA})

		if s.Contains(ip) {
			t.Errorf("❌: %s should contain %s", ipNetSet, ip)
		}
	})
}
