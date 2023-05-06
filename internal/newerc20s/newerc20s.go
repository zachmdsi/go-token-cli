package newerc20s

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/createdcontracts"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func FindERC20Tokens(ethNodeURL string) ([]string, error) {
	fmt.Println("Finding new ERC20 tokens")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, err
	}
	var toAddresses []string
	txs, err := createdcontracts.FindCreatedContracts(ethNodeURL)
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		txHash := common.HexToHash(tx)
		contractAddress, err := getContractAddress(cl, txHash)
		if err != nil {
			return nil, err
		}
		toAddresses = append(toAddresses, contractAddress.Hex())
	}
	var erc20Addresses []string
	for _, contractAddress := range toAddresses {
		isERC20, err := isERC20Contract(cl, common.HexToAddress(contractAddress))
		if err != nil {
			return nil, err
		}
		if isERC20 {
			fmt.Printf("https://etherscan.io/token/%s\n", contractAddress)
			erc20Addresses = append(erc20Addresses, contractAddress)
		}
	}
	return erc20Addresses, nil
}

func getContractAddress(client *ethclient.Client, txHash common.Hash) (common.Address, error) {
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get transaction: %v", err)
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get transaction sender: %v", err)
	}

	// Derive the contract address from the transaction sender and nonce
	contractAddress := crypto.CreateAddress(from, tx.Nonce())
	return contractAddress, nil
}

func isERC20Contract(client *ethclient.Client, contractAddress common.Address) (bool, error) {
	parsedABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return false, fmt.Errorf("failed to parse ERC-20 ABI: %v", err)
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	// Check for the existence of the `totalSupply` method
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
