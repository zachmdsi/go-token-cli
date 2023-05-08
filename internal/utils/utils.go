package utils

import (
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetTokenDecimals(cl *ethclient.Client, tokenAddress common.Address) (uint8, error) {
	tokenABI, err := abi.JSON(strings.NewReader(ERC20ABI))
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

func CalculateTotalSupply(totalSupply *big.Int, decimals uint8) (*big.Int) {
	totalSupplyFloat := new(big.Float).SetInt(totalSupply)
	tokenDecimalsFactor := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	actualTotalSupplyFloat := new(big.Float).Quo(totalSupplyFloat, tokenDecimalsFactor)
	actualTotalSupply := new(big.Int)
	actualTotalSupply, _ = actualTotalSupplyFloat.Int(actualTotalSupply)
	return actualTotalSupply
}