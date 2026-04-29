package commands

// pgUniqueViolation is the Postgres SQLSTATE code for a unique-constraint
// violation. Handlers map this to friendly per-resource messages.
const pgUniqueViolation = "23505"
