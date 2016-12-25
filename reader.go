package resp

import (
	"bufio"
	"io"
)

// Reader implements reading of messages from provided underlying io.Reader
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

// Read reads and returns Message if success, otherwise error.
func (reader *Reader) Read() (*Message, error) {
	return reader.protocol.Read(reader.br)
}
