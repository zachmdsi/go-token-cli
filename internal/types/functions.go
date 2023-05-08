package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (t *ERC20) Name(opts *bind.CallOpts) (string, error) {
	var result []interface{}
	err := t.Contract.Call(opts, &result, "name")
	if err != nil {
		return "", fmt.Errorf("\nContract call failed to name(): %s", err.Error())
	}
	return result[0].(string), nil
}

func (t *ERC20) Symbol(opts *bind.CallOpts) (string, error) {
	var result []interface{}
	err := t.Contract.Call(opts, &result, "symbol")
	if err != nil {
		return "", fmt.Errorf("\nContract call failed to symbol(): %s", err.Error())
	}
	return result[0].(string), nil
}

func (t *ERC20) Decimals(opts *bind.CallOpts) (uint8, error) {
	var result []interface{}
	err := t.Contract.Call(opts, &result, "decimals")
	if err != nil {
		return 0, fmt.Errorf("\nContract call failed to decimals(): %s", err.Error())
	}
	return uint8(result[0].(uint8)), nil	
}

func (t *ERC20) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var result []interface{}
	err := t.Contract.Call(opts, &result, "totalSupply")
	if err != nil {
		return nil, fmt.Errorf("\nContract call failed to totalSupply(): %s", err.Error())
	}
	return result[0].(*big.Int), nil
}
