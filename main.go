package main

import (
	"os"
	"fmt"
	"bufio"
	"dep/clinput"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Println("Encountered an error", err)
			continue
		}

		inp := clinput.Clean_Input(scanner.Text())

		if len(inp) == 0 {
			fmt.Println("You must input anything")
			continue
		}

		fmt.Printf("Your command was: %s\n", inp[0])
	}
}
