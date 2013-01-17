package magicport

import (
	"bytes"
	"net"
	"regexp"
)

type RegexMatchInterface struct {
	destAddr  string
	regex     *regexp.Regexp
	readUntil []byte
	readSize  int
}

func NewRegexMatchInterface(addr, regex string, until []byte, size int) *RegexMatchInterface {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil
	}
	if len(until) == 0 && size <= 0 {
		return nil
	}

	inter := new(RegexMatchInterface)
	inter.destAddr = addr
	inter.regex = r
	inter.readUntil = until
	inter.readSize = size

	return inter
}

func (self *RegexMatchInterface) IsBufferEnough(buf []byte) bool {
	if self.readSize == 0 {
		return bytes.Contains(buf, self.readUntil)
	}
	return len(buf) >= self.readSize
}

func (self *RegexMatchInterface) Match(buf []byte, net_type string) (bool, net.Conn, error) {
	is_match := func() bool {
		if self.readSize != 0 {
			return self.regex.Match(buf[:self.readSize])
		} else {
			idx := bytes.Index(buf, self.readUntil)
			if idx != -1 {
				return self.regex.Match(buf[:idx+len(self.readUntil)])
			}
		}
		return false
	}()

	if is_match {
		conn, err := net.Dial(net_type, self.destAddr)
		WriteBuf(conn, buf)
		return true, conn, err
	}
	return false, nil, nil
}
