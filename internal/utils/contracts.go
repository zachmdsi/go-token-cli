package utils

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetContractAddress(client *ethclient.Client, txHash common.Hash) (common.Address, error) {
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

func IsERC20Contract(client *ethclient.Client, contractAddress common.Address) (bool, error) {
	parsedABI, err := abi.JSON(strings.NewReader(ERC20ABI))
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
