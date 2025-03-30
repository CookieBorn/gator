package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/CookieBorn/gator/internal/config"
	"github.com/CookieBorn/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	file, _ := config.Read()
	pointFile := &file
	com := innit()
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	stat := state{configStruct: pointFile, db: dbQueries}
	args := os.Args
	if len(args) < 2 {
		fmt.Print("Missing arguments\n")
		os.Exit(1)
	}
	cmd := command{
		name:      args[1],
		arguments: args[2:],
	}
	err = com.run(&stat, cmd)
	if err != nil {
		os.Exit(1)
	}

}

func handleLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		err := fmt.Errorf("Login expecting one argument")
		return err
	}
	_, err := s.db.GetUserName(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("User does not exit: %v", err)
	}
	s.configStruct.SetUser(cmd.arguments[0])
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		err := fmt.Errorf("Register expecting one argument")
		return err
	}
	usereParam := database.CreateUserParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	_, err := s.db.CreateUser(context.Background(), usereParam)
	if err != nil {
		fmt.Printf("Creation error: %v\n", err)
		return fmt.Errorf("User creation error")
	}
	s.configStruct.SetUser(cmd.arguments[0])
	fmt.Print("User created succesfully\n")
	fmt.Printf("ID:%v, CreatedAt:%s, UpdatedAt:%s, Name:%s\n", usereParam.ID, usereParam.CreatedAt, usereParam.UpdatedAt, usereParam.Name)
	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handleUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == s.configStruct.Username {
			fmt.Printf(" - %s (current)\n", name)
		} else {
			fmt.Printf(" - %s\n", name)
		}
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var rssFeed RSSFeed
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rssFeed, err
	}
	client := http.Client{}
	req.Header.Set("User-Agent", "gator-CB")
	resp, err := client.Do(req)
	if err != nil {
		return &rssFeed, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rssFeed, err
	}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return &rssFeed, err
	}
	html.UnescapeString(rssFeed.Channel.Title)
	html.UnescapeString(rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		html.UnescapeString(item.Title)
		html.UnescapeString(item.Description)
	}
	return &rssFeed, nil
}
