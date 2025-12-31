# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Bithumb Go SDK
- Public API support (Ticker, OrderBook, Trade, Candlestick)
- Private API support (Account, Order management)
- JWT authentication for Private API
- WebSocket support (Ticker, OrderBook, Trade, MyOrder, MyAsset)
- Automatic reconnection with subscription restoration
- Thread-safe client implementation
- Context support for request cancellation and timeout control
- Options pattern for flexible client configuration
- Comprehensive error handling with typed errors
- Test utilities and mock client for testing

## [0.1.0] - 2026-01-01

### Added
- Public API endpoints:
  - GetTicker
  - GetOrderBook
  - GetRecentTrades
  - GetCandlestick
- Private API endpoints:
  - GetAccount
  - PlaceOrder
  - CancelOrder
- WebSocket subscription types:
  - Ticker
  - OrderBook
  - Trade (Transaction)
  - MyOrder
  - MyAsset
- Client configuration options:
  - WithAPIKey
  - WithBaseURL
  - WithTimeout
  - WithHTTPClient
- Error handling utilities:
  - IsAPIError
  - IsRateLimitError
  - IsNetworkError
  - HasErrorCode
