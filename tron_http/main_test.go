package main

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/tronprotocol/grpc-gateway/tron"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/jsonpb"
	"github.com/i9/bar"
	"github.com/stretchr/testify/require"
)

const (
	//wusdcontract   = "0x058753afad571881cb73e1e0F47080C441A96D88"
	wbbcontract = "TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC"
	//taiTubcontract = "0x515d0293BA0D4fe5F2E2D08ab37201b5BBF20664"

	url         = "grpc.shasta.trongrid.io:50051"
	account1    = "TXVeaD62HJ2Gfk4NYATrsW1e5mt77jBaMq"
	privateKey1 = "A313EF1C4C855FD0B59DEC39613BA541644E39016282C4C9D4422EA078713F9B"
)

func TestMarshal(t *testing.T) {
	require := require.New(t)
	cc, err := tron.NewExecutionContractCaller(url, privateKey1)
	require.NoError(err)
	//tt := &Call_Contract{
	//	Caller:   account1,
	//	Contract: wbbcontract,
	//	Method:   "transfer(address,uint256)",
	//	Params:   `[{"address":"TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC"},{"uint256":"10"}]`,
	//}
	tx, err := cc.ExecuteTransaction(account1, wbbcontract, "transfer(address, uint256)", `[{"address":"TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC"},{"uint256":"10"}]`)
	require.NoError(err)
	//tr := `{"rawData":{"refBlockBytes":"pJs=","refBlockHash":"JUIqrdzs0PI=","expiration":"1600424577000","contract":[{"type":"TransferContract","parameter":{"typeUrl":"type.googleapis.com/protocol.TransferContract","value":"ChVBUdvVl8vPNCr2iVh330KcxGXls8wSFUEP5mbUv4BD7YoTo5kfG783c2nyrhhk"}}],"timestamp":"1600424517976"}}`
	m := &jsonpb.Marshaler{AnyResolver: bar.BetterAnyResolver}
	marshaled, err := m.MarshalToString(tx)
	fmt.Println(marshaled, err)

	r := strings.NewReader(marshaled)
	tt := &core.Transaction{
		//RawData: &core.TransactionRaw{},
		//Signature: [][]byte{[]byte("xx")},
		//Ret:       []*core.Transaction_Result{{}},
	}
	require.NoError(jsonpb.Unmarshal(r, tt))

	//xx, err := proto.Marshal(tt)
	//xx, err := jsonpb.Marshal(tt)
	m = &jsonpb.Marshaler{AnyResolver: bar.BetterAnyResolver}
	//err:=m.MarshalToString(tt)
	//require.NoError(err)
	fmt.Println(m.MarshalToString(tt))
}

//func TestRest(t *testing.T) {
//	require := require.New(t)
//	{
//		contractBalanceOf(t, wbbcontract, account1)
//		contractBalanceOf(t, wbbcontract, wbbcontract)
//	}
//	tt := &tron.Call_Contract{
//		Caller:   "TXVeaD62HJ2Gfk4NYATrsW1e5mt77jBaMq",
//		Contract: "TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC",
//		Method:   "transfer(address,uint256)",
//		Params:   `[{"address":"TFcgBNqQPqX4fXbkZtKe7gBPZhtZddzpkC"},{"uint256":"10"}]`,
//	}
//	xx, err := json.Marshal(tt)
//	require.NoError(err)
//	fmt.Println(string(xx))
//	var res []byte
//	{
//		client := resty.New()
//
//		// POST JSON string
//		// No need to set content type, if you have client level setting
//		resp, err := client.R().
//			SetHeader("Content-Type", "application/json").
//			SetBody(string(xx)).
//			//SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
//			Post("http://192.168.59.128:38080/wallet/callcontract")
//		fmt.Println("here", resp, err)
//		res = resp.Body()
//		fmt.Println("body", res)
//	}
//	tran := &Transaction{}
//
//	{
//		//sign and broadcast
//		require.NoError(json.Unmarshal(res, tran))
//		rawData, err := json.Marshal(tran.RawData)
//		require.NoError(err)
//		h256h := sha256.New()
//		h256h.Write(rawData)
//		hash := h256h.Sum(nil)
//		pri, err := crypto.HexToECDSA(privateKey1)
//		require.NoError(err)
//		signature, err := crypto.Sign(hash, pri)
//		require.NoError(err)
//		tran.Signature = append(tran.Signature, signature)
//	}
//
//	{
//		xx, err := json.Marshal(tran)
//		require.NoError(err)
//		fmt.Println(string(xx))
//		client := resty.New()
//
//		resp, err := client.R().
//			SetHeader("Content-Type", "application/json").
//			SetBody(string(xx)).
//			Post("http://192.168.59.128:38080/wallet/broadcasttransaction")
//		fmt.Println(resp, err)
//	}
//	{
//		contractBalanceOf(t, wbbcontract, account1)
//		contractBalanceOf(t, wbbcontract, wbbcontract)
//	}
//}

func contractBalanceOf(t *testing.T, contract, add string) {
	require := require.New(t)
	rc, err := tron.NewReadContractCaller(url)
	require.NoError(err)
	param := `[{"address":"` + add + `"}]`
	result, err := rc.ReadContract(contract, "balanceOf(address)", param, 10)
	require.NoError(err)
	for _, i := range result {
		fmt.Println(big.NewInt(0).SetBytes(i))
	}
}
