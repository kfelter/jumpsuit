package jumpsuit

import "fmt"

var (
	Nil = fmt.Errorf("nil")
)

type Storage interface {
	Get(table string, objID int64) (any, error)
	Del(table string, objID int64) error
	Put(table string, objID int64, obj any) error
	Lst(table string) ([]any, error)
	Inc(table string) (int64, error)
}
