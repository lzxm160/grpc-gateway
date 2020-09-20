package tron

import "C"
import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
)

type executionContractCaller struct {
	endpoint   string
	ks         *keystore.KeyStore
	ksTempPath string
	acct       *keystore.Account
	feeLimit   int64
	amount     int64
}

func NewExecutionContractCaller(endpoint, private string) (ExecuteContractCaller, error) {
	var ks *keystore.KeyStore
	var tempPath string
	var acc keystore.Account
	if private != "" {
		pri, err := crypto.HexToECDSA(private)
		if err != nil {
			return nil, err
		}
		rand.Seed(time.Now().Unix())
		tempPath = os.TempDir() + fmt.Sprintf("%d", rand.Int())
		defer os.RemoveAll(tempPath)
		ks = keystore.NewKeyStore(tempPath, keystore.StandardScryptN, keystore.StandardScryptP)
		acc, err = ks.ImportECDSA(pri, tempPass)
		if err != nil {
			return nil, err
		}
		err = ks.Unlock(acc, tempPass)
		if err != nil {
			return nil, err
		}
	}

	return &executionContractCaller{
		endpoint:   endpoint,
		ks:         ks,
		ksTempPath: tempPath,
		acct:       &acc,
	}, nil
}

func (e *executionContractCaller) SetFeeLimit(f int64) ExecuteContractCaller {
	e.feeLimit = f
	return e
}

func (e *executionContractCaller) SetAmount(amount int64) ExecuteContractCaller {
	e.amount = amount
	return e
}

func (e *executionContractCaller) ExecuteTransaction(caller, contract, method, args string) (tx *core.Transaction, err error) {
	cli, err := NewTronClient(e.endpoint)
	if err != nil {
		return
	}
	defer cli.Stop()
	fmt.Println(caller, contract, method, args)
	txs, err := cli.Conn.TriggerContract(
		caller,
		contract,
		method,
		args, e.feeLimit, e.amount, "0", 0)
	if err != nil {
		fmt.Println("trigger", err)
		return
	}
	return txs.GetTransaction(), nil
}

func (e *executionContractCaller) Execute(contract, method, args string) (txid string, err error) {
	cli, err := NewTronClient(e.endpoint)
	if err != nil {
		return
	}
	defer cli.Stop()
	tx, err := cli.Conn.TriggerContract(
		e.acct.Address.String(),
		contract,
		method,
		args, e.feeLimit, e.amount, "0", 0)
	if err != nil {
		return
	}
	ctrlr := transaction.NewController(cli.Conn, e.ks, e.acct, tx.Transaction)

	err = ctrlr.ExecuteTransaction()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx.GetTxid()), nil
}
