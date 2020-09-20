package tron

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/stretchr/testify/require"
	//github.com/fbsobreirago/gotron-sdk
)

func TestMonitor(t *testing.T) {
	require := require.New(t)

	r, err := NewReadContractCaller(url)
	require.NoError(err)
	r.Monitor(10)
}

func TestAddre(t *testing.T) {
	test1 := "TCirJ7fLj9DnUNW9GVWwBBfoi5Ucq7Mmwz"
	ad, err := address.Base58ToAddress(test1)
	fmt.Println(hex.EncodeToString(ad.Bytes()), ad.Hex(), err)
	ad2 := address.HexToAddress(hex.EncodeToString(ad.Bytes()))
	fmt.Println(ad2.String(), err)
}
