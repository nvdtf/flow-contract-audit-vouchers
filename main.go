package main

import (
	"fmt"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func main() {
	fmt.Println("hi")

	_ = gwtf.NewGoWithTheFlowInMemoryEmulator()
}
