package resp

import (
	"bufio"
	"fmt"
	"io"
)

// Writer implements buffering for an io.Writer object.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
type Writer struct {
	protocol *Protocol
	bw       *bufio.Writer
}

// NewWriter returns a new Writer with provided protocol.
func NewWriter(w io.Writer, protocol *Protocol) *Writer {
	return &Writer{
		protocol: protocol,
		bw:       bufio.NewWriter(w),
	}
}

func (resp *Writer) Write(item *Item) error {
	return resp.protocol.Write(resp.bw, item)
}

func (resp *Writer) WriteOK() error {
	return resp.protocol.WriteOK(resp.bw)
}

func (resp *Writer) WriteEmptyBulk() error {
	return resp.protocol.WriteEmptyBulk(resp.bw)
}

func (resp *Writer) WriteZero() error {
	return resp.protocol.WriteZero(resp.bw)
}

func (resp *Writer) WriteOne() error {
	return resp.protocol.WriteOne(resp.bw)
}

func (resp *Writer) WriteNil() error {
	return resp.protocol.WriteNil(resp.bw)
}

func (resp *Writer) WriteNilBulk() error {
	return resp.protocol.WriteNilBulk(resp.bw)
}

func (resp *Writer) WritePing() error {
	return resp.protocol.WritePing(resp.bw)
}

func (resp *Writer) WritePong() error {
	return resp.protocol.WritePong(resp.bw)
}

func (resp *Writer) WriteArityError(cmd string) error {
	return resp.WriteErrorFormat("ERR wrong number of arguments for '%s' command", cmd)
}

func (resp *Writer) WriteUnknownCommandError(cmd string) error {
	return resp.WriteErrorFormat("ERR unknown command '%s'", cmd)
}

func (resp *Writer) WriteNotImplementedError(cmd string) error {
	return resp.WriteErrorFormat("ERR '%s' command is not implemented", cmd)
}

func (resp *Writer) WriteSimpleString(s string) error {
	return resp.protocol.WriteSimpleString(resp.bw, s)
}

func (resp *Writer) WriteError(err error) error {
	return resp.protocol.WriteError(resp.bw, err)
}

func (resp *Writer) WriteErrorString(s string) error {
	return resp.protocol.WriteErrorString(resp.bw, s)
}

func (resp *Writer) WriteErrorFormat(format string, v ...interface{}) error {
	return resp.WriteErrorString(fmt.Sprintf(format, v...))
}

func (resp *Writer) WriteInteger(value int) error {
	return resp.protocol.WriteInteger(resp.bw, value)
}

func (resp *Writer) WriteBulkString(s []byte) error {
	return resp.protocol.WriteBulkString(resp.bw, s)
}

func (resp *Writer) WriteArray(arr []interface{}) error {
	return resp.protocol.WriteArray(resp.bw, arr)
}

func (resp *Writer) WriteCmd(cmd []byte, args ...[]byte) error {
	return resp.protocol.WriteCmd(resp.bw, cmd, args...)
}

func (resp *Writer) Flush() error {
	return resp.protocol.Flush(resp.bw)
}
