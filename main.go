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
