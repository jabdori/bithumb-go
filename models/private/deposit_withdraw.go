package private

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
