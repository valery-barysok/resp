package resp

type Message struct {
	Type  byte
	Value interface{}
}

func (item *Message) IsNil() bool {
	return item.Value == nil
}

func (item *Message) IsString() bool {
	return item.Type == simpleStringType
}

func (item *Message) IsError() bool {
	return item.Type == errorType
}

func (item *Message) IsInt() bool {
	return item.Type == integerType
}

func (item *Message) IsBulkString() bool {
	return item.Type == bulkStringType
}

func (item *Message) IsArray() bool {
	return item.Type == arrayType
}

func (item *Message) String() string {
	return item.Value.(string)
}

func (item *Message) Err() error {
	return item.Value.(error)
}

func (item *Message) Int() int {
	return item.Value.(int)
}

func (item *Message) BulkString() []byte {
	return item.Value.([]byte)
}

func (item *Message) Array() []*Message {
	return item.Value.([]*Message)
}

func (item *Message) Raw() interface{} {
	return item.Value
}
