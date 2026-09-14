package handler

import (
	"errors"
	"testing"
)

func TestOutcomeCommitEvidence(t *testing.T) {
	for _, test := range []struct {
		name      string
		states    []TransactionState
		err       error
		summary   TransactionState
		confirmed bool
	}{
		{"no units", nil, nil, TransactionNone, false},
		{"no transaction", []TransactionState{TransactionNone}, nil, TransactionNone, false},
		{"committed", []TransactionState{TransactionCommitted, TransactionNone}, nil, TransactionCommitted, true},
		{"rolled back", []TransactionState{TransactionRolledBack}, nil, TransactionRolledBack, false},
		{"caller pending", []TransactionState{TransactionCallerPending}, nil, TransactionCallerPending, false},
		{"commit unknown", []TransactionState{TransactionCommitUnknown}, nil, TransactionCommitUnknown, false},
		{"rollback unknown", []TransactionState{TransactionRollbackUnknown}, nil, TransactionRollbackUnknown, false},
		{"unsupported", []TransactionState{TransactionUnknown}, nil, TransactionUnknown, false},
		{"partial", []TransactionState{TransactionCommitted, TransactionCommitUnknown}, nil, TransactionPartial, false},
		{"mixed pending", []TransactionState{TransactionCommitted, TransactionCallerPending}, nil, TransactionPartial, false},
		{"mixed no commit", []TransactionState{TransactionRolledBack, TransactionCallerPending}, nil, TransactionMixed, false},
		{"committed with operation error", []TransactionState{TransactionCommitted}, errors.New("operation"), TransactionCommitted, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			o := Outcome{Error: test.err}
			for i, state := range test.states {
				o.Transactions = append(o.Transactions, TransactionOutcome{Unit: i, State: state})
			}
			if o.State() != test.summary || o.CommitConfirmed() != test.confirmed {
				t.Fatalf("summary=%s confirmed=%v", o.State(), o.CommitConfirmed())
			}
			copy := o.Clone()
			if len(copy.Transactions) > 0 {
				copy.Transactions[0].State = "changed"
				if o.Transactions[0].State == "changed" {
					t.Fatal("snapshot shares mutable state")
				}
			}
		})
	}
}
