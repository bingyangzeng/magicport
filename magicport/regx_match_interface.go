package magicport

import (
	"bytes"
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

func (self *RegexMatchInterface) GetDestAddr() string {
	return self.destAddr
}

func (self *RegexMatchInterface) IsBufferEnough(buf []byte) bool {
	if self.readSize == 0 {
		return bytes.Contains(buf, self.readUntil)
	}
	return len(buf) >= self.readSize
}

func (self *RegexMatchInterface) IsMatch(buf []byte) bool {
	if self.readSize != 0 {
		return self.regex.Match(buf[:self.readSize])
	} else {
		idx := bytes.Index(buf, self.readUntil)
		if idx != -1 {
			return self.regex.Match(buf[:idx+len(self.readUntil)])
		}
	}
	return false
}
