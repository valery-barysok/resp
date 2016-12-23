package resp

import (
	"bufio"
	"fmt"
	"io"
)

type Writer struct {
	protocol *Protocol
	bw       *bufio.Writer
	err      error
}

func NewWriter(w io.Writer, protocol *Protocol) *Writer {
	return &Writer{
		protocol: protocol,
		bw:       bufio.NewWriter(w),
	}
}

func (resp *Writer) WriteOK() error {
	resp.err = resp.protocol.WriteOK(resp.bw)
	return resp.err
}

func (resp *Writer) WriteEmptyBulk() error {
	resp.err = resp.protocol.WriteEmptyBulk(resp.bw)
	return resp.err
}

func (resp *Writer) WriteZero() error {
	resp.err = resp.protocol.WriteZero(resp.bw)
	return resp.err
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
	err := resp.protocol.WriteSimpleString(resp.bw, s)
	return err
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

func (resp *Writer) End() error {
	return resp.protocol.End(resp.bw)
}
