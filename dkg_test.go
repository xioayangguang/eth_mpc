package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"eth_mpc/config"
	"eth_mpc/contracts/erc20"
	"eth_mpc/impl"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereumCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"testing"
)

func init() {
	yamlFile, err := os.ReadFile("config/node-1.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(yamlFile, &config.Cfg); err != nil {
		panic(err)
	}
}

func TestEthTx(t *testing.T) {
	client, err := ethclient.Dial("https://polygon-amoy.blockpi.network/v1/rpc/private")
	if err != nil {
		panic(err)
	}
	fromAddress, err := impl.AddressFromXY(config.Cfg.Pubkey.X, config.Cfg.Pubkey.Y)
	if err != nil {
		panic(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		panic(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	to := ethereumCommon.HexToAddress("0x037cBF2B684343b659fAc3A0aa3b5bD2f453F584")
	tx := types.NewTx(&types.DynamicFeeTx{
		//ChainID :   big.NewInt(80002),
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: big.NewInt(38694000460),
		Gas:       21000,
		To:        &to,
		Value:     big.NewInt(100000000),
	})
	signer := types.NewLondonSigner(chainID)
	sig, err := getSign(signer.Hash(tx).Bytes())
	if err != nil {
		panic(err)
	}
	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		panic(err)
	}
	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		panic(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func TestEthContracts(t *testing.T) {
	client, err := ethclient.Dial("https://rpc.thanos-sepolia.tokamak.network")
	if err != nil {
		panic(err)
	}
	fromAddress, err := impl.AddressFromXY(config.Cfg.Pubkey.X, config.Cfg.Pubkey.Y)
	if err != nil {
		panic(err)
	}
	fmt.Println(fromAddress.String())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		panic(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	erc20Address := ethereumCommon.HexToAddress("0x3Fd03027edF86B3c319f0367E2d3cD2540334a4b")
	erc20Instance, err := erc20.NewErc20(erc20Address, client)
	if err != nil {
		log.Fatal(err)
	}
	//// Pack the input, call and unpack the results
	//abi, err := erc20.Erc20MetaData.GetAbi()
	//client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	From: fromAddress,
	//	To:   &erc20Address,
	//	Data: abi.Pack(""),
	//})
	transactOpts, err := NewTransactorWithChainID(fromAddress, chainID)
	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.Value = big.NewInt(0)
	transactOpts.GasFeeCap = big.NewInt(38694000460)
	//transactOpts.GasLimit = uint64(300000)
	transactOpts.GasTipCap = gasTipCap
	tx, err := erc20Instance.AirDropToken(transactOpts)
	if err != nil {
		return
	}
	fmt.Printf("tx sent: %s", tx.Hash().Hex())
}

func TestDkg(t *testing.T) {
	dkgData, err := httpRequest(map[string]string{"type": "dkg"})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	data, err := json.Marshal(dkgData)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(string(data))
}

func TestReshare(t *testing.T) {
	reshareData, err := httpRequest(map[string]string{"type": "reshare"})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	data, err := json.Marshal(reshareData)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(string(data))
}

func getSign(data []byte) ([]byte, error) {
	hashBase64 := base64.StdEncoding.EncodeToString(data) //发送给其他节点，这个是签名所需的数据
	signData, err := httpRequest(map[string]string{"type": "signer", "msg": hashBase64})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	sig, err := base64.StdEncoding.DecodeString(cast.ToString(signData["sign"]))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	pk, err := crypto.SigToPub(data, sig)
	fmt.Println(crypto.PubkeyToAddress(*pk).String())
	if len(sig) != 65 {
		log.Fatal("invalid signature length")
	}
	return sig, nil
}

func httpRequest(req map[string]string) (map[string]interface{}, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://127.0.0.1:8083/hello", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil, err
	}
	log.Println("Response:", string(body))
	responseData := map[string]interface{}{}
	if err = json.Unmarshal(body, &responseData); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return responseData, nil
}

func NewTransactorWithChainID(keyAddr ethereumCommon.Address, chainID *big.Int) (*bind.TransactOpts, error) {
	if chainID == nil {
		return nil, bind.ErrNoChainID
	}
	signer := types.NewLondonSigner(chainID)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(address ethereumCommon.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, bind.ErrNotAuthorized
			}
			signature, err := getSign(signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Context: context.Background(),
	}, nil
}
