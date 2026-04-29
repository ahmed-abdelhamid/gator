package cli

import (
	"github.com/ahmed-abdelhamid/gator/internal/config"
	"github.com/ahmed-abdelhamid/gator/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *config.Config
}
