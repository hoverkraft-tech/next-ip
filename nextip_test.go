package nextip

import (
	"errors"
	"testing"
)

func TestNextIPsReturnsRequestedCount(t *testing.T) {
	next, err := NextIPs("192.168.100.13/24", 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"192.168.100.14", "192.168.100.15", "192.168.100.16"}
	if len(next) != len(expected) {
		t.Fatalf("expected %d IPs, got %d", len(expected), len(next))
	}

	for i, ip := range next {
		if ip.String() != expected[i] {
			t.Fatalf("unexpected IP at index %d: got %s want %s", i, ip.String(), expected[i])
		}
	}
}

func TestNextIPsWithStepReturnsRequestedCount(t *testing.T) {
	next, err := NextIPsWithStep("192.168.100.102/24", 3, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"192.168.100.105", "192.168.100.108", "192.168.100.111"}
	if len(next) != len(expected) {
		t.Fatalf("expected %d IPs, got %d", len(expected), len(next))
	}

	for i, ip := range next {
		if ip.String() != expected[i] {
			t.Fatalf("unexpected IP at index %d: got %s want %s", i, ip.String(), expected[i])
		}
	}
}

func TestNextIPsReturnsErrorWhenOutOfSubnet(t *testing.T) {
	_, err := NextIPs("192.168.100.255/24", 1)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, ErrOutOfSubnet) {
		t.Fatalf("expected ErrOutOfSubnet, got %v", err)
	}
}

func TestNextIPsReturnsErrorForInvalidCount(t *testing.T) {
	_, err := NextIPs("192.168.100.1/24", 0)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestNextIPsWithStepReturnsErrorForInvalidStep(t *testing.T) {
	_, err := NextIPsWithStep("192.168.100.1/24", 1, 0)
	if err == nil {
		t.Fatal("expected an error")
	}
}
