package main

import ( 
	"os"
	"fmt"
	"time"
	"bufio"
	"dep/repl"
	"dep/cache"
	"dep/clinput"
)

func main() {
	const interval = time.Second * 60

	scanner := bufio.NewScanner(os.Stdin)

	new_cache := cache.New_Cache(interval)
	config := &repl.Cfg{Ch: &new_cache}

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
