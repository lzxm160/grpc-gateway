package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tronprotocol/grpc-gateway/core"
)

func TestMarshal(t *testing.T) {
	require := require.New(t)
	tt := &core.Transaction{
		RawData:   &core.TransactionRaw{},
		Signature: [][]byte{[]byte("xx")},
		Ret:       []*core.Transaction_Result{{}},
	}
	xx, err := json.Marshal(tt)
	require.NoError(err)
	fmt.Println(string(xx))
}
