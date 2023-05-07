package utils

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func IsTokenOnUniswap(client *ethclient.Client, tokenAddress common.Address) (bool, error) {
	factoryABI, err := abi.JSON(strings.NewReader(UniswapV2FactoryABI))
	if err != nil {
		return false, err
	}

	factory := bind.NewBoundContract(UniswapFactoryAddress, factoryABI, client, client, client)

	var pairAddress common.Address
	var result []interface{}
	err = factory.Call(&bind.CallOpts{}, &result, "getPair", tokenAddress, WETHAddress)
	if err != nil {
		return false, err
	}
	if len(result) > 0 {
		pairAddress = result[0].(common.Address)
	}

	// If the pair address is the zero address, it means the token is not listed on Uniswap
	return pairAddress != common.Address{}, nil
}
