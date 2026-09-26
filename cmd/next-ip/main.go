package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	nextip "github.com/hoverkraft-tech/next-ip"
)

func main() {
	var count int
	var step int
	var mask bool
	flag.IntVar(&count, "count", 1, "number of next IP addresses to output")
	flag.IntVar(&count, "c", 1, "number of next IP addresses to output (shorthand)")
	flag.IntVar(&step, "step", 1, "step used to increase IP addresses")
	flag.IntVar(&step, "s", 1, "step used to increase IP addresses (shorthand)")
	flag.BoolVar(&mask, "mask", false, "display the netmask (prefix length) of each IP returned")
	flag.BoolVar(&mask, "m", false, "display the netmask (prefix length) of each IP returned (shorthand)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: next-ip [--count N|-c N] [--step N|-s N] [--mask|-m] <cidr>")
		os.Exit(1)
	}

	cidr := flag.Arg(0)

	var subnet *net.IPNet
	if mask {
		var err error
		if _, subnet, err = net.ParseCIDR(cidr); err != nil {
			fmt.Fprintf(os.Stderr, "invalid CIDR: %v\n", err)
			os.Exit(1)
		}
	}

	ips, err := nextip.NextIPsWithStep(cidr, count, step)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, ip := range ips {
		fmt.Println(formatIP(ip, subnet))
	}
}

// formatIP renders ip, appending the subnet prefix length when subnet is set.
func formatIP(ip net.IP, subnet *net.IPNet) string {
	if subnet == nil {
		return ip.String()
	}
	ones, _ := subnet.Mask.Size()
	return fmt.Sprintf("%s/%d", ip.String(), ones)
}
