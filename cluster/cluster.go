package cluster

import (
	"io/ioutil"

	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
)

type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	listener, err := memberlist.Create(conf)
	if err != nil {
		return nil, err
	}

	if cluster == "" {
		cluster = addr
	}

	clu := []string{cluster}
	_, err = listener.Join(clu)
	if err != nil {
		return nil, err
	}

	circle := consistent.New()
	circle.NumberOfReplicas = 256
}
