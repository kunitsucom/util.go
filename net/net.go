package netz

import (
	"net"
)

// nolint: revive,gochecknoglobals
var (
	_, LoopbackAddress, _        = net.ParseCIDR("127.0.0.0/8")
	_, PrivateIPAddressClassA, _ = net.ParseCIDR("10.0.0.0/8")
	_, PrivateIPAddressClassB, _ = net.ParseCIDR("172.16.0.0/12")
	_, PrivateIPAddressClassC, _ = net.ParseCIDR("192.168.0.0/16")
)

type IPNetSet []*net.IPNet

func (ipNetSet IPNetSet) Contains(ip net.IP) bool {
	var contains bool
	for _, ipNet := range ipNetSet {
		contains = contains || ipNet.Contains(ip)
	}

	return contains
}
