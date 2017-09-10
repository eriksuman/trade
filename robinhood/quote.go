package robinhood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var quoteURL = APIURL + "/quotes/"

type Quote struct {
	Symbol        string
	InstrumentURL string  `json:"instrument"`
	Ask           string  `json:"ask_price"`
	AskSize       int     `json:"ask_size"`
	Bid           string  `json:"bid_price"`
	BidSize       int     `json:"bid_size"`
	Last          string  `json:"last_trade_price"`
	LastExtHrs    *string `json:"last_extended_hours_trade_price"`
	PrevClose     string  `json:"previous_close"`
	TradingHalted bool    `json:"trading_halted"`
	HasTraded     bool    `json:"has_traded"`
	LastUpdate    string  `json:"updated_at"`
}

type multiQuote struct {
	Results []*Quote
}

func QuoteForSymbol(s string) (*Quote, error) {
	url := quoteURL + strings.ToUpper(s) + "/"

	b, err := getResponseBody(url)
	if err != nil {
		return nil, err
	}

	q := new(Quote)
	if err := json.Unmarshal(b, q); err != nil {
		return nil, err
	}

	return q, nil
}

func QuoteForSymbols(s []string) ([]*Quote, error) {
	for i, str := range s {
		s[i] = strings.ToUpper(str)
	}
	url := quoteURL + "?symbols=" + strings.Join(s, ",")

	b, err := getResponseBody(url)
	if err != nil {
		return nil, err
	}

	m := new(multiQuote)
	if err := json.Unmarshal(b, m); err != nil {
		return nil, err
	}

	return m.Results, nil
}

func getResponseBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("quote.go: non-OK status: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
