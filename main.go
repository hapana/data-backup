package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"os/exec"

	"github.com/urfave/cli/v2"
)

type backupConfig struct {
	Directory     string
	DockerCompose bool
	DataDirectory string
}

var Platforms = map[string]string{
	"docker": "docker-compose*",
}

func main() {
	app := &cli.App{
		Name:  "data-backup",
		Usage: "backup data with complex process",
		Action: func(c *cli.Context) error {
			config := backupConfig{
				Directory: c.String("path"),
			}
			err := config.backup()
			if err != nil {
				log.Fatalf("Error backing up: %v", err)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Value:   "path",
				Usage:   "Path of the config file",
				Aliases: []string{"p"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *backupConfig) backup() error {
	var platformImplemented string

	log.Printf("Read config file: %v", c.Directory)
	err := c.readConfig()
	if err != nil {
		return err
	}

	for platform, filePattern := range Platforms {
		found, err := c.findFile(filePattern)
		if err != nil {
			return err
		} else {
			if found {
				platformImplemented = platform
			}
		}
	}

	err = c.processFile(filepath.Dir(c.Directory), platformImplemented)
	if err != nil {
		return err
	}

	return nil
}

func (c *backupConfig) readConfig() error {
	const filename = ".backup.json"
	var filepath = filepath.Join(c.Directory, filename)
	var config = backupConfig{}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}

	c.DataDirectory = config.Directory

	log.Printf("Config is: %+v", config)
	return nil
}

func (c *backupConfig) findFile(filePattern string) (bool, error) {
	dir := filepath.Dir(c.Directory)
	var found bool

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		found, err = filepath.Match(filePattern, info.Name())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		return false, err
	}

	if found {
		return true, nil
	}

	return false, nil
}

func (c *backupConfig) processFile(dataDir, platform string) error {
	switch platform {
	case "docker":
		cmd := exec.Command(
			"docker-compose",
			"down",
		)
		cmd.Dir = c.Directory
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		log.Printf("Output is: %s\n", stdoutStderr)
	default:
		return fmt.Errorf("No platform specified\n")
	}

	return nil
}
