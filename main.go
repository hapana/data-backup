package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

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

func read_config() {
	var filepath = ".backup.json"

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json.Marshal(content))
}
