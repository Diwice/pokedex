package repl

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"dep/cache"
	"math/rand"
	"encoding/json"
)

type Cfg struct {
	next 	 string
	previous string
	Ch       *cache.Cache
	pokedex  map[string]pokemon
}

type cli_command struct {
	name	    string
	description string
	Callback    func(*Cfg, string) error
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

type pokemon struct {
	Name     string `json:"name"`
	Base_exp int    `json:"base_experience"`
}

type poke_encounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type enc_res struct {
	Encounters []poke_encounter `json:"pokemon_encounters"`
}

func Get_cmds() map[string]cli_command {
	return map[string]cli_command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex.",
			Callback:    command_exit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message.",
			Callback:    command_help,
		},
		"map": {
			name:        "map",
			description: "Shows next page available locations.",
			Callback:    command_map,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows previous page of available locations.",
			Callback:    command_mapb,
		},
		"explore": {
			name:        "explore",
			description: "Shows the Pokemons inhabiting the named area. explore <area-name>",
			Callback:    command_explore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a named pokemon. catch <pokemon-name>",
			Callback:    command_catch,
		},
	}

}

func command_exit(c *Cfg, params string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Failed to exit")
}

func command_help(c *Cfg, params string) error {
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

func command_map(c *Cfg, params string) error {
	if err := map_sub(false, c); err != nil {
		return err
	}

	return nil
}

func command_mapb(c *Cfg, params string) error {
	if err := map_sub(true, c); err != nil {
		return err
	}

	return nil
}

func command_explore(c *Cfg, params string) error {
	sub_url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", params)

	entry, ok := c.Ch.Get(sub_url)
	if ok {
		fmt.Print(string(entry))

		return nil
	}

	req, err := http.Get(sub_url)
	if err != nil {
		return err
	}

	if req.StatusCode == 404 {
		c.Ch.Add(sub_url, []byte("Location does not exist!\n"))
		return fmt.Errorf("Location does not exist!")
	} else if req.StatusCode > 299 {
		return fmt.Errorf("Request failed, status code: %d", req.StatusCode)
	}

	defer req.Body.Close()

	var poke_encs enc_res
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&poke_encs); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", params)

	res := ""
	for _, enc := range poke_encs.Encounters {
		res += fmt.Sprintf(" - %s\n", enc.Pokemon.Name)
	}

	fmt.Print(res)

	c.Ch.Add(sub_url, []byte(res))

	return nil
}

func command_catch(c *Cfg, params string) error {
	const base_modifier = 7.0

	sub_url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", params)

	req, err := http.Get(sub_url)
	if err != nil {
		return err
	}

	if req.StatusCode == 404 {
		return fmt.Errorf("Pokemon with this name doesn't exist!")
	} else if req.StatusCode > 299 {
		return fmt.Errorf("Request failed, status code: %d", req.StatusCode)
	}

	defer req.Body.Close()

	fmt.Printf("Throwing a Pokeball at %s...\n", params)

	var res pokemon
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&res); err != nil {
		return err
	}

	rand_src := rand.New(rand.NewSource(time.Now().UnixNano()))
	calc_chance := (608.0/float32(res.Base_exp))*base_modifier
	if rand_src.Float32()*100.0 <= calc_chance {
		fmt.Printf("%s was caught!\n", params)

		if c.pokedex == nil {
			c.pokedex = make(map[string]pokemon)
		}

		c.pokedex[params] = res

		return nil
	}

	fmt.Printf("%s escaped!\n", params)

	return nil
}
