package newuniswaptokens

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/newerc20s"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

type Token struct {
	addr        common.Address
	priceInWETH *big.Float
}

var (
	uniswapFactoryAddress = common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f") // Uniswap V2 Factory contract address
	wethAddress           = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") // Wrapped Ether (WETH) contract address
)

func FindNewUniswapTokens(ethNodeURL string) ([]Token, error) {
	fmt.Println("Finding new Uniswap WETH pairs")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, err
	}
	erc20Addresses, err := newerc20s.FindERC20Tokens(ethNodeURL)
	if err != nil {
		return nil, err
	}

	var uniswapPairPrices []Token
	for _, addr := range erc20Addresses {
		tokenAddress := common.HexToAddress(addr)
		isListed, err := isTokenOnUniswap(cl, tokenAddress)
		if err != nil {
			return nil, err
		} else {
			if isListed {
				priceInWETH, err := getTokenPriceInWETH(cl, tokenAddress)
				if err != nil {
					return nil, err
				} else if priceInWETH == nil {
					continue
				}
				uniswapPairPrices = append(uniswapPairPrices, Token{tokenAddress, priceInWETH})
				formattedPrice := priceInWETH.Text('f', 18)
				fmt.Printf("%s | https://app.uniswap.org/#/swap?inputCurrency=%s\n", formattedPrice, addr)

			}
		}
	}

	return uniswapPairPrices, nil
}

func isTokenOnUniswap(client *ethclient.Client, tokenAddress common.Address) (bool, error) {
	factoryABI, err := abi.JSON(strings.NewReader(utils.UniswapV2FactoryABI))
	if err != nil {
		return false, err
	}

	factory := bind.NewBoundContract(uniswapFactoryAddress, factoryABI, client, client, client)

	var pairAddress common.Address
	var result []interface{}
	err = factory.Call(&bind.CallOpts{}, &result, "getPair", tokenAddress, wethAddress)
	if err != nil {
		return false, err
	}
	if len(result) > 0 {
		pairAddress = result[0].(common.Address)
	}

	// If the pair address is the zero address, it means the token is not listed on Uniswap
	return pairAddress != common.Address{}, nil
}

func getTokenPriceInWETH(cl *ethclient.Client, tokenAddress common.Address) (*big.Float, error) {
	// Parse Uniswap V2 Factory ABI
	factoryABI, err := abi.JSON(strings.NewReader(utils.UniswapV2FactoryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Uniswap V2 Factory ABI: %v", err)
	}

	// Create a bound contract instance for the Uniswap V2 Factory
	factory := bind.NewBoundContract(uniswapFactoryAddress, factoryABI, cl, cl, cl)

	// Get the pair address for the given token and WETH
	var pairAddress common.Address
	var factoryCallResult []interface{}
	err = factory.Call(&bind.CallOpts{}, &factoryCallResult, "getPair", tokenAddress, wethAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get Uniswap pair address: %v", err)
	}
	if len(factoryCallResult) > 0 {
		pairAddress = factoryCallResult[0].(common.Address)
	}

	// Parse Uniswap V2 Pair ABI
	pairABI, err := abi.JSON(strings.NewReader(utils.UniswapV2PairABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Uniswap V2 Pair ABI: %v", err)
	}

	// Create a bound contract instance for the Uniswap V2 Pair
	pair := bind.NewBoundContract(pairAddress, pairABI, cl, cl, cl)

	// Get the reserves of the given token and WETH
	var reserves [2]*big.Int
	var pairCallResult []interface{}
	err = pair.Call(&bind.CallOpts{}, &pairCallResult, "getReserves")
	if err != nil {
		return nil, fmt.Errorf("failed to get reserves: %v", err)
	}
	if len(pairCallResult) > 1 {
		reserves[0] = pairCallResult[0].(*big.Int)
		reserves[1] = pairCallResult[1].(*big.Int)
	}

	// Get token0 address
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

	if price == big.NewFloat(0) {
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
