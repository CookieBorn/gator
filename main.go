package main

import (
	"fmt"

	"github.com/CookieBorn/gator/internal/config"
)

func main() {
	fmt.Print("Have Fun\n")
	file, _ := config.Read()
	file.SetUser("robert")
}

type state struct {
	configStruct *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		err := fmt.Errorf("Login expecting one argument")
		return err
	}
	s.configStruct.SetUser(cmd.arguments[0])
	return nil
}
