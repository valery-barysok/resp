package resp

type Message struct {
	Type  byte
	Value interface{}
}

func (msg *Message) IsNil() bool {
	return msg.Value == nil
}

func (msg *Message) IsString() bool {
	return msg.Type == simpleStringType
}

func (msg *Message) IsError() bool {
	return msg.Type == errorType
}

func (msg *Message) IsInt() bool {
	return msg.Type == integerType
}

func (msg *Message) IsBulkString() bool {
	return msg.Type == bulkStringType
}

func (msg *Message) IsArray() bool {
	return msg.Type == arrayType
}

func (msg *Message) String() string {
	return msg.Value.(string)
}

func (msg *Message) Err() error {
	return msg.Value.(error)
}

func (msg *Message) Int() int {
	return msg.Value.(int)
}

func (msg *Message) BulkString() []byte {
	return msg.Value.([]byte)
}

func (msg *Message) Array() []*Message {
	return msg.Value.([]*Message)
}

func (msg *Message) Raw() interface{} {
	return msg.Value
}
