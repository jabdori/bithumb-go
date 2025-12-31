// Package public provides data models for Bithumb Public API responses.
package public

// Market represents a market code information.
type Market struct {
	// Market is the market identifier (e.g., "KRW-BTC").
	Market string `json:"market"`
	// KoreanName is the Korean name of the asset.
	KoreanName string `json:"korean_name"`
	// EnglishName is the English name of the asset.
	EnglishName string `json:"english_name"`
	// MarketWarning is the warning status ("NONE" or "CAUTION").
	MarketWarning string `json:"market_warning"`
}
