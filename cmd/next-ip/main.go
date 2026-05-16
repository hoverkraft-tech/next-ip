package main

import (
	"flag"
	"fmt"
	"os"

	nextip "github.com/hoverkraft-tech/next-ip"
)

func main() {
	var count int
	var step int
	flag.IntVar(&count, "count", 1, "number of next IP addresses to output")
	flag.IntVar(&count, "c", 1, "number of next IP addresses to output (shorthand)")
	flag.IntVar(&step, "step", 1, "step used to increase IP addresses")
	flag.IntVar(&step, "s", 1, "step used to increase IP addresses (shorthand)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: next-ip [--count N|-c N] [--step N|-s N] <cidr>")
		os.Exit(1)
	}

	ips, err := nextip.NextIPsWithStep(flag.Arg(0), count, step)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
