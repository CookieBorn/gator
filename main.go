package main

import (
	"context"
	"database/sql"
	"fmt"
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
	com.register("login", handleLogin)
	com.register("register", handleRegister)
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

func handleLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		err := fmt.Errorf("Login expecting one argument")
		return err
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
