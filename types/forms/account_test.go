package forms

import (
	"fmt"
	"testing"
)

func TestAccountForm(t *testing.T) {
	form := AccountEnable{}
	want := "forms.AccountEnable{Hash:\"\"}"
	got := fmt.Sprintf("%#v", form)
	if got != want {
		t.Errorf("AccountEnable struct mismatch:\n got  %s\n want %s", got, want)
	}
}

func TestAccountForm_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    AccountEnable
		wantHash string
	}{
		{
			name:     "empty struct",
			input:    AccountEnable{},
			wantHash: "",
		},
		{
			name:     "struct with hash",
			input:    AccountEnable{Hash: "abc123"},
			wantHash: "abc123",
		},
		// Add more test cases here as your struct grows
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.Hash != tt.wantHash {
				t.Errorf("Hash mismatch: got %q, want %q", tt.input.Hash, tt.wantHash)
			}
		})
	}
}
