// persist/filesystem.go

package persist

import (
	"io"
	"os"
)

type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Stat(path string) (os.FileInfo, error)
	Open(path string) (io.ReadCloser, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
	Rename(oldpath, newpath string) error
}

type OSFileSystem struct{}

func NewOSFileSystem() FileSystem {
	return OSFileSystem{}
}

func (OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (OSFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (OSFileSystem) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OSFileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (OSFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (OSFileSystem) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}
