# Pokemon TCG Pocket Boosters

A utility to work out which booster pack gives you the highest chance of receiving a card you don't already own.


## Dependencies

 1. Use asdf to install required versions of system dependencies, or otherwise see `.tool-versions` to install them manually in their required versions.
 2. Setup lefthook `lefthook install`


## Running

```
go run main.go
```

## TODO

 - Make missing card #s configurable - not hardcoded
 - Provide separate %s - one including secrets and one not (since secrets are dupes of a non-secret)