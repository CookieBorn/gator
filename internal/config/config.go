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
	fileLocationcom := fileLocation + "/" + configFileName
	jsonFile, err := os.Open(fileLocationcom)
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

func (c *Config) SetUser(name string) error {
	c.Username = name
	fileLocation, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("File Location Errot\n")
		return err
	}
	fileLocationcom := fileLocation + "/" + configFileName
	file, err := os.Open(fileLocationcom)
	if err != nil {
		fmt.Print("Open File Error\n")
		return err
	}
	defer file.Close()
	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Print("Marshal Error\n")
		return err
	}
	fmt.Printf("%s\n", string(jsonData))
	var fileMode os.FileMode
	err = os.WriteFile(fileLocationcom, jsonData, fileMode)
	if err != nil {
		fmt.Print("Write Error\n")
		return err
	}
	return nil
}
