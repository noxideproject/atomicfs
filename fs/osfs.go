package fs

import "os"

type osFS struct{}

var _ FileSystem = (*osFS)(nil)
var _ File = (*os.File)(nil)

func New() FileSystem {
	return &osFS{}
}

func (*osFS) Create(name string) (File, error) {
	return os.Create(name)
}

func (*osFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (*osFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (*osFS) Open(name string) (File, error) {
	return os.Open(name)
}

func (*osFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}

func (*osFS) Remove(name string) error {
	return os.Remove(name)
}

func (*osFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (*osFS) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (*osFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
