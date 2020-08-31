# go-coinmarketcap-bot

A simple bot to fetch cryptocurrency info from CoinMarketCap and send output to Discord webhook.

To install, run:
```bash
go get -u github.com/hyperreal64/go-coinmarketcap-bot/...
```

The above should install the `go-cmc-bot` binary into `$GOPATH/bin`.
The binary requires the following environment variables be defined:
```bash
export DISCORD_WEBHOOK_URL
export CMC_PRO_API_KEY
```

## TODO
* Support potentially all cryptocurrencies
    + Abstract the cryptocurrency icon urls
    + Abstract the command line arg handling
* Support fetching other types of information about user-specified cryptocurrencies
* Possibly separate/abstract Discord-specific aspects of the JSON payload to make it more portable to other webhooks (e.g. Telegram)
