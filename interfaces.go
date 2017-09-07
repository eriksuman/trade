package main

import "github.com/eriksuman/trade/auth"

type Trader interface {
	Buy(symbol string, quantity int, account auth.Credentials) error
	Sell(symbol string, quantity int, account auth.Credentials) error
}
