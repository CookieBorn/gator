package main

import (
	"fmt"

	"github.com/CookieBorn/gator/internal/config"
)

func main() {
	fmt.Print("Have Fun\n")
	fmt.Print(config.Working("Config ") + "\n")
}
