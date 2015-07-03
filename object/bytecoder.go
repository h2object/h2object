package object

type ByteCoder interface{
	Decode(value []byte) (interface{}, error)
	Encode(value interface{}) ([]byte, error)
}

