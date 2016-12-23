package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
)

const (
	simpleStringType = '+'
	errorType        = '-'
	integerType      = ':'
	bulkStringType   = '$'
	arrayType        = '*'
)

const bulkStringMaxLength = 512 * 1024 * 1024

const lf = '\n'

var (
	crlfSlice           = []byte("\r\n")
	okSlice             = []byte("+OK\r\n")
	errSlice            = []byte("-ERR\r\n")
	emptyBulkSlice      = []byte("$0\r\n\r\n")
	zeroSlice           = []byte(":0\r\n")
	oneSlice            = []byte(":1\r\n")
	nilSlice            = []byte(":-1\r\n")
	nilBulkSlice        = []byte("$-1\r\n")
	nilMultiBulkSlice   = []byte("*-1\r\n")
	emptyMultiBulkSlice = []byte("*0\r\n")
	pingSlice           = []byte("+PING\r\n")
	pongSlice           = []byte("+PONG\r\n")
)

var errInlineCommandNotImplemented = errors.New("ERR Inline command is not implemented")
var errBulkStringIsTooLarge = errors.New("ERR Bulk string is too large")
var errInvalidFormat = errors.New("ERR invalid format")

type Protocol struct {
	rlogger *bytes.Buffer
	wlogger *bytes.Buffer
}

func NewProtocol() *Protocol {
	return &Protocol{}
}

func NewProtocolWithLogging(out io.Writer) *Protocol {
	return &Protocol{
		rlogger: &bytes.Buffer{},
		wlogger: &bytes.Buffer{},
	}
}

func (protocol *Protocol) WriteOK(bw *bufio.Writer) error {
	return protocol.write(bw, okSlice)
}

func (protocol *Protocol) WriteEmptyBulk(bw *bufio.Writer) error {
	return protocol.write(bw, emptyBulkSlice)
}

func (protocol *Protocol) WriteZero(bw *bufio.Writer) error {
	return protocol.write(bw, zeroSlice)
}

func (protocol *Protocol) WriteOne(bw *bufio.Writer) error {
	return protocol.write(bw, oneSlice)
}

func (protocol *Protocol) WriteNil(bw *bufio.Writer) error {
	return protocol.write(bw, nilSlice)
}

func (protocol *Protocol) WriteNilBulk(bw *bufio.Writer) error {
	return protocol.write(bw, nilBulkSlice)
}

func (protocol *Protocol) WriteNilMultiBulk(bw *bufio.Writer) error {
	return protocol.write(bw, nilMultiBulkSlice)
}

func (protocol *Protocol) WriteEmptyMultiBulk(bw *bufio.Writer) error {
	return protocol.write(bw, emptyMultiBulkSlice)
}

func (protocol *Protocol) WritePing(bw *bufio.Writer) error {
	return protocol.write(bw, pingSlice)
}

func (protocol *Protocol) WritePong(bw *bufio.Writer) error {
	return protocol.write(bw, pongSlice)
}

func (protocol *Protocol) WriteSimpleString(bw *bufio.Writer, s string) error {
	if err := protocol.writeByte(bw, simpleStringType); err != nil {
		return err
	}
	if err := protocol.writeString(bw, s); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}
	return nil
}

func (protocol *Protocol) WriteError(bw *bufio.Writer, err error) error {
	return protocol.WriteErrorString(bw, err.Error())
}

func (protocol *Protocol) WriteErrorString(bw *bufio.Writer, s string) error {
	if err := protocol.writeByte(bw, errorType); err != nil {
		return err
	}
	if err := protocol.writeString(bw, s); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}
	return nil
}

func (protocol *Protocol) WriteInteger(bw *bufio.Writer, value int) error {
	if err := protocol.writeByte(bw, integerType); err != nil {
		return err
	}
	if err := protocol.writeString(bw, strconv.Itoa(value)); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}
	return nil
}

func (protocol *Protocol) WriteBulkString(bw *bufio.Writer, s []byte) error {
	if err := protocol.writeByte(bw, bulkStringType); err != nil {
		return err
	}
	if err := protocol.writeString(bw, strconv.Itoa(len(s))); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}
	if err := protocol.write(bw, s); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}
	return nil
}

func (protocol *Protocol) Write(bw *bufio.Writer, item *Item) error {
	return protocol.writeRaw(bw, item)
}

func (protocol *Protocol) writeRaw(bw *bufio.Writer, value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return protocol.WriteBulkString(bw, v)

	case string:
		return protocol.WriteSimpleString(bw, v)

	case error:
		return protocol.WriteError(bw, v)

	case int:
		return protocol.WriteInteger(bw, v)

	case []interface{}:
		return protocol.WriteArray(bw, v)

	case nil:
		return protocol.WriteNil(bw)

	case *Item:
		return protocol.writeRaw(bw, v.Raw())

	default:
		return errInvalidFormat
	}
}

func (protocol *Protocol) WriteArray(bw *bufio.Writer, arr []interface{}) error {
	if err := protocol.writeByte(bw, arrayType); err != nil {
		return err
	}

	if err := protocol.writeString(bw, strconv.Itoa(len(arr))); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}

	for i := range arr {
		err := protocol.writeRaw(bw, arr[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (protocol *Protocol) WriteCmd(bw *bufio.Writer, cmd []byte, args ...[]byte) error {
	if err := protocol.writeByte(bw, arrayType); err != nil {
		return err
	}

	if err := protocol.writeString(bw, strconv.Itoa(len(args)+1)); err != nil {
		return err
	}
	if err := protocol.writeNewLine(bw); err != nil {
		return err
	}

	err := protocol.writeRaw(bw, cmd)
	if err != nil {
		return err
	}

	for i := range args {
		err := protocol.writeRaw(bw, args[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (protocol *Protocol) writeNewLine(bw *bufio.Writer) error {
	err := protocol.write(bw, crlfSlice)
	if err != nil {
		return err
	}

	if protocol.wlogger != nil && protocol.wlogger.Len() > 0 {
		log.Printf("W: %s", string(protocol.wlogger.Bytes()))
		protocol.wlogger.Reset()
	}

	return nil
}

func (protocol *Protocol) write(bw *bufio.Writer, b []byte) error {
	if protocol.wlogger != nil {
		protocol.wlogger.Write(b)
	}

	return write(bw, b)
}

func (protocol *Protocol) writeByte(bw *bufio.Writer, c byte) error {
	if protocol.wlogger != nil {
		protocol.wlogger.WriteRune(rune(c))
	}

	return writeByte(bw, c)
}

func (protocol *Protocol) writeString(bw *bufio.Writer, s string) error {
	if protocol.wlogger != nil {
		protocol.wlogger.WriteString(s)
	}

	return writeString(bw, s)
}

func write(bw *bufio.Writer, b []byte) error {
	_, err := bw.Write(b)
	return err
}

func writeByte(bw *bufio.Writer, c byte) error {
	return bw.WriteByte(c)
}

func writeString(bw *bufio.Writer, s string) error {
	_, err := bw.WriteString(s)
	return err
}

func (protocol *Protocol) nextItem(br *bufio.Reader) (*Item, error) {
	lineType, line, err := protocol.readLine(br)
	if err != nil {
		return nil, err
	}

	switch lineType {
	case simpleStringType:
		return protocol.simpleStringItem(line[1:])
	case errorType:
		return protocol.errorItem(line[1:])
	case integerType:
		return protocol.integerItem(line[1:])
	case bulkStringType:
		return protocol.bulkStringItem(line[1:], br)
	case arrayType:
		return protocol.arrayItem(line[1:], br)
	default:
		return nil, errInlineCommandNotImplemented
	}
}

func (protocol *Protocol) simpleStringItem(line []byte) (*Item, error) {
	return &Item{
		Type:  simpleStringType,
		Value: string(line),
	}, nil
}

func (protocol *Protocol) errorItem(line []byte) (*Item, error) {
	return &Item{
		Type:  errorType,
		Value: errors.New(string(line)),
	}, nil
}

func (protocol *Protocol) integerItem(line []byte) (*Item, error) {
	val, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, err
	}

	return &Item{
		Type:  integerType,
		Value: val,
	}, nil
}

func (protocol *Protocol) bulkStringItem(length []byte, br *bufio.Reader) (*Item, error) {
	ln, err := strconv.Atoi(string(length))
	if err != nil {
		return nil, err
	}

	if ln > bulkStringMaxLength {
		return nil, errBulkStringIsTooLarge
	}

	if ln < 0 {
		return &Item{
			Type:  bulkStringType,
			Value: nil,
		}, nil
	}

	_, line, err := protocol.readLine(br)
	if err != nil {
		return nil, err
	}

	l := len(line)
	if l != ln {
		return nil, errInvalidFormat
	}

	return &Item{
		Type:  bulkStringType,
		Value: line,
	}, nil
}

func (protocol *Protocol) arrayItem(line []byte, br *bufio.Reader) (*Item, error) {
	l, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, err
	}

	if l < 0 {
		return &Item{
			Type:  arrayType,
			Value: nil,
		}, nil
	}

	items := make([]*Item, l)
	for i := range items {
		item, err := protocol.nextItem(br)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return &Item{
		Type:  arrayType,
		Value: items,
	}, nil
}

func (protocol *Protocol) End(bw *bufio.Writer) error {
	if protocol.wlogger != nil && protocol.wlogger.Len() > 0 {
		log.Printf("W: %s", string(protocol.wlogger.Bytes()))
		protocol.wlogger.Reset()
	}
	return bw.Flush()
}

// TODO: fix incorrect end of line when multiline value separated by '\n'
// TODO: In most cases it is ok! :)
// TODO: https://redis.io/topics/protocol#simple-string-reply
// TODO: Simple Strings are encoded in the following way: a plus character, followed by a string
// TODO: that cannot contain a CR or LF character (no newlines are allowed), terminated by CRLF (that is "\r\n").
func (protocol *Protocol) readLine(br *bufio.Reader) (byte, []byte, error) {
	line, err := br.ReadBytes(lf)
	if err == nil {
		if protocol.rlogger != nil {
			protocol.rlogger.Write(line)
			log.Printf("R: %s", string(protocol.rlogger.Bytes()))
			protocol.rlogger.Reset()
		}

		if bytes.HasSuffix(line, crlfSlice) {
			l := len(line)
			return line[0], line[:l-2], nil
		}
		err = io.ErrUnexpectedEOF
	}

	return 0, nil, err
}
