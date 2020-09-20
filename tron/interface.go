package tron

import (
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type DeployContractCaller interface {
	SetArgs(abi, bin string, args ...interface{}) DeployContractCaller
	SetFeeLimit(int64) DeployContractCaller
	SetConsumerResourcePercent(int64) DeployContractCaller
	SetOriginEnergyLimit(int64) DeployContractCaller
	// todo sdk add amount in DeployContractWithArguments
	SetAmount(int64) DeployContractCaller
	Deploy(contractName string) (*api.TransactionExtention, error)
}

type ReadContractCaller interface {
	ReadContract(contract, method, args string, wait time.Duration) ([][]byte, error)
	GetReceipt(txString string, wait time.Duration) (map[string]interface{}, error)
	GetLogs(txString string, wait time.Duration) ([]*core.TransactionInfo_Log, error)
	Monitor(wait time.Duration) ([]*core.TransactionInfo_Log, error)
}

type ExecuteContractCaller interface {
	SetFeeLimit(int64) ExecuteContractCaller
	SetAmount(int64) ExecuteContractCaller
	Execute(contract, method, args string) (txid string, err error)
	ExecuteTransaction(caller, contract, method, args string) (tx *core.Transaction, err error)
}

type TransferCaller interface {
	Transfer(from, to string, value int64) (txid string, err error)
}

type ReadCaller interface {
	Balance(from string) (value int64, err error)
}
