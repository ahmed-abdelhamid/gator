package commands

import (
	"errors"

	"github.com/lib/pq"
)

// pgUniqueViolation is the Postgres SQLSTATE code for a unique-constraint
// violation. Handlers map this to friendly per-resource messages.
const pgUniqueViolation = "23505"

// isUniqueViolation reports whether err wraps a Postgres unique-constraint
// violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation
}
