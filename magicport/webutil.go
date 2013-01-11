package magicport

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"time"
)

const (
	encodingTable = "8O7P0R56AxjMniwBL@ekcWfGYFHrXEplND3z4a2hdQq1KvCsyUuZTg9JSoVmtIb$"
)

var coder = base64.NewEncoding(encodingTable)

func NewSessionId() string {
	h := md5.New()
	io.WriteString(h, time.Now().String())
	digests := h.Sum(nil)

	return coder.EncodeToString(digests)
}
