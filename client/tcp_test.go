package client

import (
	"fmt"
	"log"
	"testing"
)

func TestTcpClient(t *testing.T) {
	cmd := &Cmd{
		Name:  "set",
		Key:   "/test1",
		Value: []byte("this is test1"),
	}
	cli, err := New("tcp", "127.0.0.1:6430")
	if err != nil {
		log.Fatalln(err)
	}

	cli.Run(cmd)

	fmt.Println(string(cmd.Value))

}
