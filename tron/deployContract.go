package tron

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	"github.com/fbsobreira/gotron-sdk/pkg/contract"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

type deployContractCaller struct {
	endpoint                string
	ks                      *keystore.KeyStore
	ksTempPath              string
	acct                    *keystore.Account
	feeLimit                int64
	consumerResourcePercent int64
	originEnergyLimit       int64
	amount                  int64
	abiString               string
	args                    []interface{}
	binString               string
}

func NewDeployContractCaller(endpoint, private string) (DeployContractCaller, error) {
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
	return &deployContractCaller{
		endpoint:   endpoint,
		ks:         ks,
		ksTempPath: tempPath,
		acct:       &acc,
	}, nil
}

func (c *deployContractCaller) SetArgs(abi, bin string, args ...interface{}) DeployContractCaller {
	c.abiString = abi
	c.binString = bin
	c.args = args
	return c
}

func (c *deployContractCaller) SetFeeLimit(f int64) DeployContractCaller {
	c.feeLimit = f
	return c
}

func (c *deployContractCaller) SetConsumerResourcePercent(consumerResourcePercent int64) DeployContractCaller {
	c.consumerResourcePercent = consumerResourcePercent
	return c
}

func (c *deployContractCaller) SetOriginEnergyLimit(originEnergyLimit int64) DeployContractCaller {
	c.originEnergyLimit = originEnergyLimit
	return c
}

func (c *deployContractCaller) SetAmount(amount int64) DeployContractCaller {
	c.amount = amount
	return c
}

func (c *deployContractCaller) Deploy(contractName string) (*api.TransactionExtention, error) {
	if len(c.binString) == 0 {
		return nil, errors.New("contract data is empty")
	}
	cli, err := NewTronClient(c.endpoint)
	if err != nil {
		return nil, err
	}
	defer cli.Stop()

	tronabi, err := contract.JSONtoABI(c.abiString)
	if err != nil {
		return nil, err
	}
	ethabi, err := abi.JSON(strings.NewReader(c.abiString))
	if err != nil {
		return nil, err
	}
	tx, err := cli.Conn.DeployContractWithArguments(c.acct.Address.String(), contractName,
		&ethabi, tronabi, c.binString, c.amount, c.feeLimit, c.consumerResourcePercent, c.originEnergyLimit, c.args...)
	if err != nil {
		return nil, err
	}
	ctrlr := transaction.NewController(cli.Conn, c.ks, c.acct, tx.Transaction)

	err = ctrlr.ExecuteTransaction()
	if err != nil {
		return nil, err
	}
	return tx, nil
}
