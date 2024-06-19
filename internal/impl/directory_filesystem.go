package impl

import (
	"io/fs"
	"os"
	"path/filepath"
)

type WorkDirectoryFilesystem struct {
	WorkDirectory string
}

func (f WorkDirectoryFilesystem) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(f.WorkDirectory, name))
}

func (f WorkDirectoryFilesystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(f.WorkDirectory, name))
}
