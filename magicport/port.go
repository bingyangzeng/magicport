package magicport

import (
	"io"
	"log"
	"net"
)

type Interface interface {
	IsBufferEnough([]byte) bool
	Match([]byte, string) (bool, net.Conn, error)
}

type Port struct {
	NetType    string // tcp or udp
	Addr       string
	MaxBufSize int
	interfaces []Interface
}

func NewPort(net, addr string) *Port {
	port := new(Port)
	port.NetType = net
	port.Addr = addr
	port.MaxBufSize = 1024

	return port
}

func (self *Port) AddInterface(inter Interface) {
	self.interfaces = append(self.interfaces, inter)
}

func (self *Port) ListenAndServe() error {
	ln, err := net.Listen(self.NetType, self.Addr)
	if err != nil {
		return err
	}

	log.Printf("port [%s]%s listening", self.NetType, self.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go self.bindInterface(conn)
	}

	return nil
}

func (self *Port) bindInterface(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, self.MaxBufSize)
	buf_size := 0
	for _, inter := range self.interfaces {
		for !inter.IsBufferEnough(buf[:buf_size]) {
			if n, err := conn.Read(buf[buf_size:]); err != nil {
				log.Printf("%v", err)
				return
			} else {
				buf_size += n
			}
		}
		is_match, r_conn, err := inter.Match(buf[:buf_size], self.NetType)
		if is_match {
			if err == nil {
				//log.Println("start to connectInterface")
				self.connectInterface(conn, r_conn, buf)
			} else {
				log.Println("Match error:", err)
			}
			return
		}
	}
	//log.Println("bind fail")
}

func (self *Port) connectInterface(conn net.Conn, r_conn net.Conn, buf []byte) {
	defer r_conn.Close()
	go io.Copy(r_conn, conn)
	io.Copy(conn, r_conn)
}

func WriteBuf(conn net.Conn, buf []byte) error {
	n := len(buf)
	for i := 0; i < n; {
		if n, err := conn.Write(buf[i:]); err != nil {
			return err
		} else {
			i += n
		}
	}

	return nil
}
