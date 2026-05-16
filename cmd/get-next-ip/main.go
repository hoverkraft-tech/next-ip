package main

import (
	"flag"
	"fmt"
	"os"

	nextip "github.com/hoverkraft-tech/next-ip"
)

func main() {
	var count int
	flag.IntVar(&count, "count", 1, "number of next IP addresses to output")
	flag.IntVar(&count, "c", 1, "number of next IP addresses to output (shorthand)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: get-next-ip [--count N] <cidr>")
		os.Exit(1)
	}

	ips, err := nextip.NextIPs(flag.Arg(0), count)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
