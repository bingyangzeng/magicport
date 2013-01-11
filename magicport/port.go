package magicport

import (
	"io"
	"log"
	"net"
)

type Interface interface {
	IsBufferEnough([]byte) bool
	IsMatch(buf []byte) bool
	GetDestAddr() string
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
		log.Printf("enter inter: %v", inter)
		for !inter.IsBufferEnough(buf[:buf_size]) {
			if n, err := conn.Read(buf[buf_size:]); err != nil {
				log.Printf("%v", err)
				return
			} else {
				buf_size += n
				log.Printf("read %d", len(buf))
			}
		}
		if inter.IsMatch(buf[:buf_size]) {
			self.connectInterface(inter, conn, buf)
			return
		} else {
			log.Printf("not match")
		}
	}
}

func (self *Port) connectInterface(inter Interface, conn net.Conn, buf []byte) {
	r_conn, err := net.Dial(self.NetType, inter.GetDestAddr())
	if err == nil {
		defer r_conn.Close()
		for i := 0; i < len(buf); {
			n, err := r_conn.Write(buf[i:])
			if err != nil {
				return
			}
			i += n
		}
		go io.Copy(r_conn, conn)
		io.Copy(conn, r_conn)
	}
}
