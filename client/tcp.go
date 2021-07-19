package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func newTcpClient(server string) (*tcpClient, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(conn)
	return &tcpClient{conn, r}, nil

}

func (c *tcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return

	} else if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		return

	} else if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
		return
	}
}

func (c *tcpClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		} else if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value)
		} else if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}

	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
}

func (c *tcpClient) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient) sendSet(key string, value []byte) {
	klen := len(key)
	vlen := len(value)
	b := append([]byte(fmt.Sprintf("S%d %d %s", klen, vlen, key)), value...) // 可以优化
	c.Write(b)
}

func (c *tcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0
	}
	return l
}

func (c *tcpClient) recvResponse() ([]byte, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return nil, nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return nil, e
		}
		return nil, errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return nil, e
	}
	return value, nil
}
