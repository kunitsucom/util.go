package netz_test

import (
	"fmt"
	"net"
	"testing"

	errorz "github.com/kunitsuinc/util.go/pkg/errors"
	netz "github.com/kunitsuinc/util.go/pkg/net"
)

func TestCIDRToIPNet(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidr := "10.0.0.0/8"
		ipNet, err := netz.ParseCIDR(cidr)
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
		if _, err := netz.ParseCIDR(cidr); err == nil {
			t.Errorf("❌: CIDRToIPNet(%s) returned an error: %s", cidr, err)
		}
	})
}

func TestCIDRsToIPNets(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{"10.0.0.0/8"}
		ipNets, err := netz.ParseCIDRs(cidrs)
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
		if _, err := netz.ParseCIDRs(cidrs); err == nil {
			t.Errorf("❌: CIDRsToIPNets(%s) returned an error: %s", cidrs, err)
		}
	})
}

func TestMustParseCIDRs(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{
			netz.LoopbackAddress.String(),
			netz.LinkLocalAddress.String(),
			netz.PrivateIPAddressClassA.String(),
			netz.PrivateIPAddressClassB.String(),
			netz.PrivateIPAddressClassC.String(),
		}
		ipNets := netz.MustParseCIDRs(cidrs)
		for i := range ipNets {
			if expect, actual := cidrs[i], ipNets[i].String(); expect != actual {
				t.Fatalf("❌: MustParseCIDRs: expect(%s) != actual(%s)", cidrs, actual)
			}
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		const cidr = "FAILURE"
		defer func() {
			err, ok := recover().(error)
			if !ok {
				t.Fatalf("❌: MustParseCIDRs should panic with an error")
			}
			if expect := fmt.Sprintf("ParseCIDRs: ParseCIDR: net.ParseCIDR: cidr=%s: invalid CIDR address: %s", cidr, cidr); !errorz.Contains(err, expect) {
				t.Fatalf("❌: recover: expect(%s) != actual(%v)", expect, err)
			}
		}()
		netz.MustParseCIDRs([]string{cidr})
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
