package main

import (
	"./magicport"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
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
var key_block cipher.BlockMode

func sendRequest(conn net.Conn, addr string) error {
	req := make([]byte, key_block.BlockSize()*4)
	dst := make([]byte, key_block.BlockSize()*2)
	src := make([]byte, key_block.BlockSize()*2)
	copy(src, []byte(addr))

	key_block.CryptBlocks(dst, src)
	base64.StdEncoding.Encode(req, dst)

	req = append(bytes.TrimRight(req, "\x00"), '\n')
	_, err := conn.Write(req)
	return err
}

func forward(conn net.Conn) {
	if r_conn, err := net.Dial(*net_type, *server_addr); err == nil {
		if sendRequest(r_conn, *port_addr) == nil {
			go io.Copy(r_conn, conn)
			io.Copy(conn, r_conn)
		}
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

	h := sha256.New()
	h.Write([]byte(*key))
	if block, err := aes.NewCipher(h.Sum(nil)); err != nil {
		fmt.Println("create cipher fail:", err)
		return
	} else {
		key_block = cipher.NewCBCEncrypter(block, magicport.AnyPortIV)
	}

	if srv, err := net.Listen(*net_type, *bind_addr); err == nil {
		acceptLoop(srv)
	} else {
		fmt.Println(err)
	}
}
