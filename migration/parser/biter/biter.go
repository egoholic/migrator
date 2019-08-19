package biter

import "errors"

type Iterator struct {
	bytes []byte
	len   int
	cur   int
}

func New(bytes []byte) *Iterator {
	return &Iterator{
		bytes: bytes,
		len:   len(bytes),
		cur:   0,
	}
}
func (i Iterator) HasNext() bool {
	return i.len > i.cur+1
}
func (i Iterator) Next() (b byte, err error) {
	if i.HasNext() {
		b = i.bytes[i.cur]
		i.cur++
		return b, nil
	}
	return b, errors.New("")
}
