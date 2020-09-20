package main

import (
	"fmt"
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
