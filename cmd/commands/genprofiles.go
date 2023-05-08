package commands

import (
	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/internal/config"
	"github.com/zachmdsi/go-token-cli/internal/core"
	"github.com/zachmdsi/go-token-cli/internal/core/contracts"
)

func GenerateProfiles() *cli.Command {
	return &cli.Command{
		Name:    "genprofiles",
		Aliases: []string{"gp"},
		Usage:   "Generates token profiles for tokens minted in the previous 1000 blocks",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "num-blocks",
				Usage: "Number of blocks to search for created contracts",
				Value: 1000,
			},
		},
		Action: func(ctx *cli.Context) error {
			conf, err := config.LoadConfig()
			if err != nil {
				panic("Failed to load config:\n\n\t" + err.Error())
			}
			numBlocks := uint64(ctx.Int64("num-blocks"))
			txs, err := contracts.FindCreatedContracts(conf.EthNodeURL, numBlocks)
			if err != nil {
				panic("Failed to create contracts:\n\n\t" + err.Error())
			}
			erc20Addresses, err := contracts.FindERC20Tokens(conf.EthNodeURL, txs)
			if err != nil {
				panic("Failed to find ERC20 tokens:\n\n\t" + err.Error())
			}
			_, err = core.GenerateTokenProfiles(conf.EthNodeURL, numBlocks, erc20Addresses)
			if err != nil {
				panic("Failed to generate token profiles:\n\n\t" + err.Error())
			}
			return nil
		},
	}
}
