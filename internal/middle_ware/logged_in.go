package middle_ware

import (
	"context"
	"github.com/mf-stuart/agreGATOR/commands"
	"github.com/mf-stuart/agreGATOR/internal/config"
	"github.com/mf-stuart/agreGATOR/internal/database"
)

func LoggedIn(handler func(s *config.State, cmd commands.Command, user database.User) error) func(*config.State, commands.Command) error {
	return func(s *config.State, command commands.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, command, user)
	}
}
