package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// CreateClient 创建客户端
func CreateClient(s string) *ethclient.Client {
	client, _ := ethclient.Dial(s)
	return client
}
