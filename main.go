package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mf-stuart/agreGATOR/internal/database"
)

import (
	"fmt"
	"github.com/mf-stuart/agreGATOR/internal/config"
	"os"
)

var s = state{}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	configData, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		os.Exit(1)
	}
	s.cfg = &configData

	db, err := sql.Open("postgres", configData.DbUrl)
	dbQueries := database.New(db)
	s.db = dbQueries
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	cmdName := os.Args[1]
	args := os.Args[2:]

	cmdStruct := command{cmdName, args}

	err = cmds.run(&s, cmdStruct)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
