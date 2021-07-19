package tcp

import (
	"log"
	"net"

	"github.com/impact-eintr/ecache/cache"
	"github.com/impact-eintr/ecache/global"
)

type TcpHandler struct {
	cache.Cache
}

func (th *TcpHandler) Listen() {
	l, err := net.Listen("tcp", "127.0.0.1:"+global.TcpPort)
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
