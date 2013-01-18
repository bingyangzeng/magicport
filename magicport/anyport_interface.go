package magicport

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net"
	"regexp"
)

var AnyPortIV = []byte{
	0xf4, 0x11, 0x13, 0xa1, 0x0e, 0x68, 0x05, 0x66,
	0x48, 0x29, 0x4a, 0x6b, 0x71, 0x02, 0xee, 0x9f}

var addr_regex *regexp.Regexp

func init() {
	addr_regex, _ = regexp.Compile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:\\d{1,5}$")
}

type AnyPortInterface struct {
	prefix     []byte
	key_cipher cipher.BlockMode
}

func NewAnyPortInterface(prefix, key []byte) *AnyPortInterface {
	inter := new(AnyPortInterface)
	inter.prefix = prefix

	h := sha256.New()
	h.Write(key)
	b, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		return nil
	}
	inter.key_cipher = cipher.NewCBCDecrypter(b, AnyPortIV)

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
	addr, err := self.decrypt(req)
	if err != nil {
		return false, ""
	}

	if !addr_regex.Match(addr) {
		log.Println("not match:", addr)
		return false, ""
	}

	return true, string(addr)
}

func (self *AnyPortInterface) decrypt(req []byte) ([]byte, error) {
	src := make([]byte, self.key_cipher.BlockSize()*4)
	if _, err := base64.StdEncoding.Decode(src, req); err != nil {
		return nil, err
	}

	dst := make([]byte, self.key_cipher.BlockSize()*2)
	self.key_cipher.CryptBlocks(dst, src[:self.key_cipher.BlockSize()*2])

	return bytes.TrimRight(dst, "\x00"), nil
}
