package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachmdsi/go-token-cli/internal/types"
	"github.com/zachmdsi/go-token-cli/internal/utils"
)

func GetBasicContractData(cl *ethclient.Client, tokenAddress common.Address) (*types.Token, error) {
	tokenContract, err := NewERC20(cl, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("\nNewERC20() failed:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}

	opts := &bind.CallOpts{}

	name, err := tokenContract.Name(opts)
	if err != nil {
		return nil, fmt.Errorf("\nName() failed:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}

	symbol, err := tokenContract.Symbol(opts)
	if err != nil {
		return nil, fmt.Errorf("\nSymbol() failed:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}

	decimals, err := tokenContract.Decimals(opts)
	if err != nil {
		return nil, fmt.Errorf("\nDecimals() failed:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}

	totalSupply, err := tokenContract.TotalSupply(opts)
	if err != nil {
		return nil, fmt.Errorf("\nTotalSupply() failed:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}
	totalSupply = utils.CalculateTotalSupply(totalSupply, decimals)

	tokenData := &types.Token{
		Address: tokenAddress,
		Name:    name,
		Symbol:  symbol,
		Decimals: decimals,
		TotalSupply: totalSupply,
	}

	return tokenData, nil
}

func NewBoundContract(cl *ethclient.Client, tokenAddress common.Address) (*bind.BoundContract, error) {
	tokenABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return nil, fmt.Errorf("\nFailed to get token ABI:\n\tToken Address: %s\n\tError: %s", tokenAddress, err.Error())
	}

	contract := bind.NewBoundContract(tokenAddress, tokenABI, cl, cl, cl)
	return contract, nil
}

func FindCreatedContracts(ethNodeURL string, numBlocks uint64) ([]string, error) {
	fmt.Println("\nSearching for created contracts")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("\nFailed to create ethclient: %s", err.Error())
	}

	blockNum, err := cl.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("\nFailed to get block number: %s", err.Error())
	}
	startBlockNum := blockNum - numBlocks

	fmt.Printf("Iterate over %d blocks from %d -> %d\n", numBlocks, startBlockNum, blockNum)
	var addresses []string
	for i := startBlockNum; i <= blockNum; i++ {
		block, err := cl.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("\nFailed to get block number: %s", err.Error())
		}

		for _, tx := range block.Transactions() {
			if tx.To() == nil {
				addresses = append(addresses, tx.Hash().Hex())
			}
		}
	}

	fmt.Printf("Found %d newly created contracts\n", len(addresses))

	return addresses, nil
}

func GetContractAddress(cl *ethclient.Client, txHash common.Hash) (common.Address, error) {
	tx, _, err := cl.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return common.Address{}, fmt.Errorf("\nFailed to convert tx hash to a Transaction: %s", err.Error())
	}

	signer := gethtypes.LatestSignerForChainID(tx.ChainId())
	from, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("\nFailed to get the signer: %s", err.Error())
	}

	// Derive the contract address from the transaction sender and nonce
	contractAddress := crypto.CreateAddress(from, tx.Nonce())
	return contractAddress, nil
}