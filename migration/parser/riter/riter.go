package riter

import (
	"errors"
)

type Iterator struct {
	runes []rune
	len   int
	cur   int
}

func New(runes []rune) *Iterator {
	return &Iterator{
		runes: runes,
		len:   len(runes),
		cur:   0,
	}
}
func (i *Iterator) HasNext() bool {
	return i.cur < i.len
}
func (i *Iterator) Next() (r rune, err error) {
	if i.HasNext() {
		r = i.runes[i.cur]
		i.cur = i.cur + 1
		return r, nil
	}
	return r, errors.New("there is no more runes")
}
