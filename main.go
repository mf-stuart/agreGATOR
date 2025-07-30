package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mf-stuart/agreGATOR/commands"
	"github.com/mf-stuart/agreGATOR/internal/config"
	"github.com/mf-stuart/agreGATOR/internal/database"
	"github.com/mf-stuart/agreGATOR/internal/middle_ware"
	"os"
)

var s = config.State{}

func main() {
	configData, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		os.Exit(1)
	}
	s.Cfg = &configData

	db, err := sql.Open("postgres", configData.DbUrl)
	dbQueries := database.New(db)
	s.Db = dbQueries
	commands.Cmds.Register("login", commands.HandlerLogin)
	commands.Cmds.Register("register", commands.HandlerRegister)
	commands.Cmds.Register("reset", commands.HandlerReset)
	commands.Cmds.Register("users", commands.HandlerUsers)
	commands.Cmds.Register("agg", commands.HandlerAgg)
	commands.Cmds.Register("addfeed", middle_ware.LoggedIn(commands.HandlerAddFeed))
	commands.Cmds.Register("feeds", commands.HandlerFeeds)
	commands.Cmds.Register("follow", middle_ware.LoggedIn(commands.HandlerFollow))
	commands.Cmds.Register("unfollow", middle_ware.LoggedIn(commands.HandlerUnfollow))
	commands.Cmds.Register("following", middle_ware.LoggedIn(commands.HandlerFollowing))

	cmdName := os.Args[1]
	args := os.Args[2:]

	cmdstruct := commands.Command{cmdName, args}

	err = commands.Cmds.Run(&s, cmdstruct)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
