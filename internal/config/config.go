package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Working(x string) string {
	return x + "Works"
}

type Config struct {
	DB_URL   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	var config Config
	fileLocation, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("File Location Errot\n")
		return config, err
	}
	fileLocation = fileLocation + configFileName
	jsonFile, err := os.Open(configFileName)
	if err != nil {
		fmt.Print("File Open Error\n")
		return config, err
	}
	defer jsonFile.Close()
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Print("File Read Error\n")
		return config, err
	}
	json.Unmarshal(data, &config)
	return config, nil
}

func (c *Config) SetUser(name string) {
	c.Username = name
	file, err := os.Open(configFileName)
	if err != nil {
		fmt.Print("Open File Error\n")
	}
	defer file.Close()
	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Print("Marshal Error\n")
	}
	fmt.Printf("%s\n", string(jsonData))
	var fileMode os.FileMode
	err = os.WriteFile(configFileName, jsonData, fileMode)
}
