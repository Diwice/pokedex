package main

import ( 
	"os"
	"fmt"
	"bufio"
	"dep/repl"
	"dep/clinput"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	config := &repl.Cfg{}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Println("Encountered an error during read of Stdin", err)
			continue
		}

		inp := clinput.Clean_Input(scanner.Text())

		if len(inp) == 0 {
			fmt.Println("You must input anything")
			continue
		}

		e, ok := repl.Get_cmds()[inp[0]]
		if !ok {
			fmt.Println("Unknown command")

			continue
		}

		if err := e.Callback(config); err != nil {
			fmt.Println(err)
		}
	}
}
