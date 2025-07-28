package main

import (
	"errors"
	"fmt"
	"github.com/mf-stuart/agreGATOR/internal/config"
	"os"
)

var cmds = commands{make(map[string]func(*state, command) error)}
var s = state{}

type state struct {
	cfg *config.Config
}

type commands struct {
	commands map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return cmdFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no arguments specified")
	}

	err := s.cfg.SetUsername(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", cmd.args[0])
	return nil
}

func main() {
	configData, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
	}
	s.cfg = &configData
	cmds.register("login", handlerLogin)

	cmdName := os.Args[1]
	args := os.Args[2:]
	cmdStruct := command{cmdName, args}

	err = cmds.run(&s, cmdStruct)
	if err != nil {
		os.Exit(1)
	}
}
