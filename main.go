package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/sbrown3212/gator/internal/config"
	"github.com/sbrown3212/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{cfg: &cfg}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("error opening database")
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState.db = dbQueries

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUser)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide at least two arguments")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
