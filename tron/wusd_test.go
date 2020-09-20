package tron

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	//github.com/fbsobreira/gotron-sdk
)

const (
	wusdcontract = "TBRH9XU2upbkh1uza66qHFEunSbao8LY2j"
)

func TestDepolyWusd(t *testing.T) {
	require := require.New(t)

	DSTokenBin, err := ioutil.ReadFile("./abibin/WaterBridgeExToken.bin")
	if err != nil {
		fmt.Println("WaterBridgeExToken.bin not found")
		return
	}
	DSTokenAbi, err := ioutil.ReadFile("./abibin/WaterBridgeExToken.abi")
	if err != nil {
		fmt.Println("WaterBridgeExToken.abi not found")
		return
	}

	caller, err := NewDeployContractCaller(url, privateKey1)
	require.NoError(err)
	tx, err := caller.SetFeeLimit(feeLimit).SetConsumerResourcePercent(curPercent).SetOriginEnergyLimit(oeLimit).SetAmount(1111111).SetArgs(string(DSTokenAbi), string(DSTokenBin), "symbol").Deploy("contractnamewusdtest")
	require.NoError(err)
	fmt.Println(hex.EncodeToString(tx.GetTxid()))

	getReceipt(t, hex.EncodeToString(tx.GetTxid()))
}

func TestGetWusdReceipt(t *testing.T) {
	txid := "0x26e3abcf7dbdf6301512fcfcd400cc0bcfa37fd647046dc4cd84351d076126bf"
	getReceipt(t, txid)
	getLogs(t, txid)
}

func TestMint(t *testing.T) {
	require := require.New(t)
	{
		rc, err := NewExecutionContractCaller(url, privateKey1)
		require.NoError(err)
		param := `[{"address":"THS35xNZrEAodL7DQZecQ1cF5VcXcwvssc"},{"uint256":"100"}]`
		result, err := rc.SetFeeLimit(feeLimit).SetAmount(4444444).Execute(wusdcontract,
			"mint(address,uint256)", param)
		require.NoError(err)
		fmt.Println("txid", result)
		getReceipt(t, result)
	}
}

func TestWusdBalance(t *testing.T) {
	contractBalanceOf(t, wusdcontract, account1)
	contractBalanceOf(t, wusdcontract, normalUser)
	contractBalanceOf(t, wusdcontract, normalUser2)
}

func TestTronBalance(t *testing.T) {
	balance(t, account1)
	balance(t, wusdcontract)
	balance(t, normalUser)
	balance(t, normalUser2)
}

func TestWusdApproveAndTransfer(t *testing.T) {
	require := require.New(t)
	{
		rc, err := NewExecutionContractCaller(url, privateKey1)
		require.NoError(err)
		param := `[{"address":"` + normalUser + `"},{"uint256":"100"}]`
		result, err := rc.SetFeeLimit(feeLimit).SetAmount(3333333).Execute(wusdcontract,
			"approve(address,uint256)", param)
		require.NoError(err)
		fmt.Println("txid", result)
		getReceipt(t, result)
	}

	{
		rc, err := NewExecutionContractCaller(url, normalUserPrivate)
		require.NoError(err)
		param := `[{"address":"` + account1 + `"},{"address":"` + normalUser2 + `"},{"uint256":"10"}]`
		result, err := rc.SetFeeLimit(feeLimit).SetAmount(2222222).Execute(wusdcontract, "transferFrom(address,address,uint256)", param)
		require.NoError(err)
		fmt.Println("txid", result)
	}
}

func balance(t *testing.T, from string) {
	require := require.New(t)
	r, _ := NewReadCaller(url)
	b, err := r.Balance(from)
	require.NoError(err)
	fmt.Println("balance", b)
}
