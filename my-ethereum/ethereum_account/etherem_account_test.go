package etherem_account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	myethclient "my-ethereum/client"
	"my-ethereum/erc20/token"
	"os"
	"regexp"
	"testing"
)

// 账户余额
func Test_AccountBalance(t *testing.T) {
	account := common.HexToAddress("0x59608B282D1758BdA41CeD2e439DF78Ea92e3560")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance) // 25893180161173005034
}

// 得到账户代币余额
func Test_AccountTokenBalance(t *testing.T) {
	tokenAddress := common.HexToAddress("0x8c14DD246462B6b4D4A7132915b1782f47eC0808")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal("NewToken", err)
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, tokenAddress)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)
	fmt.Printf("symbol: %s\n", symbol)
	fmt.Printf("decimals: %v\n", decimals)
	fmt.Printf("wei: %s\n", bal)

}

// 生成新钱包
func Test_CreateNewWallet(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("private", hexutil.Encode(privateKeyBytes), "->", hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("public", hexutil.Encode(publicKeyBytes), "->", hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("public address", address)

	// crypto.PubkeyToAddress 逻辑
	//hash := sha3.NewLegacyKeccak256()
	//hash.Write(publicKeyBytes[1:])
	//fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}

// 密钥库生成
func Test_KeyStoresCreate(t *testing.T) {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "123456"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())
}

// 密钥库导入，得到私钥
func Test_KeyStoresImport(t *testing.T) {
	file := "./wallets/UTC--2024-03-20T07-43-57.332355000Z--609ca7a998d1692b5ed8409f64cd4e995ed75dfc"
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := "123456"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}

// todo分层确认性钱包
func Test_HDWallet(t *testing.T) {

}

// 地址验证
func Test_CheckAddress(t *testing.T) {
	// 正则校验地址格式
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	fmt.Printf("is valid: %v\n", re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d")) // is valid: true
	fmt.Printf("is valid: %v\n", re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d")) // is valid: false

	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	// 检查地址是否为账户或智能合约
	// 当地址上没有字节码时，我们知道它不是一个智能合约，它是一个标准的以太坊账户
	address := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	isContract := len(bytecode) > 0
	fmt.Printf("is contract: %v\n", isContract) // is contract: true

	// a random user account address
	address = common.HexToAddress("0x8e215d06ea7ec1fdb4fc5fd21768f4b34ee92ef4")
	bytecode, err = client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	isContract = len(bytecode) > 0
	fmt.Printf("is contract: %v\n", isContract) // is contract: false
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
