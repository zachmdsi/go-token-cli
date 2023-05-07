package functions

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/types"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func GenerateTokenProfiles(ethNodeURL string, numBlock uint64, tokens []types.Token) ([]types.Token, error) {
	fmt.Println("\nGenerating token profiles")
	var tokenProfiles []types.Token
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}
	for _, token := range tokens {
		tokenProfile, err := utils.GetTokenProfile(cl, token)
		if err != nil {
			return nil, errors.New("getTokenProfile() failed: " + err.Error())
		} else if tokenProfile.UniswapPriceInWETH == nil {
			continue
		}
		tokenProfiles = append(tokenProfiles, tokenProfile)	
	}

	fmt.Printf("Generated %d token profiles\n", len(tokenProfiles))
	for _, profile := range tokenProfiles {
		fmt.Printf("\nAddress: %s\n", profile.Address)
		fmt.Printf("Name: %s\n", profile.Name)
		fmt.Printf("Symbol: %s\n", profile.Symbol)
		fmt.Printf("Decimals: %d\n", profile.Decimals)
		fmt.Printf("Price in WETH: %.18f\n", profile.UniswapPriceInWETH)
	}

	return tokenProfiles, nil
}
