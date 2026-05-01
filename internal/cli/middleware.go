package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/database"
)

// AuthedHandlerFunc is the contract for handlers that require a logged-in
// user. Wrap one with RequireLoggedIn to adapt it to HandlerFunc for the
// dispatcher.
type AuthedHandlerFunc func(*State, Command, database.User) error

// RequireLoggedIn adapts h to a HandlerFunc by resolving the user named
// in s.Cfg.CurrentUserName before invoking h. It rejects an empty config
// and translates a missing user row into a clear error.
func RequireLoggedIn(h AuthedHandlerFunc) HandlerFunc {
	return func(s *State, cmd Command) error {
		if s.Cfg.CurrentUserName == "" {
			return fmt.Errorf("no user logged in; run `gator login` first")
		}
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("logged-in user %q does not exist", s.Cfg.CurrentUserName)
			}
			return fmt.Errorf("get user %q: %w", s.Cfg.CurrentUserName, err)
		}
		return h(s, cmd, user)
	}
}
