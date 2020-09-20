package main

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/jsonpb"

	"github.com/tronprotocol/grpc-gateway/core"
)

func TestMarshal(t *testing.T) {
	//require := require.New(t)
	tt := &core.Transaction{
		RawData:   &core.TransactionRaw{},
		Signature: [][]byte{[]byte("xx")},
		Ret:       []*core.Transaction_Result{{}},
	}
	//xx, err := proto.Marshal(tt)
	//xx, err := jsonpb.Marshal(tt)
	m := &jsonpb.Marshaler{}
	//err:=m.MarshalToString(tt)
	//require.NoError(err)
	fmt.Println(m.MarshalToString(tt))
}
