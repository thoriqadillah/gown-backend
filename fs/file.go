package fs

import (
	"os"
	"path/filepath"

	"github.com/thoriqadillah/gown/config"
)

type File struct {
	Data [][]byte
	*config.Config
}

func New(size int, config *config.Config) *File {
	return &File{
		Data:   make([][]byte, size),
		Config: config,
	}
}

func (f *File) Combine(data []byte, index int) {
	f.Data[index] = data
}

func (f *File) Save(name string) error {
	if _, err := os.Stat(f.SaveLocation); err != nil {
		if err := os.MkdirAll(f.SaveLocation, os.ModePerm); err != nil {
			return err
		}
	}

	combined := []byte{}
	for i := 0; i < len(f.Data); i++ {
		combined = append(combined, f.Data[i]...)
	}

	return os.WriteFile(filepath.Join(f.SaveLocation, name), combined, os.ModePerm)
}
