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

	"github.com/otiai10/copy"
	"github.com/urfave/cli/v2"
)

type backupConfig struct {
	Directories   directory
	DockerCompose bool
}

type directory struct {
	Root   string
	Data   string
	Backup string
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
				Directories: directory{
					Root: c.String("path"),
				},
			}
			err := config.backup()
			if err != nil {
				log.Fatalf("Error backing up: %v", err)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Value:       "path",
				Usage:       "Path of the config file",
				DefaultText: ".",
				Aliases:     []string{"p"},
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

	log.Printf("Read config file: %v", c.Directories.Data)
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

	err = c.processFile(filepath.Dir(c.Directories.Data), platformImplemented)
	if err != nil {
		return err
	}

	return nil
}

func (c *backupConfig) readConfig() error {
	const filename = ".backup.json"
	var filepath = filepath.Join(c.Directories.Root, filename)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	log.Printf("Config is: %+v", c)
	return nil
}

func (c *backupConfig) findFile(filePattern string) (bool, error) {
	dir := filepath.Dir(c.Directories.Root)
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
		log.Println("Pulling down docker-compose")
		cmd := exec.Command(
			"docker-compose",
			"down",
		)
		cmd.Dir = c.Directories.Root
		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		log.Printf("Copying data from %s to %s", c.Directories.Data, c.Directories.Backup)
		err = copy.Copy(filepath.Join(c.Directories.Root, c.Directories.Data), filepath.Join(c.Directories.Root, c.Directories.Backup))
		if err != nil {
			return err
		}

		log.Println("Starting docker-compose again")
		cmd = exec.Command(
			"docker-compose",
			"up",
			"-d",
		)
		cmd.Dir = c.Directories.Root
		_, err = cmd.CombinedOutput()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No platform specified\n")
	}

	return nil
}
