package contracts

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/types"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func NewERC20(cl *ethclient.Client, tokenAddress common.Address) (*types.ERC20, error) {
	contract, err := NewBoundContract(cl, tokenAddress)
	if err != nil {
		return nil, err
	}

	return &types.ERC20{
		Contract: contract,
	}, nil
}

func FindERC20Tokens(ethNodeURL string, txs []string) ([]string, error) {
	fmt.Println("\nFinding new ERC20 tokens")

	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("\nFailed to create ethclient: %s", err.Error())
	}

	// Iterate through the list of contract creation txs to get the actual contract address
	var toAddresses []string
	for _, tx := range txs {
		txHash := common.HexToHash(tx)
		contractAddress, err := GetContractAddress(cl, txHash)
		if err != nil {
			return nil, fmt.Errorf("\nGetContractAddress() failed: %s", err.Error())
		}
		toAddresses = append(toAddresses, contractAddress.Hex())
	}

	// Iterate throught the contract addresses to check its' ABI to see if it as an ERC20 token
	var erc20Addresses []string
	for _, contractAddress := range toAddresses {
		isERC20, err := IsERC20Contract(cl, common.HexToAddress(contractAddress))
		if err != nil {
			return nil, fmt.Errorf("\nIsERC20Address() failed:\n\tContract Address: %s\n\tError: %s", contractAddress, err.Error())
		}
		if isERC20 {
			erc20Addresses = append(erc20Addresses, contractAddress)
		}
	}

	fmt.Printf("Found %d new ERC20 tokens\n", len(erc20Addresses))

	return erc20Addresses, nil
}

func IsERC20Contract(cl *ethclient.Client, contractAddress common.Address) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	parsedABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return false, fmt.Errorf("\nFailed to parse ERC20ABI: %s", err.Error())
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, cl, cl, cl)

	callOpts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}

	// we panic after an error because the function was not found.
	// this is so we can recover from the runtime error.
	err = contract.Call(callOpts, nil, "totalSupply")
	if err != nil {
		panic(err)
	}

	err = contract.Call(callOpts, nil, "balanceOf", common.Address{})
	if err != nil {
		panic(err)
	}

	err = contract.Call(callOpts, nil, "allowance", common.Address{}, common.Address{})
	if err != nil {
		panic(err)
	}

	return true, nil
}

