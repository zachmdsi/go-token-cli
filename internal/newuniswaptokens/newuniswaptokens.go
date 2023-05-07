package newuniswaptokens

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"unicode"

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

func FindNewUniswapTokens(ethNodeURL string, numBlocks uint64) ([]Token, error) {
	fmt.Println("Finding new Uniswap WETH pairs")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, err
	}
	erc20Addresses, err := newerc20s.FindERC20Tokens(ethNodeURL, numBlocks)
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
				formattedPrice := priceInWETH.Text('f', 18)
				if !onlyZeroAndNonDigits(formattedPrice) {
					uniswapPairPrices = append(uniswapPairPrices, Token{tokenAddress, priceInWETH})
					fmt.Printf("%s | https://app.uniswap.org/#/swap?inputCurrency=%s\n", formattedPrice, addr)
				}

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

	factory := bind.NewBoundContract(utils.UniswapFactoryAddress, factoryABI, client, client, client)

	var pairAddress common.Address
	var result []interface{}
	err = factory.Call(&bind.CallOpts{}, &result, "getPair", tokenAddress, utils.WETHAddress)
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

func onlyZeroAndNonDigits(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) && r != '0' {
			return false
		}
	}
	return true
}
