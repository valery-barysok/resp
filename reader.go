package resp

import (
	"bufio"
	"io"
)

type Reader struct {
	protocol *Protocol
	br       *bufio.Reader
}

func NewReader(r io.Reader, protocol *Protocol) *Reader {
	return &Reader{
		protocol: protocol,
		br:       bufio.NewReader(r),
	}
}

func (reader *Reader) Read() (*Item, error) {
	return reader.protocol.nextItem(reader.br)
}
