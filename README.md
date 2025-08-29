# Pokemon TCG Pocket Boosters

[![Build](https://github.com/danielholmes/pokemon-tcg-pocket-boosters/actions/workflows/test.yml/badge.svg)](https://github.com/danielholmes/pokemon-tcg-pocket-boosters/actions/workflows/test.yml)

A utility to work out which booster pack gives you the highest chance of receiving a card you don't have. Also runs a monte carlo simulation opening boosters until you have a full collection, reporting how long that will take.


## Development setup

 1. Use asdf to install required versions of system dependencies, or otherwise see `.tool-versions` to install them manually in their required versions.
 2. Setup lefthook `lefthook install`


## Configuring

Add a `/data.json` file which contains a map of `%expansionId% data.ExpansionId` : `[]data.ExpansionCardNumber`. See `/data.json.example` for an example.


## Building


## Running

Execute with default options:
```
./ptcgpocket
```

Execute (simulation of 200 runs):
```
./ptcgpocket -r 200
```

Execute (simulation of 15 runs with a known simulation seed):
```
./ptcgpocket -r 15 -s 123
```

Static analysis:
```
go fmt ./...
go vet ./...
go tool staticcheck ./...
```

Tests:
```
go test ./...
```

## TODO

 - Switch to tcgdex
  - https://api.tcgdex.net/v2/en/sets?id=A*
  - https://api.tcgdex.net/v2/en/sets/A2
  - https://api.tcgdex.net/v2/en/cards?id=A*
  - https://api.tcgdex.net/v2/en/cards/A1-005
 - Implement regular + 1 pack probabilities
 - Put trades in the simulation. Could make ideal strategy kind of complex, e.g. should maybe ignore 4D, 3D and 1* 
   when deciding pack openings.
 - Too much info - better ui. Maybe menu choices.
 - Base cards not de-duped between boosters in same expansion. Careful with Eevee example - need moves to work out.
 - Show fractional open packs value.
 - Handle special case of 283 genetic apex. Not in any boosters.
 - Include wonder picks? 
    - pick a rate such as average 3 boosters opened per day.
    - each 3 boosters = 1.25 wonder stamina (includes some for quests)
    - view wonder picks 2x per day = X random booster instances
    - find probability of one of those instances having a missing card
    - then apply probability of 1/5 of picking missing card, then consume stamina
 - Include trades?