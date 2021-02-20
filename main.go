package main

import (
  "fmt"
	"github.com/urfave/cli/v2"
	"os"
	"log"
)

func main() {
	app := &cli.App{
    Name: "data-backup",
    Usage: "backup data with complex process",
    Action: func(c *cli.Context) error {
      fmt.Println("hi")
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
