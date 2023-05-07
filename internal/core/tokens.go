package core

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


func ExistsOnADEX(cl *ethclient.Client, token types.Token) (bool, error) {
	isTokenOnSushi, err := IsTokenOnSushi(cl, token.Address)
	if err != nil {
		return false, err
	}

	isTokenOnUniswap, err := IsTokenOnUniswap(cl, token.Address)
	if err != nil {
		return false, err
	}

	switch {
	case isTokenOnSushi: return true, nil
	case isTokenOnUniswap: return true, nil
	}

	return false, nil
}

func GetContractAddress(client *ethclient.Client, txHash common.Hash) (common.Address, error) {
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return common.Address{}, err
	}

	signer := gethtypes.LatestSignerForChainID(tx.ChainId())
	from, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}

	// Derive the contract address from the transaction sender and nonce
	contractAddress := crypto.CreateAddress(from, tx.Nonce())
	return contractAddress, nil
}

func IsERC20Contract(client *ethclient.Client, contractAddress common.Address) (bool, error) {
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

func GetTokenProfile(cl *ethclient.Client, token types.Token) (types.Token, error) {
	tokenABI, err := abi.JSON(strings.NewReader(utils.ERC20ABI))
	if err != nil {
		return types.Token{}, err
	}

	tokenContract := bind.NewBoundContract(token.Address, tokenABI, cl, cl, cl)

	var tokenName, tokenSymbol, tokenTotalSupply []interface{}

	opts := &bind.CallOpts{}

	if err := tokenContract.Call(opts, &tokenName, "name"); err != nil {
		return types.Token{}, err
	}
	if err := tokenContract.Call(opts, &tokenSymbol, "symbol"); err != nil {
		return types.Token{}, err
	}
	if err := tokenContract.Call(opts, &tokenTotalSupply, "totalSupply"); err != nil {
		return types.Token{}, err
	}
	decimals, err := utils.GetTokenDecimals(cl, token.Address)
	if err != nil {
		return types.Token{}, err
	}
	uniswapTokenPriceInWETH, uniswapLink, err := GetUniswapData(cl, token.Address)
	if err != nil {
		return types.Token{}, err
	}
	sushiTokenPriceInWETH, sushiLink, err := GetSushiData(cl, token.Address)
	if err != nil {
		return types.Token{}, err
	}

	fmt.Println("SUSHI PRICE IN WETH:", sushiTokenPriceInWETH)
	fmt.Println("SUSHI LINK:", sushiLink)

	tokenProfile := types.Token{
		Address:            token.Address,
		Name:               tokenName[0].(string),
		Symbol:             tokenSymbol[0].(string),
		TotalSupply:        tokenTotalSupply[0].(*big.Int),
		Decimals:           decimals,
		UniswapPriceInWETH: uniswapTokenPriceInWETH,
		SushiPriceInWETH:   sushiTokenPriceInWETH,
		UniswapLink:        uniswapLink,
		SushiLink:          sushiLink,
	}

	return tokenProfile, nil
}

func CreateTokens(ethNodeURL string, erc20Addresses []string) []types.Token {
	var tokens []types.Token
	for _, addr := range erc20Addresses {
		tokenAddress := common.HexToAddress(addr)
		token := types.Token{Address: tokenAddress}
		tokens = append(tokens, token)
	}

	return tokens
}
