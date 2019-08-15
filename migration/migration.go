package migration

import (
	"io/ioutil"

	"github.com/egoholic/migrator/migration/parser"
)

type Migration parser.Result

func New(up, down string) *Migration {
	return &Migration{
		Up:   up,
		Down: down,
	}
}

func NewFromFile(absPath string) *Migration {
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	l1 := len()

}
