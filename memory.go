package jumpsuit

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"
)

type Memory struct {
	sync.Mutex
	Data map[int64]any
}

func (m *Memory) Get(objID int64) (any, error) {
	m.Lock()
	data, ok := m.Data[objID]
	m.Unlock()
	if !ok {
		return nil, Nil
	}
	return data, nil
}

func (m *Memory) Put(objID int64, obj any) error {
	m.Lock()
	m.Data[objID] = obj
	m.Unlock()
	return nil
}

func (m *Memory) Del(objID int64) error {
	m.Lock()
	delete(m.Data, objID)
	m.Unlock()
	return nil
}

func (m *Memory) Lst() ([]any, error) {
	lst := []any{}
	m.Lock()
	for _, v := range m.Data {
		lst = append(lst, v)
	}
	m.Unlock()
	return lst, nil
}

func (m *Memory) Dump(wr io.Writer) (int, error) {
	m.Lock()
	raw, err := json.Marshal(m.Data)
	m.Unlock()

	if err != nil {
		return 0, errors.Wrap(err, "json.Marshal m.data")
	}
	return wr.Write(raw)
}

func (m *Memory) Load(r io.Reader) error {
	newData := make(map[int64]any)
	raw, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "reading r")
	}
	json.Unmarshal(raw, &newData)
	m.Lock()
	m.Data = newData
	m.Unlock()
	return nil
}

func (m *Memory) Inc() int64 {
	return int64(len(m.Data))
}
