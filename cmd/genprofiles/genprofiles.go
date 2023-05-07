package genprofiles

import (
	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/internal/config"
	"github.com/zachmdsi/go-token-cli/internal/createdcontracts"
	"github.com/zachmdsi/go-token-cli/internal/genprofiles"
	"github.com/zachmdsi/go-token-cli/internal/newerc20s"
	"github.com/zachmdsi/go-token-cli/internal/newuniswaptokens"
)

func Command() *cli.Command {
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
				panic("Failed to load config:\n\t\t" + err.Error())
			}
			numBlocks := uint64(ctx.Int64("num-blocks"))
			txs, err := createdcontracts.FindCreatedContracts(conf.EthNodeURL, numBlocks)
			if err != nil {
				panic("Failed to create contracts:\n\t\t" + err.Error())
			}
			erc20Addresses, err := newerc20s.FindERC20Tokens(conf.EthNodeURL, txs, numBlocks)
			if err != nil {
				panic("Failed to find ERC20 tokens:\n\t\t" + err.Error())
			}
			tokens, err := newuniswaptokens.FindNewUniswapTokens(conf.EthNodeURL, erc20Addresses, numBlocks)
			if err != nil {
				panic("Failed to find new Uniswap tokens:\n\t\t" + err.Error())
			}
			_, err = genprofiles.GenerateTokenProfiles(conf.EthNodeURL, numBlocks, tokens)
			if err != nil {
				panic("Failed to generate token profiles:\n\t\t" + err.Error())
			}
			return nil
		},
	}
}
