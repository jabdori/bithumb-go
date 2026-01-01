package public

// NetworkFee represents deposit/withdrawal fee information for a network
type NetworkFee struct {
	NetName                string `json:"net_name"`
	DepositFeeQuantity     string `json:"deposit_fee_quantity"`
	DepositMinimumQuantity string `json:"deposit_minimum_quantity"`
	WithdrawFeeQuantity    string `json:"withdraw_fee_quantity"`
	WithdrawMinimumQuantity string `json:"withdraw_minimum_quantity"`
}

// ChainFee represents fee information for a currency
type ChainFee struct {
	Name     string       `json:"name"`
	Currency string       `json:"currency"`
	Networks []NetworkFee `json:"networks"`
}
