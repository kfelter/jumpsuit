package jumpsuit

import (
	"os"
	"sync"
)

type FileStore struct {
	sync.Mutex
	BasePath string
}

func loadf(path string) (*Memory, error) {
	inf, err := os.Open(path)
	if os.IsNotExist(err) {
		inf, err = os.Create(path)
	}
	if err != nil {

		return nil, err
	}
	ms := NewMemoryStore()
	if err = ms.Load(inf); err != nil {
		return nil, err
	}
	inf.Close()
	return ms, nil
}

func savef(ms *Memory, path string) error {
	outf, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = ms.Dump(outf); err != nil {
		return err
	}

	return outf.Close()
}

func (f *FileStore) Get(table string, objID int64) (any, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.path(table))
	if err != nil {
		return nil, err
	}
	obj, err := ms.Get(objID)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (f *FileStore) Del(table string, objID int64) error {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.path(table))
	if err != nil {
		return err
	}

	if err = ms.Del(objID); err != nil {
		return err
	}

	return savef(ms, f.path(table))
}

func (f *FileStore) Put(table string, objID int64, obj any) error {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.path(table))
	if err != nil {
		return err
	}

	if err = ms.Put(objID, obj); err != nil {
		return err
	}

	return savef(ms, f.path(table))
}

func (f *FileStore) Lst(table string) ([]any, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.path(table))
	if err != nil {
		return nil, err
	}
	return ms.Lst()
}

func (f *FileStore) Inc(table string) (int64, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.path(table))
	if err != nil {
		return 0, err
	}
	return ms.Inc(), nil
}

func (f *FileStore) path(table string) string {
	return f.BasePath + "/" + table + ".json"
}
