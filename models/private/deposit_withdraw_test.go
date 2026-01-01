package private

import (
	"testing"
)

func TestWalletState_String(t *testing.T) {
	tests := []struct {
		value   WalletState
		wantStr string
	}{
		{WalletStateWorking, "working"},
		{WalletStateWithdrawOnly, "withdraw_only"},
		{WalletStateDepositOnly, "deposit_only"},
		{WalletStatePaused, "paused"},
	}

	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			if string(tt.value) != tt.wantStr {
				t.Errorf("WalletState = %s, want %s", tt.value, tt.wantStr)
			}
		})
	}
}

func TestDepositState_Values(t *testing.T) {
	states := []DepositState{
		DepositStateRequestedPending,
		DepositStateAccepted,
		DepositStateCancelled,
	}

	for _, state := range states {
		if state == "" {
			t.Error("DepositState should not be empty")
		}
	}
}
