package functions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func FindCreatedContracts(ethNodeURL string, numBlocks uint64) ([]string, error) {
	fmt.Println("\nSearching for created contracts")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, errors.New("Failed to create ethclient: " + err.Error())
	}

	blockNum, err := cl.BlockNumber(context.Background())
	if err != nil {
		return nil, errors.New("Failed to get block number: " + err.Error())
	}
	startBlockNum := blockNum - numBlocks

	fmt.Printf("Iterate over %d blocks from %d -> %d\n", numBlocks, startBlockNum, blockNum)
	var addresses []string
	for i := startBlockNum; i <= blockNum; i++ {
		block, err := cl.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			return nil, errors.New("Failed to get block number: " + err.Error())
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
