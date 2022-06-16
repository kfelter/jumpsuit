package jumpsuit

import (
	"os"
	"sync"
)

type FileStore struct {
	sync.Mutex
	Path string
}

func loadf(path string) (*Memory, error) {
	ms := NewMemoryStore()
	inf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
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

func (f *FileStore) Get(objID int64) (any, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.Path)
	if err != nil {
		return nil, err
	}
	obj, err := ms.Get(objID)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (f *FileStore) Del(objID int64) error {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.Path)
	if err != nil {
		return err
	}

	if err = ms.Del(objID); err != nil {
		return err
	}

	return savef(ms, f.Path)
}

func (f *FileStore) Put(objID int64, obj any) error {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.Path)
	if err != nil {
		return err
	}

	if err = ms.Put(objID, obj); err != nil {
		return err
	}

	return savef(ms, f.Path)
}

func (f *FileStore) Lst() ([]any, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.Path)
	if err != nil {
		return nil, err
	}
	return ms.Lst()
}

func (f *FileStore) Inc() (int64, error) {
	f.Lock()
	defer f.Unlock()
	ms, err := loadf(f.Path)
	if err != nil {
		return 0, err
	}
	return ms.Inc(), nil
}
