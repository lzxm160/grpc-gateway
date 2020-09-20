package tron

import "C"
import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
)

type transferCaller struct {
	endpoint   string
	ks         *keystore.KeyStore
	ksTempPath string
	acct       *keystore.Account
	amount     *big.Int
}

func NewTransferCaller(endpoint, private string) (TransferCaller, error) {
	pri, err := crypto.HexToECDSA(private)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().Unix())
	tempPath := os.TempDir() + fmt.Sprintf("%d", rand.Int())
	defer os.RemoveAll(tempPath)
	ks := keystore.NewKeyStore(tempPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.ImportECDSA(pri, tempPass)
	if err != nil {
		return nil, err
	}
	err = ks.Unlock(acc, tempPass)
	if err != nil {
		return nil, err
	}
	return &transferCaller{
		endpoint:   endpoint,
		ks:         ks,
		ksTempPath: tempPath,
		acct:       &acc,
	}, nil
}

func (t *transferCaller) Transfer(from, to string, value int64) (txid string, err error) {
	cli, err := NewTronClient(t.endpoint)
	if err != nil {
		return
	}
	defer cli.Stop()
	tx, err := cli.Conn.Transfer(from, to, value)
	if err != nil {
		return
	}
	//j := jsonpb.Marshaler{}
	//coreTransaction, err := j.MarshalToString(tx.Transaction)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(coreTransaction))
	//
	//{
	//	var protoReq core.Transaction
	//	err := jsonpb.Unmarshal(ioutil.NopCloser(strings.NewReader(coreTransaction)), &protoReq)
	//	if err != nil {
	//		fmt.Println("unmarshal", err)
	//	}
	//}

	ctrlr := transaction.NewController(cli.Conn, t.ks, t.acct, tx.Transaction)

	err = ctrlr.ExecuteTransaction()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx.GetTxid()), nil
}
