module pokedex

go 1.25.3

replace dep/clinput => ./dep/clinput

replace dep/repl => ./dep/repl

replace dep/cache => ./dep/cache

require dep/clinput v0.0.0-00010101000000-000000000000

require (
	dep/cache v0.0.0-00010101000000-000000000000
	dep/repl v0.0.0-00010101000000-000000000000
)
