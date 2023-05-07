package functions

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/types"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func FindNewUniswapTokens(ethNodeURL string, erc20Addresses []string, numBlocks uint64) ([]types.Token, error) {
	fmt.Println("\nFinding new Uniswap WETH pairs")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}

	var uniswapTokens []types.Token
	for _, addr := range erc20Addresses {
		tokenAddress := common.HexToAddress(addr)
		isListed, err := utils.IsTokenOnUniswap(cl, tokenAddress)
		if err != nil {
			return nil, errors.New("isTokenOnUniswap() failed: " + err.Error())
		} else if isListed {
			uniswapTokens = append(uniswapTokens, types.Token{Address: tokenAddress})
		}
	}

	fmt.Printf("Found %d new tokens on Uniswap\n", len(uniswapTokens))
	for _, token := range uniswapTokens {
		fmt.Printf("https://app.uniswap.org/#/swap?inputCurrency=%s\n", token.Address)
	}

	return uniswapTokens, nil
}

