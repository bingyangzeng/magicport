package magicport

import (
	"bytes"
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

func (self *RawMatchInterface) GetDestAddr() string {
	return self.destAddr
}

func (self *RawMatchInterface) IsBufferEnough(buf []byte) bool {
	return len(buf) >= self.patternEnd
}

func (self *RawMatchInterface) IsMatch(buf []byte) bool {
	return bytes.Equal(self.pattern, buf[self.patternStart:self.patternEnd])
}
