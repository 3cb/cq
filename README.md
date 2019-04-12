# Crypto Quotes
Provides streaming quotes from Coinbase, Bitfinex, and HitBTC cryptocurrency exchanges.  Main page has overview of prices from all three markets while exchange specific pages give detailed quotes with bid-ask spread, volume, etc.

Installation:

```
go get github.com/3cb/cq
```

The package provides a flag to change the color theme from the default dark setting to light:

```
cq -t light
```

It also allows the user to set the flash of the price as either full-cell or numbers only.  While full-cell is the default, you can use a flag to change it:

```
cq -f=false
```