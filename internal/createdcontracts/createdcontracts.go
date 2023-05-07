package createdcontracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func FindCreatedContracts(ethNodeURL string, numBlocks uint64) ([]string, error) {
	fmt.Println("Searching for created contracts")
	cl, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, err
	}

	blockNum, err := cl.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	startBlockNum := blockNum - numBlocks

	fmt.Printf("Iterate over %d blocks from %d\n", numBlocks, startBlockNum)
	var addresses []string
	for i := startBlockNum; i <= blockNum; i++ {
		block, err := cl.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions() {
			if tx.To() == nil {
				addresses = append(addresses, tx.Hash().Hex())
			}
		}
	}
	return addresses, nil
}
