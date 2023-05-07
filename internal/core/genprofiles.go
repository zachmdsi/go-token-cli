package core

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/types"
)

func GenerateTokenProfiles(ethNodeURL string, numBlock uint64, erc20addresses []string) ([]types.Token, error) {
	fmt.Println("\nGenerating token profiles")
	tokens := CreateTokens(ethNodeURL, erc20addresses)
	var tokenProfiles []types.Token
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}
	for _, token := range tokens {
		exists, err := ExistsOnADEX(cl, token)
		if err != nil {
			return nil, errors.New("ExistsOnADEX() failed: " + err.Error())
		}
		if exists {
			tokenProfile, err := GetTokenProfile(cl, token)
			if err != nil {
				return nil, errors.New("getTokenProfile() failed: " + err.Error())
			} else if tokenProfile.UniswapPriceInWETH == nil {
				continue
			}
			tokenProfiles = append(tokenProfiles, tokenProfile)
		}
	}

	fmt.Printf("Generated %d token profiles\n", len(tokenProfiles))
	for _, profile := range tokenProfiles {
		fmt.Printf("\nAddress: %s\n", profile.Address)
		fmt.Printf("Name: %s\n", profile.Name)
		fmt.Printf("Symbol: %s\n", profile.Symbol)
		fmt.Printf("Decimals: %d\n", profile.Decimals)
		fmt.Printf("Total Supply: %d\n", profile.TotalSupply)
		fmt.Printf("Uniswap Price in WETH: %.18f\n", profile.UniswapPriceInWETH)
		fmt.Printf("Sushi Price in WETH: %.18f\n", profile.SushiPriceInWETH)
		fmt.Printf("Uniswap Link: %s\n", profile.UniswapLink)
		fmt.Printf("Sushi Link: %s\n", profile.SushiLink)
	}

	return tokenProfiles, nil
}
