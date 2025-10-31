package main

import (
	"os"
	"fmt"
	"bufio"
	"dep/clinput"
)

type cli_command struct {
	name	    string
	description string
	callback    func() error
}

func get_cmds() map[string]cli_command {
	return map[string]cli_command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    command_exit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    command_help,
		},
	}

}

func command_exit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Failed to exit")
}

func command_help() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")

	for _, v := range get_cmds() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

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

		e, ok := get_cmds()[inp[0]]
		if !ok {
			fmt.Println("Unknown command")

			continue
		}

		if err := e.callback(); err != nil {
			fmt.Println(err)
		}
	}
}
