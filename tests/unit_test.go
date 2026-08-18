package tests

import (
	"testing"
	appErrors "wallet-api/errors"
	"wallet-api/services"
)

func TestValidateWithdraw(t *testing.T) {
	tests := []struct {
		name    string
		balance int
		amount  int
		wantErr error
	}{
		{"happy path", 100, 50, nil},
		{"exact balance", 100, 100, nil},
		{"insufficient funds", 50, 100, appErrors.ErrInsufficientFunds},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := services.ValidateWithdraw(tt.balance, tt.amount)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
