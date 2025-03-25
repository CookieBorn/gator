package main

import (
	"fmt"
	"os"

	"github.com/CookieBorn/gator/internal/config"
)

func main() {
	file, _ := config.Read()
	pointFile := &file
	stat := state{configStruct: pointFile}
	com := innit()
	com.register("login", handleLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Print("Missing arguments\n")
		os.Exit(1)
	}
	cmd := command{
		name:      args[1],
		arguments: args[2:],
	}
	err := com.run(&stat, cmd)
	if err != nil {
		os.Exit(1)
	}

}

type state struct {
	configStruct *config.Config
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
		fmt.Print("Run Error\n")
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
