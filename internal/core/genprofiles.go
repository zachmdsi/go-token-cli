package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/core/contracts"
	"github.com/zachmdsi/go-token-cli/internal/core/dexes"
	"github.com/zachmdsi/go-token-cli/internal/types"
)

/*
	Generating a token profile will occur in several steps:

	We are already given a list of ERC20 contract addresses so we start with:

		1. Get basic contract data:

			- Name
			- Symbol
			- Total Supply
			- Decimals

			Directly interact with the contract at the given address to get this data

		2. Get DEX data:

			- Uniswap
			- Sushi
			- Curve

			Each DEX will include:

				- Timestamp
				- Price in WETH
				- Price in USDC
				- Link to trade

			Use the factory contracts for each given DEX to calculate

			Eventually, each token profile is added to a PostgreSQL db that stores historical price data
			Since this program is focused on newly minted tokens (past 7 days), the db
			is pruned of old data weekly to an archived db. Which leads to the third step:

		3. Calculate necessary data:

			- Circulating Supply
			- Market Cap
			- Volume 1H
			- Price Change 1H
			- Holders
			- Largest Holders
			- Token Transfers

			There may possibly be other data points due to the uniqueness of this data (i.e. Price Change per Tx)
			We use the 1H interval since the lifetime of these newly minted tokens is quite low so we need a bigger microscope
			This will only be possible once the db has been integrated

*/

func GenerateTokenProfiles(ethNodeURL string, numBlock uint64, erc20addresses []string) ([]*types.Token, error) {
	fmt.Println("\nGenerating token profiles")

	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("\nFailed to dial eth node: %v", err.Error())
	}

	var tokens []*types.Token
	for _, address := range erc20addresses {
		tokenAddress := common.HexToAddress(address)
		tokenContractData, err := contracts.GetBasicContractData(cl, tokenAddress)
		if err != nil {
			return nil, fmt.Errorf("\nGetBasicContractData() failed:\n\tToken Address: %s\n\tError: %v", tokenAddress, err)
		} else if tokenContractData != nil {
			tokenUniswapPriceInWETH, _, err := dexes.GetUniswapData(cl, tokenAddress)
			if err != nil {
				return nil, fmt.Errorf("\nGetUniswapData() failed:\n\tToken Address: %s\n\tError: %v", tokenAddress, err)
			} else if tokenUniswapPriceInWETH != nil {
				newToken := &types.Token{
					Address: tokenAddress,
					Name: tokenContractData.Name,
					Symbol: tokenContractData.Symbol,
					Decimals: tokenContractData.Decimals,
					TotalSupply: tokenContractData.TotalSupply,
					UniswapPriceInWETH: tokenUniswapPriceInWETH,
				}
				tokens = append(tokens, newToken)
			}
		}
	}

	for _, token := range tokens {
		fmt.Printf("\nAddress:               %s\n", token.Address)
		fmt.Printf("Name:                  %s\n", token.Name)
		fmt.Printf("Symbol:                %s\n", token.Symbol)
		fmt.Printf("Decimals:              %d\n", token.Decimals)
		fmt.Printf("Total Supply:          %s\n", token.TotalSupply)
		fmt.Printf("Uniswap Price in WETH: %.18f\n", token.UniswapPriceInWETH)
		fmt.Println()
	}

	return tokens, nil
}
