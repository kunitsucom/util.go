package netz

import (
	"net"
)

// nolint: revive,gochecknoglobals
var (
	_, PrivateIPClassA, _ = net.ParseCIDR("10.0.0.0/8")
	_, PrivateIPClassB, _ = net.ParseCIDR("172.16.0.0/12")
	_, PrivateIPClassC, _ = net.ParseCIDR("192.168.0.0/16")
)

type IPNetSet []*net.IPNet

func (ipNetSet IPNetSet) Contains(ip net.IP) bool {
	var contains bool
	for _, ipNet := range ipNetSet {
		contains = contains || ipNet.Contains(ip)
	}

	return contains
}
