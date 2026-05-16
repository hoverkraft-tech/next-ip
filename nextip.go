package nextip

import (
	"errors"
	"fmt"
	"net"
)

var ErrOutOfSubnet = errors.New("next IP is out of subnet")

func NextIPs(cidr string, count int) ([]net.IP, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %w", err)
	}

	current := cloneIP(ip)
	result := make([]net.IP, 0, count)

	for range count {
		current = incrementIP(current)
		if !subnet.Contains(current) {
			return nil, fmt.Errorf("%w for subnet %s", ErrOutOfSubnet, subnet.String())
		}
		result = append(result, cloneIP(current))
	}

	return result, nil
}

func incrementIP(ip net.IP) net.IP {
	next := cloneIP(ip)
	for i := len(next) - 1; i >= 0; i-- {
		next[i]++
		if next[i] != 0 {
			break
		}
	}
	return next
}

func cloneIP(ip net.IP) net.IP {
	cloned := make(net.IP, len(ip))
	copy(cloned, ip)
	return cloned
}
