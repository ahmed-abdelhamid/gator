package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/commands"
	"github.com/ahmed-abdelhamid/gator/internal/config"
	"github.com/ahmed-abdelhamid/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	s := &cli.State{Cfg: &cfg, DB: database.New(db), Conn: db}

	cmds := cli.NewCommands()
	cmds.Register("login", commands.Login)
	cmds.Register("register", commands.Register)
	cmds.Register("reset", commands.Reset)
	cmds.Register("users", commands.Users)
	cmds.Register("agg", commands.Agg)
	cmds.Register("addfeed", cli.RequireLoggedIn(commands.AddFeed))
	cmds.Register("feeds", commands.Feeds)
	cmds.Register("follow", cli.RequireLoggedIn(commands.Follow))
	cmds.Register("following", cli.RequireLoggedIn(commands.Following))

	if len(os.Args) < 2 {
		return fmt.Errorf("usage: gator <command> [args...]")
	}

	cmd := cli.Command{Name: os.Args[1], Args: os.Args[2:]}
	return cmds.Run(s, cmd)
}
