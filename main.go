package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type backupConfig struct {
	Directory string
}

func main() {
	app := &cli.App{
		Name:  "data-backup",
		Usage: "backup data with complex process",
		Action: func(c *cli.Context) error {
			read_config()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func read_config() (backupConfig, error) {
	var filepath = ".backup.json"
	var config = backupConfig{}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return backupConfig{}, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return backupConfig{}, err
	}

	fmt.Printf("Config is: %+v", config)
	return config, nil
}
