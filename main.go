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
	"strconv"
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
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("Missing time")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFetch(s)
		if err != nil {
			return err
		}
	}
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
	defer resp.Body.Close()
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

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("Require 2 arguments name and url\n")
	}
	feedParam := database.CreateFeedParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}
	_, errm := s.db.CreateFeed(context.Background(), feedParam)
	if errm != nil {
		return errm
	}
	fmt.Print("Feed Successfully Added:\n")
	fmt.Printf("%s\n", feedParam)
	comm := command{
		name:      "url",
		arguments: cmd.arguments[1:],
	}
	err := handleFollow(s, comm, user)
	if err != nil {
		return err
	}
	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		User, err := s.db.GetUserId(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Print("Feed:\n")
		fmt.Printf("- Name: %s\n", feed.Name)
		fmt.Printf("- URL: %s\n", feed.Url)
		fmt.Printf("- Created by: %s\n", User)
	}
	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("Missing url\n")
	}
	feed, err := s.db.GetFeedName(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	followParam := database.CreateFollowFeedsParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	followFeed, err := s.db.CreateFollowFeeds(context.Background(), followParam)
	if err != nil {
		return err
	}
	fmt.Print("Follow Successfully Created\n")
	fmt.Printf("%s\n", followFeed)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFollowFeeds(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%s is following:\n", s.configStruct.Username)
	for _, follow := range follows {
		fmt.Printf("-%s\n", follow.FeedName)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("Missing url\n")
	}
	feed, err := s.db.GetFeedName(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	unfollowParam := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), unfollowParam)
	if err != nil {
		return err
	}
	return nil
}

func scrapeFetch(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	markedFetchFeedParam := database.MarkedFetchFeedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}
	err = s.db.MarkedFetchFeed(context.Background(), markedFetchFeedParam)
	if err != nil {
		fmt.Print("FetchFeedError\n")
		return err
	}
	RSSfeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", RSSfeed.Channel.Title)
	for _, item := range RSSfeed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		desc := false
		if item.Description != "" {
			desc = true
		}
		descpt := sql.NullString{
			String: item.Description,
			Valid:  desc,
		}
		pubT := sql.NullTime{
			Time:  pubTime,
			Valid: pubTime.IsZero(),
		}
		Post := database.CreatePostParams{
			ID:          int32(uuid.New().ID()),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descpt,
			PublishedAt: pubT,
			FeedID:      feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), Post)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleBrowse(s *state, cmd command) error {
	lim := int32(0)
	if len(cmd.arguments) < 1 {
		lim = int32(2)
	} else {
		lim2, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		lim = int32(lim2)
	}
	posts, err := s.db.GetPostsUser(context.Background(), lim)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		if post.Description.Valid {
			fmt.Printf("%s\n", post.Description.String)
		}
		fmt.Printf("-%s\n", post.Url)
	}
	return nil
}
