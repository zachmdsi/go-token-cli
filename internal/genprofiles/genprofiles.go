package genprofiles

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
		tokenProfile, err := getTokenProfile(cl, token)
		if err != nil {
			return nil, errors.New("getTokenProfile() failed: " + err.Error())
		} else if tokenProfile.PriceInWETH == nil {
			continue
		}
		tokenProfiles = append(tokenProfiles, tokenProfile)	
	}

	fmt.Printf("Generated %d token profiles\n", len(tokenProfiles))
	for _, profile := range tokenProfiles {
		fmt.Printf("\nAddress: %s\n", profile.Address)
		fmt.Printf("Name: %s\n", profile.Name)
		fmt.Printf("Symbol: %s\n", profile.Symbol)
		fmt.Printf("Price in WETH: %.18f\n", profile.PriceInWETH)
	}

	return tokenProfiles, nil
}

func getTokenProfile(cl *ethclient.Client, token types.Token) (types.Token, error) {
	tokenABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return types.Token{}, err
	}

	tokenContract := bind.NewBoundContract(token.Address, tokenABI, cl, cl, cl)

	var tokenName, tokenSymbol []interface{}

	opts := &bind.CallOpts{}

	if err := tokenContract.Call(opts, &tokenName, "name"); err != nil {
		return types.Token{}, err
	}
	if err := tokenContract.Call(opts, &tokenSymbol, "symbol"); err != nil {
		return types.Token{}, err
	}

	tokenPriceInWETH, err := getTokenPriceInWETH(cl, token.Address)
	if err != nil {
		return types.Token{}, err
	}
	decimals, err := getTokenDecimals(cl, token.Address)
	if err != nil {
		return types.Token{}, err
	}

	tokenProfile := types.Token{
		Address: token.Address,
		Name: tokenName[0].(string),
		Symbol: tokenSymbol[0].(string),
		Decimals: decimals,
		PriceInWETH: tokenPriceInWETH,
	}

	return tokenProfile, nil
}

func getTokenPriceInWETH(cl *ethclient.Client, tokenAddress common.Address) (*big.Float, error) {
	factoryABI, err := abi.JSON(strings.NewReader(utils.UniswapV2FactoryABI))
	if err != nil {
		return nil, err
	}

	factory := bind.NewBoundContract(utils.UniswapFactoryAddress, factoryABI, cl, cl, cl)

	var pairAddress common.Address
	var factoryCallResult []interface{}
	err = factory.Call(&bind.CallOpts{}, &factoryCallResult, "getPair", tokenAddress, utils.WETHAddress)
	if err != nil {
		return nil, err
	}
	if len(factoryCallResult) > 0 {
		pairAddress = factoryCallResult[0].(common.Address)
	}

	pairABI, err := abi.JSON(strings.NewReader(utils.UniswapV2PairABI))
	if err != nil {
		return nil, err
	}

	pair := bind.NewBoundContract(pairAddress, pairABI, cl, cl, cl)

	var reserves [2]*big.Int
	var pairCallResult []interface{}
	err = pair.Call(&bind.CallOpts{}, &pairCallResult, "getReserves")
	if err != nil {
		return nil, err
	}
	if len(pairCallResult) > 1 {
		reserves[0] = pairCallResult[0].(*big.Int)
		reserves[1] = pairCallResult[1].(*big.Int)
	}

	var token0Address common.Address
	var pairCallResultToken0 []interface{}
	err = pair.Call(&bind.CallOpts{}, &pairCallResultToken0, "token0")
	if err != nil {
		return nil, err
	}

	var tokenReserve, wethReserve *big.Int
	if token0Address == tokenAddress {
		tokenReserve, wethReserve = reserves[0], reserves[1]
	} else {
		tokenReserve, wethReserve = reserves[1], reserves[0]
	}

	if wethReserve.Cmp(big.NewInt(0)) == 0 {
		return nil, nil
	}

	tokenDecimals, err := getTokenDecimals(cl, tokenAddress)
	if err != nil {
		return nil, err
	}

	tokenDecimalsFactor := new(big.Float).SetFloat64(math.Pow10(int(tokenDecimals)))

	normalizedTokenReserve := new(big.Float).Quo(new(big.Float).SetInt(tokenReserve), tokenDecimalsFactor)

	price := new(big.Float).Quo(normalizedTokenReserve, new(big.Float).SetInt(wethReserve))

	minPrice := 0.000000000000000001
	maxPrice := 1000.0
	if price.Cmp(big.NewFloat(0)) == 0 || price.Cmp(big.NewFloat(minPrice)) < 0 || price.Cmp(big.NewFloat(maxPrice)) > 0 {
        return nil, nil
    }

	return price, nil

}

func getTokenDecimals(cl *ethclient.Client, tokenAddress common.Address) (uint8, error) {
	tokenABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return 0, err
	}

	token := bind.NewBoundContract(tokenAddress, tokenABI, cl, cl, cl)

	var decimals uint8
	var tokenCallResult []interface{}
	err = token.Call(&bind.CallOpts{}, &tokenCallResult, "decimals")
	if err != nil {
		return 0, err
	}
	if len(tokenCallResult) > 0 {
		decimals = tokenCallResult[0].(uint8)
	}

	return decimals, nil
}

