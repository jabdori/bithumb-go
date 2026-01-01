package public

// WarningType represents warning alert types
type WarningType string

const (
	WarningPriceSuddenFluctuation          WarningType = "PRICE_SUDDEN_FLUCTUATION"
	WarningTradingVolumeSuddenFluctuation  WarningType = "TRADING_VOLUME_SUDDEN_FLUCTUATION"
	WarningDepositAmountSuddenFluctuation  WarningType = "DEPOSIT_AMOUNT_SUDDEN_FLUCTUATION"
	WarningPriceDifferenceHigh             WarningType = "PRICE_DIFFERENCE_HIGH"
	WarningSpecificAccountHighTransaction  WarningType = "SPECIFIC_ACCOUNT_HIGH_TRANSACTION"
	WarningExchangeTradingConcentration    WarningType = "EXCHANGE_TRADING_CONCENTRATION"
)

// Warning represents a market warning alert
type Warning struct {
	Market      string     `json:"market"`
	WarningType WarningType `json:"warning_type"`
	EndDate     string     `json:"end_date"`
}
