package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
			readConfig(c.String("path"))
			detectCompose(c.String("path"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Value:       "path",
				Usage:       "Path of the config file",
				DefaultText: ".backup.json",
				Aliases:     []string{"p"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig(path string) (backupConfig, error) {
	const filename = ".backup.json"
	var filepath = filepath.Join(path, filename)
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

func detectCompose(path string) (bool, error) {

	dir := filepath.Dir(path)
	var found bool

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		found, err = filepath.Match("docker-compose*", info.Name())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		return false, err
	}

	fmt.Printf("\nCompose status: %v", found)
	if found {
		return true, nil
	}

	return false, nil
}
