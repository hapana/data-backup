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
			_ = backup(c.String("path"))
			readConfig(c.String("path"))
			findFile(c.String("path"), "docker-compose*")
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

func backup(path string) error {
	_, err := readConfig(path)
	if err != nil {
		log.Fatalf("Can't read error: %v", err)
	}

	_, err = findFile(path, "docker-compose*")
	if err != nil {
		log.Fatalf("Can't detect compose: %v", err)
	}
	return nil
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

func findFile(path, filePattern string) (bool, error) {

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
