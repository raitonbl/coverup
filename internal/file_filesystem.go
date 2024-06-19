package internal

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Filesystem struct {
	WorkDirectory string
}

func (f Filesystem) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(f.WorkDirectory, name))
}

func (f Filesystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(f.WorkDirectory, name))
}
