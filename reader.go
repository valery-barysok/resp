package resp

import (
	"bufio"
	"io"
)

// Reader implements reading of resp items from provided underlying io.Reader
type Reader struct {
	protocol *Protocol
	br       *bufio.Reader
}

// NewReader returns a new Reader with provided protocol.
func NewReader(r io.Reader, protocol *Protocol) *Reader {
	return &Reader{
		protocol: protocol,
		br:       bufio.NewReader(r),
	}
}

// Read reads next resp Item.
// It returns resp Item if success and error if not.
func (reader *Reader) Read() (*Item, error) {
	return reader.protocol.nextItem(reader.br)
}
