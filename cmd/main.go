package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/suyash0x/jump/pkg/jump"
)

func main() {

	listTargets := flag.Bool("l", false, "List all the jump targets")
	addTarget := flag.String("a", "", "Adds jump target with specified jump path")
	deleteTarget := flag.String("d", "", "deletes specified jump target")

	help := flag.Bool("h", false, "show help")

	flag.Parse()

	target := flag.Args()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *listTargets {
		jump.ListTargets()
		os.Exit(0)
	}

	if len(*addTarget) > 0 {
		jump.AddTarget(*addTarget)
		os.Exit(0)
	}

	if len(*deleteTarget) > 0 {
		jump.DeleteTarget(*deleteTarget)
		os.Exit(0)
	}

	if len(target) == 1 {
		jump.InitiateJump(target[0])
		os.Exit(0)
	}

	if len(target) > 1 {
		fmt.Println("Jump to multiple paths is not possible :)")
		os.Exit(0)
	}

	fmt.Println("No target provided")
	os.Exit(1)

}
