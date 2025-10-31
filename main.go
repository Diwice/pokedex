package main

import (
	"fmt"
	"dep/clinput"
)

func main() {
	em_string := clinput.Clean_Input("lol")
	fmt.Print(em_string)
}
