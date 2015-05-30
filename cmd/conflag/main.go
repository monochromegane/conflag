package main

import (
	"fmt"
	"os"

	"github.com/monochromegane/conflag"
)

func main() {
	args, err := conflag.ArgsFrom(os.Args[1], os.Args[2:]...)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", args)
}
