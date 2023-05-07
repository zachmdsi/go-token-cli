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
		Name:     "low-cap-token-cli",
		Usage:    "Query data about tokens on the Ethereum blockchain.",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			commands.NewUniswapTokens(),
			commands.GenerateProfiles(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
