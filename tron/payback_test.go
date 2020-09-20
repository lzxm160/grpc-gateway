package tron

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/stretchr/testify/require"
	//github.com/fbsobreirago/gotron-sdk
)

const (
	wusdcontractpayback = "TQfJkqTs8EPGRmxCrDBK286zUWaTMHgFhY"
	wbbcontractpayback  = "TLo13hxvsSJBTdk7YZzwutZt6wtAek2kno"
	contractpayback     = "THa4xniry9R3xN1ykK2eL5BAAuNTg7vjah"

	accountpayback1    = "TXVeaD62HJ2Gfk4NYATrsW1e5mt77jBaMq"
	privateKeypayback1 = "A313EF1C4C855FD0B59DEC39613BA541644E39016282C4C9D4422EA078713F9B"

	wusdBuyBackAccount        = "TXVeaD62HJ2Gfk4NYATrsW1e5mt77jBaMq"
	wusdBuyBackAccountPrivate = "A313EF1C4C855FD0B59DEC39613BA541644E39016282C4C9D4422EA078713F9B"
	wbbFreezeAccount          = "TBC5gKh7ijpDBaWLXtGXheGGNEHfuryX9v"
	wbbFreezeAccountPrivate   = "4cc0501aa7f0f9791e0fdd717a3f9a79c357cf51b6cc422d023ee2275c71c9e1"

	wbbUser        = "TXpYcao92SReYTHrM9yWi23ve9SXeg8rzA"
	wbbUserPrivate = "49bfe1ad6a12421f5d99fb3d51ca98844309311504ef22f6b92630ddf5ed7d03"
)

func TestDepolyPayback(t *testing.T) {
	require := require.New(t)

	DSTokenBin, err := ioutil.ReadFile("./abibin/buyback.bin")
	if err != nil {
		fmt.Println("buyback.bin not found")
		return
	}
	DSTokenAbi, err := ioutil.ReadFile("./abibin/buyback.abi")
	if err != nil {
		fmt.Println("buyback.abi not found")
		return
	}

	caller, err := NewDeployContractCaller(url, privateKeypayback1)
	require.NoError(err)
	tx, err := caller.SetFeeLimit(feeLimit).SetConsumerResourcePercent(curPercent).SetOriginEnergyLimit(oeLimit).SetAmount(0).SetArgs(string(DSTokenAbi), string(DSTokenBin), wusdcontractpayback, wbbcontractpayback, wusdBuyBackAccount, wbbFreezeAccount).Deploy("contractnametest")
	require.NoError(err)
	fmt.Println(hex.EncodeToString(tx.GetTxid()))

	getReceipt(t, hex.EncodeToString(tx.GetTxid()))
}

func TestBuyBack(t *testing.T) {
	//	contractBalanceOf(t, wusdcontractpayback, accountpayback1)
	//	contractBalanceOf(t, wbbcontractpayback, accountpayback1)
	//setOwner(t, wbbcontractpayback, contractpayback)
	//approve(t, wbbcontractpayback, contractpayback)
	//approve(t, wusdcontractpayback, contractpayback)
	//approve(t, wbbFreezeAccountPrivate, wbbcontractpayback, contractpayback)
	//feed(t)
	//transfer(t)
	//buyback(t)
	setOwner(t, contractpayback, accountpayback1)
}
func TestBalance(t *testing.T) {
	//contractBalanceOf(t, wusdcontractpayback, accountpayback1)
	//contractBalanceOf(t, wbbcontractpayback, accountpayback1)
	//contractBalanceOf(t, wusdcontractpayback, wusdBuyBackAccount)
	//contractBalanceOf(t, wbbcontractpayback, wbbFreezeAccount)
	//contractBalanceOf(t, wbbcontractpayback, wbbUser)
	//contractBalanceOf(t, wusdcontractpayback, wbbUser)
	contractFreezeOf(t, wbbcontractpayback, wbbFreezeAccount)
}

func TestGetPrice(t *testing.T) {
	require := require.New(t)

	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	result, err := rc.ReadContract(contractpayback, "price()", "", 10)
	require.NoError(err)
	for _, i := range result {
		fmt.Println(big.NewInt(0).SetBytes(i))
	}
}

func TestGetOwner(t *testing.T) {
	require := require.New(t)

	rc, err := NewReadContractCaller(url)
	require.NoError(err)
	result, err := rc.ReadContract(wbbcontractpayback, "owner()", "", 10)
	require.NoError(err)
	for _, i := range result {
		ii := hex.EncodeToString(i)
		fmt.Println(ii[24:])
		fmt.Println(address.HexToAddress("41" + ii[24:]).String())
	}
}

func transfer(t *testing.T) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, privateKeypayback1)
	require.NoError(err)
	param := `[{"address":"` + wbbUser + `"},{"uint256":"10000"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(wbbcontractpayback, "transfer(address,uint256)", param)
	require.NoError(err)
	fmt.Println("txid", result)
	getReceipt(t, result)
}

func feed(t *testing.T) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, privateKeypayback1)
	require.NoError(err)
	param := `[{"uint256":"100"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(contractpayback, "feedPrice(uint256)", param)
	require.NoError(err)
	fmt.Println("txid", result)
	getReceipt(t, result)
}

func buyback(t *testing.T) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, wbbUserPrivate)
	require.NoError(err)
	param := `[{"uint256":"1000"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(contractpayback, "buyBack(uint256)", param)
	require.NoError(err)
	fmt.Println("txid", result)
	getReceipt(t, result)
}

func approve(t *testing.T, pri, contract, add string) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, pri)
	require.NoError(err)
	param := `[{"address":"` + add + `"},{"uint256":"10000"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(contract, "approve(address,uint256)", param)
	require.NoError(err)
	fmt.Println("txid", result)
	getReceipt(t, result)
}

func setOwner(t *testing.T, contract, add string) {
	require := require.New(t)
	rc, err := NewExecutionContractCaller(url, privateKeypayback1)
	require.NoError(err)
	param := `[{"address":"` + add + `"}]`
	result, err := rc.SetFeeLimit(feeLimit).Execute(contract, "setOwner(address)", param)
	require.NoError(err)
	fmt.Println("txid", result)
	getReceipt(t, result)
}
