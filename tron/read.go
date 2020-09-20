package tron

import "C"

type readCaller struct {
	endpoint string
}

func NewReadCaller(endpoint string) (ReadCaller, error) {
	return &readCaller{endpoint: endpoint}, nil
}

func (t *readCaller) Balance(acc string) (value int64, err error) {
	cli, err := NewTronClient(t.endpoint)
	if err != nil {
		return
	}
	defer cli.Stop()
	a, err := cli.Conn.GetAccount(acc)
	if err != nil {
		return
	}
	value = a.Balance
	return
}
