package magicport

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"net"
	"regexp"
)

var addr_regex *regexp.Regexp

func init() {
	addr_regex, _ = regexp.Compile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:\\d{1,5}$")
}

type AnyPortInterface struct {
	prefix []byte
	key    []byte
}

func NewAnyPortInterface(prefix, key []byte) *AnyPortInterface {
	inter := new(AnyPortInterface)
	inter.prefix = prefix
	inter.key = key

	return inter
}

func (self *AnyPortInterface) IsBufferEnough(buf []byte) bool {
	if len(buf) >= len(self.prefix) {
		if !bytes.HasPrefix(buf, self.prefix) {
			return true
		}
		return bytes.IndexByte(buf, '\n') != -1
	}
	return false
}

func (self *AnyPortInterface) Match(buf []byte, net_type string) (bool, net.Conn, error) {
	if bytes.HasPrefix(buf, self.prefix) {
		idx := bytes.IndexByte(buf, '\n')
		req := buf[len(self.prefix):idx]
		if ok, addr := self.getRequest(req); ok {
			conn, err := net.Dial(net_type, addr)
			WriteBuf(conn, buf[idx+1:])
			return true, conn, err
		}
	}
	return false, nil, nil
}

func (self *AnyPortInterface) getRequest(req []byte) (bool, string) {
	segs := bytes.SplitN(req, []byte(" "), 2)

	if len(segs) != 2 || !addr_regex.Match(segs[0]) {
		return false, ""
	}

	h := hmac.New(sha1.New, self.key)
	h.Write(segs[0])
	digests := h.Sum(nil)

	return bytes.Equal(segs[1], digests), string(segs[0])
}
