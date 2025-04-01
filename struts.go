package main

import (
	"fmt"

	"github.com/CookieBorn/gator/internal/config"
	"github.com/CookieBorn/gator/internal/database"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

type state struct {
	configStruct *config.Config
	db           *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	comd map[string]func(*state, command) error
}

func innit() commands {
	com := make(map[string]func(*state, command) error)
	comm := commands{comd: com}
	comm.register("login", handleLogin)
	comm.register("register", handleRegister)
	comm.register("reset", handleReset)
	comm.register("users", handleUsers)
	comm.register("agg", handleAgg)
	comm.register("addfeed", middlewareLoggedIn(handleAddFeed))
	comm.register("feeds", handleFeeds)
	comm.register("following", middlewareLoggedIn(handleFollowing))
	comm.register("follow", middlewareLoggedIn(handleFollow))
	return comm
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.comd[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	fun, ok := c.comd[cmd.name]
	if !ok {
		fmt.Print("Command incorrect\n")
		err := fmt.Errorf("Command incorrect %v", 1)
		return err
	}
	err := fun(s, cmd)
	if err != nil {
		fmt.Printf("Run Error: %v\n", err)
		return err
	}
	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
