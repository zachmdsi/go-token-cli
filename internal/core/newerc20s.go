package core

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FindERC20Tokens(ethNodeURL string, txs []string) ([]string, error) {
	fmt.Println("\nFinding new ERC20 tokens")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}
	var toAddresses []string
	for _, tx := range txs {
		txHash := common.HexToHash(tx)
		contractAddress, err := GetContractAddress(cl, txHash)
		if err != nil {
			return nil, errors.New("getContractAddress() failed: " + err.Error())
		}
		toAddresses = append(toAddresses, contractAddress.Hex())
	}
	var erc20Addresses []string
	for _, contractAddress := range toAddresses {
		isERC20, err := IsERC20Contract(cl, common.HexToAddress(contractAddress))
		if err != nil {
			return nil, errors.New("isERC20Contract() failed: " + err.Error())
		}
		if isERC20 {
			erc20Addresses = append(erc20Addresses, contractAddress)
		}
	}

	fmt.Printf("Found %d new ERC20 tokens\n", len(erc20Addresses))

	return erc20Addresses, nil
}
