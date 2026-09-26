package main

import (
	"net"
	"testing"
)

func TestFormatIPWithoutMask(t *testing.T) {
	ip := net.ParseIP("192.168.0.3")
	if got := formatIP(ip, nil); got != "192.168.0.3" {
		t.Fatalf("expected 192.168.0.3, got %s", got)
	}
}

func TestFormatIPWithMask(t *testing.T) {
	_, subnet, err := net.ParseCIDR("192.168.0.2/20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ip := net.ParseIP("192.168.0.3")
	if got := formatIP(ip, subnet); got != "192.168.0.3/20" {
		t.Fatalf("expected 192.168.0.3/20, got %s", got)
	}
}

func TestFormatIPWithMaskIPv6(t *testing.T) {
	_, subnet, err := net.ParseCIDR("2001:db8::/48")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ip := net.ParseIP("2001:db8::1")
	if got := formatIP(ip, subnet); got != "2001:db8::1/48" {
		t.Fatalf("expected 2001:db8::1/48, got %s", got)
	}
}
