package main

import (
	"flag"
	"fmt"
	"os"

	nextip "github.com/hoverkraft-tech/next-ip"
)

func main() {
	var count int
	var countShort int
	var step int
	var stepShort int
	flag.IntVar(&count, "count", 1, "number of next IP addresses to output")
	flag.IntVar(&countShort, "c", 0, "number of next IP addresses to output (shorthand)")
	flag.IntVar(&step, "step", 1, "step used to increase IP addresses")
	flag.IntVar(&stepShort, "s", 0, "step used to increase IP addresses (shorthand)")
	flag.Parse()

	if countShort > 0 {
		if count != 1 && count != countShort {
			fmt.Fprintln(os.Stderr, "count and c must have the same value when both are provided")
			os.Exit(1)
		}
		count = countShort
	}

	if stepShort > 0 {
		if step != 1 && step != stepShort {
			fmt.Fprintln(os.Stderr, "step and s must have the same value when both are provided")
			os.Exit(1)
		}
		step = stepShort
	}

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
