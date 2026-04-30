package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/database"
)

// requireCurrentUser resolves the user named in s.Cfg.CurrentUserName.
// It rejects an empty config and translates a missing user row into a
// clear "logged-in user does not exist" message.
func requireCurrentUser(ctx context.Context, s *cli.State) (database.User, error) {
	if s.Cfg.CurrentUserName == "" {
		return database.User{}, fmt.Errorf("no user logged in; run `gator login` first")
	}
	user, err := s.DB.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.User{}, fmt.Errorf("logged-in user %q does not exist", s.Cfg.CurrentUserName)
		}
		return database.User{}, fmt.Errorf("get user %q: %w", s.Cfg.CurrentUserName, err)
	}
	return user, nil
}
