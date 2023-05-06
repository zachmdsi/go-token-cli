package createdcontracts

import (
	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/internal/config"
	"github.com/zachmdsi/go-token-cli/internal/createdcontracts"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "createdcontracts",
		Aliases: []string{"nt"},
		Usage:   "Finds txs that created a contract in the previous 1000 blocks",
		Action: func(ctx *cli.Context) error {
			conf, err := config.LoadConfig()
			if err != nil {
				panic("Failed to load config: " + err.Error())
			}
			_, err = createdcontracts.FindCreatedContracts(conf.EthNodeURL)
			if err != nil {
				panic("Failed to search for created contracts: " + err.Error())
			}
			return nil
		},
	}
}
