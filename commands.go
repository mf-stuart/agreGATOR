package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mf-stuart/agreGATOR/internal/database"
	"time"
)

var cmds = commands{make(map[string]func(*state, command) error)}

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
		return errors.New("no username specified")
	}
	name := cmd.args[0]

	if _, err := s.db.GetUser(context.Background(), name); err == nil {
		err := s.cfg.SetUsername(cmd.args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Logged in as %s\n", cmd.args[0])
	} else {
		return fmt.Errorf("username unregistered: %s", cmd.args[0])
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no username specified")
	}
	name := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), name); errors.Is(err, sql.ErrNoRows) {
		newUser := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name}
		_, err := s.db.CreateUser(context.Background(), newUser)
		if err != nil {
			return err
		}
		err = s.cfg.SetUsername(name)
		if err != nil {
			return err
		}
		fmt.Printf("User registered as %s\n", name)
	} else {
		return fmt.Errorf("username already registered: %s", cmd.args[0])
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		printString := fmt.Sprintf("* %s", user.Name)
		if s.cfg.CurrentUsername == user.Name {
			printString += " (current)"
		}
		fmt.Println(printString)
	}
	return nil
}
