package dexes

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func GetUniswapData(cl *ethclient.Client, tokenAddress common.Address) (*big.Float, string, error) {
	factory, err := NewUniswapV2Factory(cl)
	if err != nil {
		return nil, "", fmt.Errorf("\nNewUniswapV2Factory() failed: %v", err)
	}
	pair, err := NewUniswapPairContract(cl, factory, tokenAddress)
	if err != nil {
		return nil, "", fmt.Errorf("\nNewUniswapPairContract() failed: %v", err)
	} else if pair == nil {
		return nil, "", nil
	}
	priceInWETH, err := GetTokenPriceInWETH(cl, pair, tokenAddress)	
	if err != nil {
		return nil, "", fmt.Errorf("\nNewUniswapV2Factory() failed: %v", err)
	}
	return priceInWETH, "", nil
}

func NewUniswapV2Factory(cl *ethclient.Client) (*bind.BoundContract, error) {
	parsedABI, err := abi.JSON(strings.NewReader(utils.UniswapV2FactoryABI))
	if err != nil {
		return nil, fmt.Errorf("\nFailed to parse UniswapV2FactoryABI: %v", err)
	}
	factory := bind.NewBoundContract(utils.UniswapFactoryAddress, parsedABI, cl, cl, cl)
	return factory, nil
}

func NewUniswapPairContract(cl *ethclient.Client, factory*bind.BoundContract, tokenAddress common.Address) (*bind.BoundContract, error){
	var pairAddress common.Address
	var factoryCallResult []interface{}
	err := factory.Call(&bind.CallOpts{}, &factoryCallResult, "getPair", tokenAddress, utils.WETHAddress)
	if err != nil {
		return nil, err
	}
	if len(factoryCallResult) > 0 {
		pairAddress = factoryCallResult[0].(common.Address)
	} else {
		return nil, nil
	}

	pairABI, err := abi.JSON(strings.NewReader(utils.UniswapV2PairABI))
	if err != nil {
		return nil, err
	}

	pair := bind.NewBoundContract(pairAddress, pairABI, cl, cl, cl)

	return pair, nil
}

func GetTokenPriceInWETH(cl *ethclient.Client, pair *bind.BoundContract, tokenAddress common.Address) (*big.Float, error) {
	var reserves [2]*big.Int
	var pairCallResult []interface{}
	err := pair.Call(&bind.CallOpts{}, &pairCallResult, "getReserves")
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

	tokenDecimals, err := utils.GetTokenDecimals(cl, tokenAddress)
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