package event_log

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"my-ethereum/erc20/token"
	"os"
	"strings"
	"testing"
)

// 监听事件日志
func Test_SubscriptEvent(t *testing.T) {
	contractAddress := common.HexToAddress("0x8c14DD246462B6b4D4A7132915b1782f47eC0808")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenMetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				fmt.Printf("Log Name: Transfer\n")

				var transferEvent LogTransfer
				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
				fmt.Printf("From: %s\n", transferEvent.From.Hex())
				fmt.Printf("To: %s\n", transferEvent.To.Hex())
				fmt.Printf("Value: %s\n", transferEvent.Value.String())
			}
		}
	}
}

// LogTransfer ..
type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// 读取事件日志
func Test_ReadEvent(t *testing.T) {
	tokenAddress := common.HexToAddress("0x8c14DD246462B6b4D4A7132915b1782f47eC0808") // DAI Address
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			tokenAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenMetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Value: %s\n", transferEvent.Value.String())
		}

		fmt.Printf("\n")
	}
}

/*** 生成客户端 */

var client *ethclient.Client

func setup() {
	var err error
	client, err = ethclient.Dial("ws://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

/*** end */
