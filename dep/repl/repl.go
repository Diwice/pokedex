package repl

import (
	"os"
	"fmt"
	"net/http"
	"dep/cache"
	"encoding/json"
)

type Cfg struct {
	next 	 string
	previous string
	Ch       *cache.Cache
}

type cli_command struct {
	name	    string
	description string
	Callback    func(*Cfg) error
}

type map_obj struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type map_struct struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []map_obj `json:"results"`
}

func Get_cmds() map[string]cli_command {
	return map[string]cli_command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    command_exit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    command_help,
		},
		"map": {
			name:        "map",
			description: "Shows next page available locations",
			Callback:    command_map,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows previous page of available locations",
			Callback:    command_mapb,
		},
	}

}

func command_exit(c *Cfg) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Failed to exit")
}

func command_help(c *Cfg) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")

	for _, v := range Get_cmds() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func map_sub(reverse bool, c *Cfg) error {
	var sub_url string
	if reverse {
		if c.previous == "" {
			return fmt.Errorf("You are already on the first page")
		} else {
			sub_url = c.previous
		}
	} else {
		if c.next != "" {
			sub_url = c.next
		} else {
			sub_url = "https://pokeapi.co/api/v2/location-area"
		}
	}

	entry, ok := c.Ch.Get(sub_url)
	if ok {
		if reverse { // Don't blame me for these checks. Couldn't find a better fix not rewriting half of the whole package.
			if len(c.previous) > 0 {
				num_slice := c.previous[len(c.previous)-11:len(c.previous)-9]
				if num_slice == "=0" {
					c.previous = ""
				} else {
					num_re := 0
					for _, v := range num_slice {
						num_re = num_re*10 + int(v - '0')
					}
					num_re -= 20

					c.previous = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", num_re)
				}
			}

			c.next = sub_url
		} else {
			c.previous = sub_url

			num_slice := c.next[len(c.next)-11:len(c.next)-9]
			
			if num_slice != "=0" {
				num_re := 0
				for _, v := range num_slice {
					num_re = num_re*10 + int(v - '0')
				}
				num_re += 20
				
				c.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", num_re)
			} else {
				c.next = "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
			}
		}

		fmt.Print(string(entry))

		return nil
	}

	req, err := http.Get(sub_url)
	if err != nil {
		return err
	}

	if req.StatusCode > 299 {
		return fmt.Errorf("Request failed, status code: %d", req.StatusCode)
	}

	defer req.Body.Close()

	var ms map_struct
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&ms); err != nil {
		return err
	}

	if ms.Next != nil {
		c.next = *ms.Next
	}
	
	if ms.Previous != nil {
		c.previous = *ms.Previous
	} else {
		c.previous = ""
	}

	res := ""
	for _, v := range ms.Results {
		res += v.Name+"\n"
	}

	fmt.Print(res)

	c.Ch.Add(sub_url, []byte(res))

	return nil
}

func command_map(c *Cfg) error {
	if err := map_sub(false, c); err != nil {
		return err
	}

	return nil
}

func command_mapb(c *Cfg) error {
	if err := map_sub(true, c); err != nil {
		return err
	}

	return nil
}
