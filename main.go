package main

import (
	"flag"
	"fmt"
	"os"
	"x86runtime/vm"
)

func main() {
	var (
		fileName = flag.String("f", "", "Linux-8086 executable file")
		debug    = flag.Bool("m", false, "Enable debug mode")
		err      error
		vm       vm.VM
	)
	flag.Parse()
	err = vm.Load(*fileName, *debug)
	if err != nil {
		goto exception
	}
	vm.Run()
	return
exception:
	fmt.Fprintf(os.Stderr, "%s\n", err)
}
