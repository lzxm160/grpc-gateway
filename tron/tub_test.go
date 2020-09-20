package tron

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/stretchr/testify/require"
	//github.com/fbsobreirago/gotron-sdk
)

const (
	//wusdcontract   = "0x058753afad571881cb73e1e0F47080C441A96D88"
	wbbcontract = "TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC"
	//taiTubcontract = "0x515d0293BA0D4fe5F2E2D08ab37201b5BBF20664"

	url         = "grpc.shasta.trongrid.io:50051"
	account1    = "TXVeaD62HJ2Gfk4NYATrsW1e5mt77jBaMq"
	privateKey1 = "A313EF1C4C855FD0B59DEC39613BA541644E39016282C4C9D4422EA078713F9B"
	//wusdIssue   = "0xAf320A89E6d4743960F367f0AEb12b4205DA17Bd"
	//wusdVault   = "0x2551033F68bD83f50049A37BFC59E39fbAB762E4"
	//feeAccount  = "0xb719c9fEeA8c16bA4871F54ac3Fa041A76D47Dda"
	//
	//wusdIssuePrivate        = "cc7a6f167d3bd25324240b66039b804e23d8e5fe5a779aa279da8d52324225de"
	//wusdVaultPrivate        = "d750b63601a6715eb63279da446e586153acbf51a2fcbab4bef91b3340afbc0e"
	//feeAccountPrivate       = "dc3f2bb8c7f2cc9e85757f332c3a9bcd426ca1a6ad804ce4aff2253e03217baf"
	//wbbsystemaccountPrivate = "bb24136652f7b0ad578d2dd59d512c5139d8859ed9712b00f19d3beac2545fa1"
	//wbbsystemaccount        = "0xA8bEF1B87aAAf64c9a86910dDAB0eF860730aa7c"

	normalUser        = "THS35xNZrEAodL7DQZecQ1cF5VcXcwvssc"
	normalUserPrivate = "891330772afbbe2c32587be69765a22109f6df76ac7869619be11adc7de24ca3"
	normalUser2       = "TBC5gKh7ijpDBaWLXtGXheGGNEHfuryX9v"
	//normalUserPrivate2 = "4cc0501aa7f0f9791e0fdd717a3f9a79c357cf51b6cc422d023ee2275c71c9e1"
	feeLimit   = 1000000000
	curPercent = 100
	oeLimit    = 1000000
)

func TestInfo(t *testing.T) {
	require := require.New(t)
	conn := client.NewGrpcClient(url)
	require.NoError(conn.Start())
	defer conn.Stop()
	testacc := normalUser
	acc, err := conn.GetAccount(testacc)
	require.NoError(err)
	require.NotNil(acc)
	rewards, err := conn.GetRewardsInfo(testacc)
	require.NoError(err)

	result := make(map[string]interface{})
	result["address"] = testacc
	result["type"] = acc.GetType()
	result["balance"] = float64(acc.GetBalance()) / 1000000
	result["allowance"] = float64(acc.GetAllowance()+rewards) / 1000000
	result["rewards"] = float64(acc.GetAllowance()) / 1000000
	asJSON, _ := json.Marshal(result)
	fmt.Println(common.JSONPrettyFormat(string(asJSON)))
}

func TestDepolyWbb(t *testing.T) {
	require := require.New(t)

	DSTokenBin, err := ioutil.ReadFile("./abibin/WaterBridgeToken.bin")
	if err != nil {
		fmt.Println("WaterBridgeToken.bin not found")
		return
	}
	DSTokenAbi, err := ioutil.ReadFile("./abibin/WaterBridgeToken.abi")
	if err != nil {
		fmt.Println("WaterBridgeToken.abi not found")
		return
	}

	caller, err := NewDeployContractCaller(url, privateKey1)
	require.NoError(err)
	tx, err := caller.SetFeeLimit(feeLimit).SetConsumerResourcePercent(curPercent).SetOriginEnergyLimit(oeLimit).SetAmount(0).SetArgs(string(DSTokenAbi), string(DSTokenBin), "symbol", "wbb name", big.NewInt(2100000000)).Deploy("contractnametest")
	require.NoError(err)
	fmt.Println(hex.EncodeToString(tx.GetTxid()))

	getReceipt(t, hex.EncodeToString(tx.GetTxid()))
}

func TestGetReceipt(t *testing.T) {
	getReceipt(t, "6bb64a28d8b243664f09df90304d04fe1c48c3188197c4d9d82255ac69633191")
	getReceipt(t, "95cc392b591596ab7946aac59a9a1d375998369415571d25c0de9cdc4f6ae17a")
}

func getReceipt(t *testing.T, tx string) {
	require := require.New(t)
	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	fmt.Println("txid", tx)
	result, err := rc.GetReceipt(tx, 10)
	require.NoError(err)
	asJSON, _ := json.Marshal(result)
	fmt.Println(common.JSONPrettyFormat(string(asJSON)))
}

func getLogs(t *testing.T, tx string) {
	require := require.New(t)
	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	fmt.Println("txid", tx)
	result, err := rc.GetLogs(tx, 0)
	require.NoError(err)
	for _, log := range result {
		hexAddress := "41" + hex.EncodeToString(log.GetAddress())
		fmt.Println("address", hexAddress, address.HexToAddress(hexAddress).String())
		fmt.Println("data", hex.EncodeToString(log.GetData()))
		for _, topic := range log.GetTopics() {
			fmt.Println(hex.EncodeToString(topic))
		}
	}
}

func TestReadContractBalanceOf(t *testing.T) {
	contractBalanceOf(t, wbbcontract, account1)
	contractBalanceOf(t, wbbcontract, normalUser)
	contractBalanceOf(t, wbbcontract, normalUser2)
}

func contractBalanceOf(t *testing.T, contract, add string) {
	require := require.New(t)
	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	param := `[{"address":"` + add + `"}]`
	result, err := rc.ReadContract(contract, "balanceOf(address)", param, 10)
	require.NoError(err)
	for _, i := range result {
		fmt.Println(big.NewInt(0).SetBytes(i))
	}
}

func contractFreezeOf(t *testing.T, contract, add string) {
	require := require.New(t)
	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	param := `[{"address":"` + add + `"}]`
	result, err := rc.ReadContract(contract, "freezeOf(address)", param, 10)
	require.NoError(err)
	for _, i := range result {
		fmt.Println(big.NewInt(0).SetBytes(i))
	}
}

func TestReadName(t *testing.T) {
	require := require.New(t)

	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	result, err := rc.ReadContract(wbbcontract, "name()", "", 10)
	require.NoError(err)
	for _, i := range result {
		fmt.Println(string(i))
	}
}

func TestContractTransfer(t *testing.T) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, privateKey1)
	require.NoError(err)
	param := `[{"address":"` + normalUser + `"},{"uint256":"10"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(wbbcontract, "transfer(address,uint256)", param)
	require.NoError(err)
	fmt.Println("txid", result)
}

func TestTransfer(t *testing.T) {
	require := require.New(t)
	rc, err := NewReadCaller(url)
	require.NoError(err)
	b, err := rc.Balance(normalUser)
	require.NoError(err)
	fmt.Println(b)
	b, err = rc.Balance(wusdcontract)
	require.NoError(err)
	fmt.Println(b)
	e, err := NewTransferCaller(url, normalUserPrivate)
	//transfer to contract
	result, err := e.Transfer(normalUser, wusdcontract, 100)
	require.NoError(err)
	fmt.Println("txid", result)
	time.Sleep(10)
	b, err = rc.Balance(normalUser)
	require.NoError(err)
	fmt.Println(b)
	b, err = rc.Balance(wusdcontract)
	require.NoError(err)
	fmt.Println(b)
}

func TestApproveAndTransfer(t *testing.T) {
	require := require.New(t)
	{
		rc, err := NewExecutionContractCaller(url, privateKey1)
		require.NoError(err)
		param := `[{"address":"` + normalUser + `"},{"uint256":"100"}]`
		result, err := rc.SetFeeLimit(feeLimit).Execute(wbbcontract,
			"approve(address,uint256)", param)
		require.NoError(err)
		fmt.Println("txid", result)
		getReceipt(t, result)
	}

	{
		rc, err := NewExecutionContractCaller(url, normalUserPrivate)
		require.NoError(err)
		param := `[{"address":"` + account1 + `"},{"address":"` + normalUser2 + `"},{"uint256":"10"}]`
		result, err := rc.SetFeeLimit(feeLimit).Execute(wbbcontract, "transferFrom(address,address,uint256)", param)
		require.NoError(err)
		fmt.Println("txid", result)
	}
}
