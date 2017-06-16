# Cryptotrend
## Description
Small utility written in `golang` to gather cryptocurrency price trends.

Currently using the cryptonator API : https://api.cryptonator.com/api/ticker/. See https://www.cryptonator.com/api for usage 

## TODOS

- Use go routines to fetch data for each currency pair
- Implement a time.Ticker() to fetch data at defined intervals
