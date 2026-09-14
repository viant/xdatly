package handler

import "context"

// TransactionStarterKey selects the invocation-owned transaction startup
// capability. It never grants commit/rollback or raw transaction access.
const TransactionStarterKey ValueKey = "transactionStarter"

// TransactionStarter starts or joins the existing managed Data transaction.
// Repeated calls and downstream components join the same database unit. The
// invocation owner completes it; a caller-supplied transaction remains owned
// by the caller. Distinct databases remain distinct units, not one distributed
// transaction.
type TransactionStarter interface {
	Start(context.Context) error
}
