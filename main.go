package main

import (
	"fmt"
	"os"
	"x86runtime/vm"
)

const (
	EXECUTABLE = iota
	FLAG
	PASSTHROUGH
	FINISH
)

func main() {
	var (
		debug = false
		err   error
		vm    vm.VM
		args  []string
	)
	args = make([]string, 0, len(os.Args))
	argStatus := EXECUTABLE
	for _, v := range os.Args {
		switch argStatus {
		case EXECUTABLE:
			argStatus = FLAG
		case FLAG:
			switch v {
			case "-m":
				debug = true
			case "--":
				argStatus = PASSTHROUGH
			default:
				args = append(args, v)
			}
		case PASSTHROUGH:
			args = append(args, v)
		}
	}
	err = vm.Load(args[0], debug, args)
	if err != nil {
		goto exception
	}
	vm.Run()
	return
exception:
	fmt.Fprintf(os.Stderr, "%s\n", err)
}
