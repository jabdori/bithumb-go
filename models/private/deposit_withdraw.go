package private

import (
	"fmt"
	"strings"
)

// DepositState represents deposit state
type DepositState string

const (
	DepositStateRequestedPending       DepositState = "REQUESTED_PENDING"
	DepositStateRequestedSystemRejected DepositState = "REQUESTED_SYSTEM_REJECTED"
	DepositStateRequestedProcessing     DepositState = "REQUESTED_PROCESSING"
	DepositStateRequestedAdminRejected  DepositState = "REQUESTED_ADMIN_REJECTED"
	DepositStateProcessing              DepositState = "DEPOSIT_PROCESSING"
	DepositStateAccepted                DepositState = "DEPOSIT_ACCEPTED"
	DepositStateCancelled               DepositState = "DEPOSIT_CANCELLED"
)

// WithdrawalState represents withdrawal state
type WithdrawalState string

const (
	WithdrawalStateProcessing WithdrawalState = "PROCESSING"
	WithdrawalStateDone       WithdrawalState = "DONE"
	WithdrawalStateCanceled   WithdrawalState = "CANCELED"
)

// WalletState represents wallet state
type WalletState string

const (
	WalletStateWorking      WalletState = "working"
	WalletStateWithdrawOnly WalletState = "withdraw_only"
	WalletStateDepositOnly  WalletState = "deposit_only"
	WalletStatePaused       WalletState = "paused"
)

// BlockState represents block state
type BlockState string

const (
	BlockStateNormal   BlockState = "normal"
	BlockStateDelayed  BlockState = "delayed"
	BlockStateInactive BlockState = "inactive"
)

// CoinDeposit represents a coin deposit record
type CoinDeposit struct {
	UUID        string       `json:"uuid"`
	Currency    string       `json:"currency"`
	NetType     string       `json:"net_type"`
	Address     string       `json:"address"`
	Amount      string       `json:"amount"`
	State       DepositState `json:"state"`
	TXID        string       `json:"txid"`
	CreatedAt   string       `json:"created_at"`
	CompletedAt string       `json:"completed_at"`
}

// CoinWithdrawal represents a coin withdrawal record
type CoinWithdrawal struct {
	UUID             string          `json:"uuid"`
	Currency         string          `json:"currency"`
	NetType          string          `json:"net_type"`
	Address          string          `json:"address"`
	SecondaryAddress string          `json:"secondary_address"`
	Amount           string          `json:"amount"`
	Fee              string          `json:"fee"`
	State            WithdrawalState `json:"state"`
	TXID             string          `json:"txid"`
	CreatedAt        string          `json:"created_at"`
	CompletedAt      string          `json:"completed_at"`
}

// AssetStatus represents deposit/withdrawal status for a currency
type AssetStatus struct {
	Currency            string      `json:"currency"`
	WalletState         WalletState `json:"wallet_state"`
	BlockState          BlockState  `json:"block_state"`
	BlockHeight         int         `json:"block_height"`
	BlockUpdatedAt      string      `json:"block_updated_at"`
	BlockElapsedMinutes int         `json:"block_elapsed_minutes"`
	NetType             string      `json:"net_type"`
	NetworkName         string      `json:"network_name"`
}

// DepositAddress represents deposit address information
type DepositAddress struct {
	Currency         string `json:"currency"`
	NetType          string `json:"net_type"`
	DepositAddress   string `json:"deposit_address"`
	SecondaryAddress string `json:"secondary_address"`
}

// KRWDeposit represents a KRW deposit record
type KRWDeposit struct {
	UUID        string `json:"uuid"`
	Amount      string `json:"amount"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}

// KRWWithdrawal represents a KRW withdrawal record
type KRWWithdrawal struct {
	UUID        string `json:"uuid"`
	Amount      string `json:"amount"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}

// GetCoinDepositsRequest represents request for coin deposit list
type GetCoinDepositsRequest struct {
	Currency string
	State    DepositState
	UUIDs    []string
	Page     int
	Limit    int
	OrderBy  string
}

// Validate validates GetCoinDepositsRequest
func (r *GetCoinDepositsRequest) Validate() error {
	if r.Limit < 0 || r.Limit > 100 {
		return fmt.Errorf("limit must be between 0 and 100")
	}
	if r.State != "" {
		validStates := []DepositState{
			DepositStateRequestedPending,
			DepositStateRequestedSystemRejected,
			DepositStateRequestedProcessing,
			DepositStateRequestedAdminRejected,
			DepositStateProcessing,
			DepositStateAccepted,
			DepositStateCancelled,
		}
		valid := false
		for _, s := range validStates {
			if r.State == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid state: %s", r.State)
		}
	}
	if r.OrderBy != "" {
		orderByLower := strings.ToLower(r.OrderBy)
		if orderByLower != "asc" && orderByLower != "desc" {
			return fmt.Errorf("order_by must be 'asc' or 'desc', got: %s", r.OrderBy)
		}
	}
	return nil
}

// GetCoinWithdrawalsRequest represents request for coin withdrawal list
type GetCoinWithdrawalsRequest struct {
	Currency string
	State    WithdrawalState
	UUIDs    []string
	TXIDs    []string
	Page     int
	Limit    int
	OrderBy  string
}

// Validate validates GetCoinWithdrawalsRequest
func (r *GetCoinWithdrawalsRequest) Validate() error {
	if r.Limit < 0 || r.Limit > 100 {
		return fmt.Errorf("limit must be between 0 and 100")
	}
	if r.State != "" {
		validStates := []WithdrawalState{
			WithdrawalStateProcessing,
			WithdrawalStateDone,
			WithdrawalStateCanceled,
		}
		valid := false
		for _, s := range validStates {
			if r.State == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid state: %s", r.State)
		}
	}
	if r.OrderBy != "" {
		orderByLower := strings.ToLower(r.OrderBy)
		if orderByLower != "asc" && orderByLower != "desc" {
			return fmt.Errorf("order_by must be 'asc' or 'desc', got: %s", r.OrderBy)
		}
	}
	return nil
}

// GetAssetStatusRequest represents request for asset status
type GetAssetStatusRequest struct {
	Currency string
}

// Validate validates GetAssetStatusRequest (empty for consistency)
func (r *GetAssetStatusRequest) Validate() error {
	return nil
}

// GetDepositAddressRequest represents request for deposit address
type GetDepositAddressRequest struct {
	Currency string
}

// Validate validates GetDepositAddressRequest
func (r *GetDepositAddressRequest) Validate() error {
	if r.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}
