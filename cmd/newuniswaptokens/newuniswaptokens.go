package newuniswaptokens

import (
	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/internal/config"
	"github.com/zachmdsi/go-token-cli/internal/newuniswaptokens"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "newuniswaptokens",
		Aliases: []string{"nut"},
		Usage:   "Finds new WETH/ERC20 pairs on Uniswap from the previous 1000 blocks",
		Action: func(ctx *cli.Context) error {
			conf, err := config.LoadConfig()
			if err != nil {
				panic("Failed to load config: " + err.Error())
			}
			_, err = newuniswaptokens.FindNewUniswapTokens(conf.EthNodeURL)
			if err != nil {
				panic("Failed to search for created contracts: " + err.Error())
			}
			return nil
		},
	}
}
