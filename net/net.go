package netz

import (
	"fmt"
	"net"
)

//nolint:revive,gochecknoglobals
var (
	_, LoopbackAddress, _        = net.ParseCIDR("127.0.0.0/8")
	_, LinkLocalAddress, _       = net.ParseCIDR("169.254.0.0/16")
	_, PrivateIPAddressClassA, _ = net.ParseCIDR("10.0.0.0/8")
	_, PrivateIPAddressClassB, _ = net.ParseCIDR("172.16.0.0/12")
	_, PrivateIPAddressClassC, _ = net.ParseCIDR("192.168.0.0/16")
)

func ParseCIDR(cidr string) (*net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("net.ParseCIDR: cidr=%s: %w", cidr, err)
	}
	ipNet.IP = ip
	return ipNet, nil
}

func MustParseCIDRs(cidrs []string) []*net.IPNet {
	ipNets, err := ParseCIDRs(cidrs)
	if err != nil {
		panic(fmt.Errorf("ParseCIDRs: %w", err))
	}
	return ipNets
}

func ParseCIDRs(cidrs []string) ([]*net.IPNet, error) {
	ipNets := make([]*net.IPNet, len(cidrs))
	for idx, cidr := range cidrs {
		ipNet, err := ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("ParseCIDR: %w", err)
		}
		ipNets[idx] = ipNet
	}
	return ipNets, nil
}

type IPNetSet []*net.IPNet

func (ipNetSet IPNetSet) Contains(ip net.IP) bool {
	var contains bool
	for _, ipNet := range ipNetSet {
		contains = contains || ipNet.Contains(ip)
	}

	return contains
}
