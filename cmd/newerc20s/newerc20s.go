package newerc20s

import (
	"github.com/urfave/cli/v2"
	"github.com/zachmdsi/go-token-cli/internal/config"
	"github.com/zachmdsi/go-token-cli/internal/newerc20s"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "newerc20s",
		Aliases: []string{"nerc20"},
		Usage:   "Finds new ERC20 tokens in the previous 1000 blocks",
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
			_, err = newerc20s.FindERC20Tokens(conf.EthNodeURL, numBlocks)
			if err != nil {
				panic("Failed to search for created contracts:\n\t\t" + err.Error())
			}
			return nil
		},
	}
}
