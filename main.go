package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ahmed-abdelhamid/gator/internal/cli"
	"github.com/ahmed-abdelhamid/gator/internal/config"
	"github.com/ahmed-abdelhamid/gator/internal/database"
	"github.com/ahmed-abdelhamid/gator/internal/handlers"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("cannot open DB: %v", err)
	}
	defer db.Close()

	s := &cli.State{Cfg: &cfg, DB: database.New(db)}

	cmds := cli.NewCommands()
	cmds.Register("login", handlers.Login)
	cmds.Register("register", handlers.Register)
	cmds.Register("reset", handlers.Reset)

	if len(os.Args) < 2 {
		log.Fatal("usage: gator <command> [args...]")
	}

	cmd := cli.Command{Name: os.Args[1], Args: os.Args[2:]}
	if err := cmds.Run(s, cmd); err != nil {
		log.Fatal(err)
	}
}
