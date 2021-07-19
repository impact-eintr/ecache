package tcp_test

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
)

func TestTCP(t *testing.T) {

	key := "tcp1"
	value := "this is tcp cache test 1!"

	klen := strconv.Itoa(len(key))
	vlen := strconv.Itoa(len(value))
	test := "S" + klen + " " + vlen + " " + key + value

	serverAddr := fmt.Sprintf("127.0.0.1:%d", 6430)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(test))
	if err != nil {
		log.Fatalln(err)
		return
	}

	reply := make([]byte, 2)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(string(reply))

}
