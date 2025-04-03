# Pokemon TCG Pocket Boosters

[![Build](https://github.com/danielholmes/pokemon-tcg-pocket-boosters/actions/workflows/test.yml/badge.svg)](https://github.com/danielholmes/pokemon-tcg-pocket-boosters/actions/workflows/test.yml)

A utility to work out which booster pack gives you the highest chance of receiving a card you don't have.


## Development setup

 1. Use asdf to install required versions of system dependencies, or otherwise see `.tool-versions` to install them manually in their required versions.
 2. Setup lefthook `lefthook install`


## Configuring

Add a `/collection.json` file which contains a map of `%expansionId% data.ExpansionId` : `[]data.ExpansionNumber`. See `/collection.json.example` for an example.


## Running

```
go run main.go
```

## Running tests

```
go test ./...
```

## TODO

 - Improve simulation performance. Booster instance creation is a good opportunity.
 - Improve booster selection in sim - open booster giving biggest chance.