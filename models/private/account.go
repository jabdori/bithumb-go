// Package private provides data models for Bithumb Private API responses.
package private

// Account represents account information for a specific currency.
type Account struct {
	// Currency is the currency code (e.g., "BTC", "ETH").
	Currency string `json:"currency"`
	// Balance is the available balance.
	Balance string `json:"balance"`
	// Locked is the locked/hold balance.
	Locked string `json:"locked"`
	// AvgBuyPrice is the average buy price.
	AvgBuyPrice string `json:"avg_buy_price"`
	// AvgBuyPriceModified indicates if avg buy price was modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified"`
	// UnitCurrency is the unit currency.
	UnitCurrency string `json:"unit_currency"`
}

// GetAccountRequest is a request to get account information.
type GetAccountRequest struct {
	// Currency is the currency code. Empty string returns all accounts.
	Currency string
}

// Validate checks if the request is valid.
func (r *GetAccountRequest) Validate() error {
	return nil // All fields are optional
}
