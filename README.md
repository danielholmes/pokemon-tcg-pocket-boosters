# Pokemon TCG Pocket Boosters

A utility to work out which booster pack gives you the highest chance of receiving a card you don't already own.


## Dependencies

 1. Use asdf to install required versions of system dependencies, or otherwise see `.tool-versions` to install them manually in their required versions.
 2. Setup lefthook `lefthook install`


## Configuring

Add a `/collection.json` file which contains a map of `%cardSetId% data.CardSetId` : `[]data.CardSetNumber`. See `/collection.json.example` for an example.


## Running

```
go run main.go
```