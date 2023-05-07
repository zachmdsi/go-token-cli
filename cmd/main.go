package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/cmd/createdcontracts"
	"github.com/zachmdsi/go-token-cli/cmd/genprofiles"
	"github.com/zachmdsi/go-token-cli/cmd/newerc20s"
	"github.com/zachmdsi/go-token-cli/cmd/newuniswaptokens"
)

func main() {
	app := &cli.App{
		Name:     "token-cli",
		Usage:    "Query data about tokens on the Ethereum blockchain.",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			createdcontracts.Command(),
			newerc20s.Command(),
			newuniswaptokens.Command(),
			genprofiles.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
