package client

import "errors"

type Cmd struct {
	Name  string
	Key   string
	Value []byte
	Error error
}

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

func New(typ, server string) (Client, error) {
	if typ == "http" {
		return newHttpClient(server), nil
	} else if typ == "tcp" {
		return newTcpClient(server)
	} else {
		return nil, errors.New("不支持的服务类型")
	}
}
