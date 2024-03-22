package smart_contract

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	myethclient "my-ethereum/client"
	"my-ethereum/erc20/token"
	"os"
	"testing"
	"time"
)

// 读取合约
func Test_ReadContract(t *testing.T) {
	// 1. load contract
	address := common.HexToAddress("0x8c14DD246462B6b4D4A7132915b1782f47eC0808")
	instance, err := token.NewToken(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// 2. read contract
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
}

// 写入合约
func Test_WriteContract(t *testing.T) {
	// 1. load contract
	address := common.HexToAddress("0x8c14DD246462B6b4D4A7132915b1782f47eC0808")
	instance, err := token.NewToken(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// 2. write contract
	privateKey, err := crypto.HexToECDSA("70453ba30ee20ca7890200686efc07b5ec50cef46fef17ccb53015da8cfa1fe1")
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTipCap

	// estimate gas
	data, err := os.ReadFile("../erc20/erc20_sol_ERC20.abi")
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}
	parsed, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	encodedData, err := parsed.Pack("mint", big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}
	estimatedGas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:      fromAddress,
		To:        &address,
		Data:      encodedData,
		GasFeeCap: gasPrice,
		GasTipCap: gasTipCap,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	auth.GasLimit = estimatedGas // in units

	tx, err := instance.Mint(auth, big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}

	txHash := common.HexToHash(tx.Hash().Hex())
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	fmt.Println("contract address =", address.Hex())
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
	fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
	fmt.Println("transaction data:", hex.EncodeToString(tx.Data()))
	fmt.Println("transaction to:", tx.To().Hex())
	fmt.Println("isPending:", isPending)

	for isPending {
		time.Sleep(time.Second * 5)
		fmt.Println("pending...")
		_, isPending, _ = client.TransactionByHash(context.Background(), txHash)
	}

	result, err := instance.BalanceOf(nil, fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result) // 100
}

/*** 生成客户端 */

var client *ethclient.Client

func setup() {
	client = myethclient.CreateClient("HTTP://127.0.0.1:7545")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

/*** end */
