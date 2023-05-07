package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Address              common.Address // .
	Name                 string         // .
	Symbol               string         // .
	Decimals             uint8          // .
	TotalSupply          *big.Int       // .
	CirculatingSupply    *big.Int       //
	MarketCap            *big.Float     //
	Volume24h            *big.Float     //
	PriceChange24h       *big.Float     //
	Holders              uint64         //
	LargestHolders       []TokenHolder  //
	ContractCreationDate time.Time      //
	ContractCreator      common.Address //
	TokenTransfers       uint64         //
	AverageTransferValue *big.Float     //
	UniswapPriceInWETH   *big.Float     // .
	UniswapLink          string         // .
	SushiPriceInWETH     *big.Float     //
	SushiLink            string         //
}

type TokenHolder struct {
	Address common.Address
	Balance *big.Int
	Share   float64
}
