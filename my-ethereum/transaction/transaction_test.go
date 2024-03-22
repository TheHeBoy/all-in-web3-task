package transaction

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
	myethclient "my-ethereum/client"
	"os"
	"testing"
)

// 查询区块
func Test_QueryBlock(t *testing.T) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	headNum := header.Number
	fmt.Println(header.Number.String())

	blockNumber := headNum
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144
}

// 查询交易
func Test_QueryTransaction(t *testing.T) {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().String())
		fmt.Println(tx.Gas())
		fmt.Println(tx.GasPrice().Uint64())
		fmt.Println(tx.Nonce())
		fmt.Println(tx.Data())
		fmt.Println(tx.To().Hex())

		chainID, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// @change 读取发送方的地址
		if from, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println(from.Hex())
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status) // 1
	}
}

// ETH转账
func Test_ETHTransfer(t *testing.T) {
	client = myethclient.CreateClient("HTTP://127.0.0.1:7545")

	privateKey, err := crypto.HexToECDSA("70453ba30ee20ca7890200686efc07b5ec50cef46fef17ccb53015da8cfa1fe1")
	if err != nil {
		log.Fatal("HexToECDSA ", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("PendingNonceAt ", err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	chainID, err := client.ChainID(context.Background())

	toAddress := common.HexToAddress("0x4143495F0f869A4559b6df40A3C238aaFa4910Ad")
	var data []byte

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatal("SignTx ", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("SendTransaction ", err)
	}

	fmt.Println(signedTx.Hash().Hex())
}

// 裸交易 rlp对消息进行打包和解包
func Test_RawTransaction(t *testing.T) {
	// 打包
	txHash := common.HexToHash("0xbbd121c4948848c7c40a0e0c67956b81a8949d7b05b34fd6564e237f30b1529e")
	tx, _, _ := client.TransactionByHash(context.Background(), txHash)
	rawTx, _ := rlp.EncodeToBytes(tx)

	// 解包
	tx = new(types.Transaction)
	rawTxBytes, _ := hex.DecodeString(hex.EncodeToString(rawTx))
	_ = rlp.DecodeBytes(rawTxBytes, &tx)
}

/*** 生成客户端 */

var client *ethclient.Client

func setup() {
	client = myethclient.CreateClient("https://cloudflare-eth.com")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

/*** end */
