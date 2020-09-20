package tron

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

const (
	tempPass = "password"
)

type TronClient struct {
	Conn *client.GrpcClient
}

func NewTronClient(endpoint string) (*TronClient, error) {
	conn := client.NewGrpcClient(endpoint)
	err := conn.Start()
	if err != nil {
		return nil, err
	}
	return &TronClient{conn}, nil
}

func (t *TronClient) Stop() {
	t.Conn.Stop()
}
