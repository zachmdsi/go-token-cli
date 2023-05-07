package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/cmd/commands"
)

func main() {
	app := &cli.App{
		Name:     "go-token-cli",
		Usage:    "View data about newly minted tokens on the Ethereum blockchain.",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			commands.GenerateProfiles(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
