// Package private provides data models for Bithumb Private API responses.
package private

import "fmt"

// OrderChance represents order chance information.
type OrderChance struct {
	// BidFee is the bid fee rate.
	BidFee string `json:"bid_fee"`
	// AskFee is the ask fee rate.
	AskFee string `json:"ask_fee"`
	// MakerBidFee is the maker bid fee rate.
	MakerBidFee string `json:"maker_bid_fee"`
	// MakerAskFee is the maker ask fee rate.
	MakerAskFee string `json:"maker_ask_fee"`
	// Market is the market information.
	Market MarketInfo `json:"market"`
	// BidAccount is the bid account information.
	BidAccount AccountInfo `json:"bid_account"`
	// AskAccount is the ask account information.
	AskAccount AccountInfo `json:"ask_account"`
}

// MarketInfo represents market information.
type MarketInfo struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	OrderTypes []string   `json:"order_types"`
	AskTypes   []string   `json:"ask_types"`
	BidTypes   []string   `json:"bid_types"`
	OrderSides []string   `json:"order_sides"`
	Bid        Constraint `json:"bid"`
	Ask        Constraint `json:"ask"`
	MaxTotal   string     `json:"max_total"`
	State      string     `json:"state"`
}

// Constraint represents trading constraints.
type Constraint struct {
	Currency  string `json:"currency"`
	PriceUnit string `json:"price_unit"`
	MinTotal  string `json:"min_total"`
}

// AccountInfo represents account information.
type AccountInfo struct {
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	Locked              string `json:"locked"`
	AvgBuyPrice         string `json:"avg_buy_price"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"`
	UnitCurrency        string `json:"unit_currency"`
}

// GetOrderChanceRequest is a request to get order chance.
type GetOrderChanceRequest struct {
	// Market is the market code (e.g., "KRW-BTC").
	Market string
}

// Validate checks if the request is valid.
func (r *GetOrderChanceRequest) Validate() error {
	if r.Market == "" {
		return fmt.Errorf("market is required")
	}
	return nil
}
