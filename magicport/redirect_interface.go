package magicport

import "net"

type RedirectInterface struct {
	destAddr string
}

func NewRedirectInterface(addr string) *RedirectInterface {
	inter := new(RedirectInterface)
	inter.destAddr = addr
	return inter
}

func (self *RedirectInterface) IsBufferEnough(buf []byte) bool {
	return true
}

func (self *RedirectInterface) Match(buf []byte, net_type string) (bool, net.Conn, error) {
	conn, err := net.Dial(net_type, self.destAddr)
	WriteBuf(conn, buf)
	return true, conn, err
}
