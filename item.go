package resp

type Item struct {
	Type  byte
	Value interface{}
}

func (item *Item) IsNil() bool {
	return item.Value == nil
}

func (item *Item) IsString() bool {
	return item.Type == simpleStringType
}

func (item *Item) IsError() bool {
	return item.Type == errorType
}

func (item *Item) IsInt() bool {
	return item.Type == integerType
}

func (item *Item) IsBulkString() bool {
	return item.Type == bulkStringType
}

func (item *Item) IsArray() bool {
	return item.Type == arrayType
}

func (item *Item) String() string {
	return item.Value.(string)
}

func (item *Item) Err() error {
	return item.Value.(error)
}

func (item *Item) Int() int {
	return item.Value.(int)
}

func (item *Item) BulkString() []byte {
	return item.Value.([]byte)
}

func (item *Item) Array() []*Item {
	return item.Value.([]*Item)
}

func (item *Item) Raw() interface{} {
	return item.Value
}
