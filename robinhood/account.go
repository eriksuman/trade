package robinhood

import (
	"encoding/json"
)

var (
	userURL     = APIURL + "user/"
	userInfoURL = APIURL + userURL + "info/"
	accountsURL = APIURL + "accounts/"
)

type accountType string

const (
	cashAccountType   = "cash"
	marginAccountType = "margin"
)

type Account struct {
	Number                    string `json:"account_number"`
	Cash                      string
	Type                      accountType
	CashBalances              CashBalances   `json:"cash_balances"`
	MarginBalances            MarginBalances `json:"margin_balances"`
	Deactivated               bool
	WithdrawalHalted          bool   `json:"withdrawl_halted"`
	SweepEnabled              bool   `json:"sweep_enabled"`
	DepositHalted             bool   `json:"deposit_halted"`
	OnlyPositionClosingTrades bool   `json:"only_position_closing_trades"`
	CashWithdrawalable        string `json:"cash_available_for_withdrawal"`
	SMA                       string
	SMAHeldForOrders          string `json:"sma_held_for_orders"`
	BuyingPower               string `json:"buying_power"`
	MaxACHEarlyAccessAmount   string `json:"max_ach_early_access_amount"`
	CashHeldForOrders         string `json:"cash_held_for_orders"`
	UnclearedDeposits         string `json:"uncleared_deposits"`
	UnsettledFunds            string `json:"unsettled_funds"`
	UserURL                   string `json:"user"`
	PositionsURL              string `json:"positions"`
	PortfolioURL              string `json:"portfolio"`
}

type CashBalances struct {
	CashHeldForOrders  string `json:"cash_held_for_orders"`
	Cash               string
	BuyingPower        string `json:"buying_power"`
	CashWithdrawalable string `json:"cash_available_for_withdrawl"`
	UnclearedDeposits  string `json:"uncleared_deposits"`
	UnsettledFunds     string `json:"unsettled_funds"`
}

type MarginBalances struct {
	DayTradeBuyingPower               string  `json:"day_trade_buying_power"`
	OvernightBuyingPowerHeldForOrders string  `json:"overnight_buying_power_held_for_orders"`
	CashHeldForOrders                 string  `json:"cash_held_for_orders"`
	DayTradeBuyingPowerHeldForOrders  string  `json:"day_trade_buying_power_held_for_orders"`
	MarkedPatternDayTraderDate        *string `json:"marked_pattern_day_trader_date"`
	Cash                              string
	UnallocatedMarginCash             string `json:"unallocated_margin_cash"`
	CashWithdrawalable                string `json:"cash_available_for_withdrawl"`
	MarginLimit                       string `json:"margin_limit"`
	OvernightBuyingPower              string `json:"overnight_buying_power"`
	UnclearedDeposits                 string `json:"uncleared_deposits"`
	UnsettledFunds                    string `json:"unsettled_funds"`
	DayTradeRatio                     string `json:"day_trade_ratio"`
	OvernightRatio                    string `json:"overnight_ratio"`
	OutstandingInterest               string `json:"outstanding_interest"`
	UnsettledDebit                    string `json:"unsettled_debit"`
}

type User struct {
	Username             string
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string
	ID                   string
	IDInfoURL            string `json:"id_info"`
	BasicInfoURL         string `json:"basic_info"`
	ProfileURL           string `json:"investment_profile"`
	InternationalInfoURL string `json:"international_info"`
	EmploymentURL        string `json:"employment"`
	AdditionalInfoURL    string `json:"additional_info"`
}

type UserInfo struct {
	PhoneNumber        string `json:"phone_number"`
	City               string
	DependentsCount    int `json:"number_dependents"`
	Citizenship        string
	MaritalStatus      string `json:"marital_status"`
	ZIP                string `json:"zipcode"`
	CountryOfResidence string `json:"country_of_residence"`
	State              string
	DOB                string `json:"date_of_birth"`
	Address            string
	TaxID              string `json:"tax_id_ssn"`
	UserURL            string `json:"user"`
}

type accountsResponse struct {
	Results  []Account
	Previous string
	Next     string
}

func UserWithCredentials(creds *Credentials) (*User, error) {
	b, err := getWithAuthorization(creds, userURL)
	if err != nil {
		return nil, err
	}

	u := new(User)
	if err := json.Unmarshal(b, u); err != nil {
		return nil, err
	}

	return u, nil
}

func AccountWithCredentials(creds *Credentials) (*Account, error) {
	b, err := getWithAuthorization(creds, accountsURL)
	if err != nil {
		return nil, err
	}

	a := new(accountsResponse)
	if err := json.Unmarshal(b, a); err != nil {
		return nil, err
	}

	return &a.Results[0], nil
}

func UserInfoWithCredentials(creds *Credentials) (*UserInfo, error) {
	b, err := getWithAuthorization(creds, userInfoURL)
	if err != nil {
		return nil, err
	}

	u := new(UserInfo)
	if err := json.Unmarshal(b, u); err != nil {
		return nil, err
	}

	return u, nil
}
