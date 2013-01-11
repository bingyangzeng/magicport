package magicport

type RedirectInterface struct {
	destAddr string
}

func NewRedirectInterface(addr string) *RedirectInterface {
	inter := new(RedirectInterface)
	inter.destAddr = addr
	return inter
}

func (self *RedirectInterface) GetDestAddr() string {
	return self.destAddr
}

func (self *RedirectInterface) IsBufferEnough(buf []byte) bool {
	return true
}

func (self *RedirectInterface) IsMatch(buf []byte) bool {
	return true
}
