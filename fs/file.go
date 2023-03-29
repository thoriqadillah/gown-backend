package fs

import (
	"os"
	"path/filepath"

	"github.com/thoriqadillah/gown/setting"
)

type File struct {
	Data [][]byte
	*setting.Setting
}

func New(size int, setting *setting.Setting) *File {
	return &File{
		Data:    make([][]byte, size),
		Setting: setting,
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
