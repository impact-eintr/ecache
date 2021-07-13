package tcp

import (
	"log"
	"net"

	"github.com/impact-eintr/ecache/cache"
)

type TcpHandler struct {
	cache.Cache
}

func (th *TcpHandler) Listen() {
	l, err := net.Listen("tcp", "127.0.0.1:7895")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go th.process(conn)
	}

}

func New(c cache.Cache) *TcpHandler {
	return &TcpHandler{c}
}
