package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	// Contract Data
	Address              common.Address
	Name                 string
	Symbol               string
	Decimals             uint8
	TotalSupply          *big.Int
	ContractCreationDate time.Time
	ContractCreator      common.Address

	// Calculated Data
	CirculatingSupply    *big.Int
	MarketCap            *big.Float
	Volume1h            *big.Float
	PriceChange1h       *big.Float
	Holders              uint64
	LargestHolders       []TokenHolder
	TokenTransfers       uint64

	// DEX Data
	UniswapPriceInWETH   *big.Float
	UniswapLink          string
	SushiPriceInWETH     *big.Float
	SushiLink            string
}

type TokenHolder struct {
	Address common.Address
	Balance *big.Int
	Share   float64
}

type ERC20 struct {
	Contract *bind.BoundContract
}
