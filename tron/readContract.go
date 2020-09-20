package tron

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type readContractCaller struct {
	endpoint   string
	terminate  chan bool
	lastHeight int64
}

func NewReadContractCaller(endpoint string) (ReadContractCaller, error) {
	cli, err := NewTronClient(endpoint)
	if err != nil {
		return nil, err
	}
	defer cli.Stop()
	b, err := cli.Conn.GetNowBlock()
	if err != nil {
		return nil, err
	}
	return &readContractCaller{
		endpoint:   endpoint,
		lastHeight: b.BlockHeader.RawData.Number,
	}, nil
}

func (r *readContractCaller) ReadContract(contract, method, args string, wait time.Duration) ([][]byte, error) {
	cli, err := NewTronClient(r.endpoint)
	if err != nil {
		return nil, err
	}
	defer cli.Stop()
	if wait > 0 {
		time.Sleep(time.Second * wait)
	}
	tx, err := cli.Conn.TriggerConstantContract(
		"",
		contract,
		method,
		args,
	)
	if err != nil {
		return nil, err
	}
	return tx.GetConstantResult(), nil
}

func (r *readContractCaller) GetReceipt(txString string, wait time.Duration) (map[string]interface{}, error) {
	if !strings.EqualFold(txString[:2], "0x") {
		txString = "0x" + txString
	}
	if wait > 0 {
		time.Sleep(time.Second * wait)
	}
	cli, err := NewTronClient(r.endpoint)
	if err != nil {
		return nil, err
	}
	defer cli.Stop()
	tx, err := cli.Conn.GetTransactionByID(txString)
	if err != nil {
		return nil, err
	}
	contracts := tx.GetRawData().GetContract()
	if len(contracts) != 1 {
		return nil, fmt.Errorf("invalid contracts")
	}
	contract := contracts[0]

	info, err := cli.Conn.GetTransactionInfoByID(txString)
	if err != nil {
		return nil, err
	}
	addrResult := address.Address(info.ContractAddress).String()
	result := make(map[string]interface{})
	t := time.Unix(info.GetBlockTimeStamp()/1000, 0)
	result["txID"] = common.BytesToHexString(info.Id)
	result["block"] = info.GetBlockNumber()
	result["timestamp"] = info.GetBlockTimeStamp()
	result["date"] = t.UTC().Format(time.RFC3339)
	result["contractAddress"] = addrResult
	mes := string(info.ResMessage)
	if mes != "" {
		mes = " " + mes
	}
	result["status"] = info.Result.String() + mes
	result["receipt"] = map[string]interface{}{
		"fee":               info.GetFee(),
		"energyFee":         info.GetReceipt().GetEnergyFee(),
		"energyUsage":       info.GetReceipt().GetEnergyUsage(),
		"originEnergyUsage": info.GetReceipt().GetOriginEnergyUsage(),
		"energyUsageTotal":  info.GetReceipt().GetEnergyUsageTotal(),
		"netFee":            info.GetReceipt().GetNetFee(),
		"netUsage":          info.GetReceipt().GetNetUsage(),
	}

	result["contractName"] = contract.Type.String()

	return result, nil
}

func (r *readContractCaller) GetLogs(txString string, wait time.Duration) ([]*core.TransactionInfo_Log, error) {
	if !strings.EqualFold(txString[:2], "0x") {
		txString = "0x" + txString
	}
	if wait > 0 {
		time.Sleep(time.Second * wait)
	}
	cli, err := NewTronClient(r.endpoint)
	if err != nil {
		return nil, err
	}
	defer cli.Stop()
	info, err := cli.Conn.GetTransactionInfoByID(txString)
	if err != nil {
		return nil, err
	}
	return info.Log, nil
}

func (r *readContractCaller) Monitor(wait time.Duration) ([]*core.TransactionInfo_Log, error) {
	heightChan := make(chan int64)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case <-r.terminate:
				fmt.Println("terminate")
				r.terminate <- true
				return
			case tipHeight := <-heightChan:
				r.dealwithHeight(tipHeight)
				//fmt.Println("new height", tipHeight)
			case err := <-errChan:
				fmt.Println("report error", err)
				return
			}
		}
	}()
	r.subscribeNewBlock(heightChan, errChan, r.terminate, wait)
	return nil, nil
}

func (r *readContractCaller) dealwithHeight(hei int64) error {
	cli, err := NewTronClient(r.endpoint)
	if err != nil {
		return err
	}

	defer cli.Stop()
	ret, err := cli.Conn.GetBlockByNum(hei)
	if err != nil {
		return err
	}
	for _, tx := range ret.GetTransactions() {
		fmt.Println(hei, "tx", hex.EncodeToString(tx.GetTxid()))
		result, err := cli.Conn.GetTransactionInfoByID(hex.EncodeToString(tx.GetTxid()))
		if err != nil {
			return err
		}
		for _, log := range result.Log {
			hexAddress := "41" + hex.EncodeToString(log.Address)
			fmt.Println(hei, "contract address", hexAddress, address.HexToAddress(hexAddress).String())
			fmt.Println(hei, "data", hex.EncodeToString(log.GetData()))
			for _, topic := range log.GetTopics() {
				fmt.Println(hei, "topics", hex.EncodeToString(topic))
			}
		}
	}
	return nil
}

func (r *readContractCaller) subscribeNewBlock(heightChan chan int64, errChan chan error, unsubscribe chan bool, wait time.Duration) {
	ticker := time.NewTicker(wait * time.Second)
	for {
		select {
		case <-unsubscribe:
			unsubscribe <- true
			return
		case <-ticker.C:
			cli, err := NewTronClient(r.endpoint)
			if err != nil {
				errChan <- err
			}

			defer cli.Stop()
			info, err := cli.Conn.GetNowBlock()
			if err != nil {
				errChan <- err
			}
			heightChan <- info.BlockHeader.RawData.Number
		}
	}
}
