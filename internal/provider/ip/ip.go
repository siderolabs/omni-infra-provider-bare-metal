// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ip provides IP address related functionality.
package ip

import (
	"fmt"
	"net"
)

// RoutableIPs returns a list of routable IP addresses.
func RoutableIPs() ([]string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	routableIPs := make([]string, 0, len(addresses))

	for _, addr := range addresses {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if isRoutableIP(ipNet.IP) {
			routableIPs = append(routableIPs, ipNet.IP.String())
		}
	}

	return routableIPs, nil
}

func isRoutableIP(ip net.IP) bool {
	isReservedIPv4 := func(ip net.IP) bool {
		return ip[0] >= 240
	}

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}

	if ip.To4() != nil {
		return !isReservedIPv4(ip)
	}

	return true
}
