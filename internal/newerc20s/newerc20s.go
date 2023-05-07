package newerc20s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func FindERC20Tokens(ethNodeURL string, txs []string, numBlocks uint64) ([]string, error) {
	fmt.Println("\nFinding new ERC20 tokens")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}
	var toAddresses []string
	for _, tx := range txs {
		txHash := common.HexToHash(tx)
		contractAddress, err := getContractAddress(cl, txHash)
		if err != nil {
			return nil, errors.New("getContractAddress() failed: " + err.Error())
		}
		toAddresses = append(toAddresses, contractAddress.Hex())
	}
	var erc20Addresses []string
	for _, contractAddress := range toAddresses {
		isERC20, err := isERC20Contract(cl, common.HexToAddress(contractAddress))
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

func getContractAddress(client *ethclient.Client, txHash common.Hash) (common.Address, error) {
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return common.Address{}, err
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}

	// Derive the contract address from the transaction sender and nonce
	contractAddress := crypto.CreateAddress(from, tx.Nonce())
	return contractAddress, nil
}

func isERC20Contract(client *ethclient.Client, contractAddress common.Address) (bool, error) {
	parsedABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return false, err
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	var totalSupplyResult []interface{}
	callOpts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}
	err = contract.Call(callOpts, &totalSupplyResult, "totalSupply")
	if err != nil {
		return false, nil
	}

	return true, nil
}
