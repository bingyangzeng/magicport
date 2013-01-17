package magicport

import (
	"bytes"
	"net"
)

type RawMatchInterface struct {
	destAddr     string
	pattern      []byte
	patternStart int
	patternEnd   int
}

func NewRawMatchInterface(addr string, pattern []byte, start int) *RawMatchInterface {
	inter := new(RawMatchInterface)
	inter.destAddr = addr
	inter.pattern = pattern
	inter.patternStart = start
	inter.patternEnd = start + len(pattern)
	return inter
}

func (self *RawMatchInterface) IsBufferEnough(buf []byte) bool {
	return len(buf) >= self.patternEnd
}

func (self *RawMatchInterface) Match(buf []byte, net_type string) (bool, net.Conn, error) {
	if bytes.Equal(self.pattern, buf[self.patternStart:self.patternEnd]) {
		conn, err := net.Dial(net_type, self.destAddr)
		WriteBuf(conn, buf)
		return true, conn, err
	}
	return false, nil, nil
}
