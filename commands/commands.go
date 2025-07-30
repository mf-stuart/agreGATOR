package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mf-stuart/agreGATOR/internal/api"
	"github.com/mf-stuart/agreGATOR/internal/config"
	"github.com/mf-stuart/agreGATOR/internal/database"
	"time"
)

var Cmds = commands{make(map[string]func(*config.State, Command) error)}

type commands struct {
	commands map[string]func(*config.State, Command) error
}

type Command struct {
	Name string
	Args []string
}

func (c *commands) Run(s *config.State, cmd Command) error {
	cmdFunc, ok := c.commands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown Command: %s", cmd.Name)
	}
	return cmdFunc(s, cmd)
}

func (c *commands) Register(name string, f func(*config.State, Command) error) {
	c.commands[name] = f
}

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("no username specified")
	}
	name := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("username unregistered: %s", cmd.Args[0])
	}
	err = s.Cfg.SetUsername(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Logged in as %s\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("no username specified")
	}
	name := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("username already registered: %s", cmd.Args[0])
	}
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	_, err = s.Db.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}
	err = s.Cfg.SetUsername(name)
	if err != nil {
		return err
	}
	fmt.Printf("User registered as %s\n", name)
	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("too many arguments")
	}
	err := s.Db.Reset(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("too many arguments")
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		printString := fmt.Sprintf("* %s", user.Name)
		if s.Cfg.CurrentUsername == user.Name {
			printString += " (current)"
		}
		fmt.Println(printString)
	}
	return nil
}

func HandlerAgg(s *config.State, cmd Command) error {
	var url string
	if len(cmd.Args) == 0 {
		url = "https://www.wagslane.dev/index.xml"
	} else {
		url = cmd.Args[0]
	}

	rssFeed, err := api.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func HandlerAddFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("usage addfeed <Name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	userId := user.ID

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userId,
	}
	feed, err := s.Db.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}
	feedId := feed.ID

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		FeedID:    feedId,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("Added feed %s to user %s\n", name, user.Name)
	return nil
}

func HandlerFeeds(s *config.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("too many arguments")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.Db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* %s - %s - %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("usage follow <url>")
	}
	url := cmd.Args[0]

	feed, err := s.Db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return err
	}
	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	newRow, err := s.Db.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("%s followed by %s\n", newRow.FeedName, newRow.UserName)
	return nil
}

func HandlerUnfollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("usage unfollow <url>")
	}
	url := cmd.Args[0]
	deleteFollowParams := database.DeleteFeedFollowParams{
		Url:    url,
		UserID: user.ID,
	}

	err := s.Db.DeleteFeedFollow(context.Background(), deleteFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("%s unfollowed by %s\n", url, user.Name)
	return nil
}

func HandlerFollowing(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return errors.New("too many arguments")
	}
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUsername)
	if err != nil {
		return err
	}
	userId := user.ID
	followList, err := s.Db.GetFeedFollowsForUser(context.Background(), userId)
	if err != nil {
		return err
	}
	for _, follow := range followList {
		feed, err := s.Db.GetFeedFromId(context.Background(), follow.FeedID)
		if err != nil {
			return err
		}
		fmt.Printf("* %s\n", feed.Name)
	}
	return nil
}
