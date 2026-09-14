package handler

// TransactionState reports managed transaction completion, never a planned
// action. Unknown states deliberately do not imply rollback or commit.
type TransactionState string

const (
	TransactionNone            TransactionState = "none"
	TransactionCommitted       TransactionState = "committed"
	TransactionRolledBack      TransactionState = "rolled_back"
	TransactionCallerPending   TransactionState = "caller_pending"
	TransactionCommitUnknown   TransactionState = "commit_unknown"
	TransactionRollbackUnknown TransactionState = "rollback_unknown"
	TransactionUnknown         TransactionState = "unknown"
	// Partial means some units committed while other participating units did not.
	TransactionPartial TransactionState = "partial"
	TransactionMixed   TransactionState = "mixed"
)

type TransactionOutcome struct {
	// Unit is an invocation-local ordinal, not a database or transaction handle.
	Unit  int
	State TransactionState
	// Error is a unit-local preparation/completion error, when available.
	Error error
}

// Outcome is a detached completion snapshot for one root invocation. Error is
// the operation/completion error. A later finalizer/publication error does not
// change a reported database commit into a rollback.
type Outcome struct {
	Error        error
	Transactions []TransactionOutcome
}

func (o Outcome) Clone() Outcome {
	o.Transactions = append([]TransactionOutcome(nil), o.Transactions...)
	return o
}

// CommitConfirmed is safe for post-commit publication only when every actual
// managed transaction committed, at least one exists, and completion succeeded.
// No transaction, caller ownership and unsupported custom Data are not proof.
func (o Outcome) CommitConfirmed() bool {
	if o.Error != nil {
		return false
	}
	committed := false
	for _, unit := range o.Transactions {
		if unit.Error != nil {
			return false
		}
		switch unit.State {
		case TransactionNone:
		case TransactionCommitted:
			committed = true
		default:
			return false
		}
	}
	return committed
}

// State summarizes unit reports without flattening mixed multi-database work
// into a misleading success or rollback. Transactions retain the full detail.
func (o Outcome) State() TransactionState {
	state := TransactionNone
	mixed, committed := false, false
	for _, unit := range o.Transactions {
		actual := unit.State
		if actual == "" {
			actual = TransactionUnknown
		}
		if actual == TransactionNone {
			continue
		}
		if actual == TransactionCommitted {
			committed = true
		}
		if state == TransactionNone {
			state = actual
		} else if state != actual {
			mixed = true
		}
	}
	if mixed {
		if committed {
			return TransactionPartial
		}
		return TransactionMixed
	}
	return state
}
