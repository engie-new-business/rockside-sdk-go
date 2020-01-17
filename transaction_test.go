package rockside

import (
	"strings"
	"testing"
)

func TestValidateTransactionFields(t *testing.T) {
	endpoint := new(TransactionEndpoint)
	validAddress := "0x268ba693540A7176ae5d3ba9256A18efbe0A63FF"
	tests := []struct {
		tx          Transaction
		errContains string
	}{
		{tx: Transaction{From: "123", To: ""}, errContains: "'from' address"},
		{tx: Transaction{From: validAddress, To: "34567898"}, errContains: "'to' address"},
		{tx: Transaction{From: "1245", To: "34567898"}, errContains: "'from' address"},
		{tx: Transaction{From: validAddress, Data: "456a789"}, errContains: "'data' bytes"},
		{tx: Transaction{From: validAddress, Value: "456a789"}, errContains: "'value' number"},
	}

	for i, test := range tests {
		_, err := endpoint.Send(test.tx)
		if test.errContains == "" && err != nil {
			t.Fatalf("case %d: unexpected error %s", i+1, err)
		}
		if test.errContains != "" && err == nil {
			t.Fatalf("case %d: expected error, got none", i+1)
		}
		if sub := test.errContains; sub != "" && !strings.Contains(err.Error(), sub) {
			t.Fatalf("case %d: expecting error %q to contains %q", i+1, err, sub)
		}
	}
}
