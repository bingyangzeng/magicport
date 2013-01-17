package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"net"
)

var net_type = flag.String("net", "tcp", "tcp/udp")
var bind_addr = flag.String("bind", "", "address to bind")
var server_addr = flag.String("server", "", "address of port server")
var port_addr = flag.String("addr", "", "address to port")
var key = flag.String("key", "", "server key")

func sendRequest(conn net.Conn, addr string, key []byte) {
	h := hmac.New(sha1.New, key)
	h.Write([]byte(addr))
	digests := h.Sum(nil)

	req := append(bytes.Join([][]byte{[]byte(addr), digests}, []byte(" ")), '\n')
	conn.Write(req)
}

func forward(conn net.Conn) {
	if r_conn, err := net.Dial(*net_type, *server_addr); err == nil {
		sendRequest(r_conn, *port_addr, []byte(*key))
		go io.Copy(r_conn, conn)
		io.Copy(conn, r_conn)
	}
}

func acceptLoop(srv net.Listener) {
	for {
		if conn, err := srv.Accept(); err == nil {
			go forward(conn)
		} else {
			break
		}
	}
}

func main() {
	flag.Parse()
	if srv, err := net.Listen(*net_type, *bind_addr); err == nil {
		acceptLoop(srv)
	} else {
		fmt.Println(err)
	}
}
