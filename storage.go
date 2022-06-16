package jumpsuit

import "fmt"

var (
	Nil = fmt.Errorf("nil")
)

type Storage interface {
	Get(objID int64) (any, error)
	Del(objID int64) error
	Put(objID int64, obj any) error
	Lst() ([]any, error)
	Inc() (int64, error)
}
