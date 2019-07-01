package atomicfs

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"

	"go.gophers.dev/pkgs/atomicfs/fs"
	"go.gophers.dev/pkgs/atomicfs/sys"
)

func setup(t *testing.T) string {
	dir, err := ioutil.TempDir("", "atomicfs-test-")
	require.NoError(t, err)
	return dir
}

func cleanup(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	require.NoError(t, err)
}

func Test_FileWriter_Write_osFS(t *testing.T) {
	tmpDir := setup(t)
	defer cleanup(t, tmpDir)

	writer := NewFileWriter(Options{
		TmpDirectory: tmpDir,
		TmpExtension: ".tempfile",
		Mode:         0600,
		FS:           fs.New(),
		Sys:          sys.New(),
	})

	input := strings.NewReader("foobar")
	filePath := filepath.Join(tmpDir, "out.txt")
	err := writer.Write(input, filePath)
	require.NoError(t, err)
}

type mocks struct {
	fs   *fs.FileSystemMock
	file *fs.FileMock
	sys  *sys.SyscallMock
}

func (m *mocks) assertions(t *testing.T) {
	t.Log("test over, asserting mock behavior")
	m.fs.MinimockFinish()
	m.file.MinimockFinish()
	m.sys.MinimockFinish()
}

func newMocks(t *testing.T) *mocks {
	return &mocks{
		fs:   fs.NewFileSystemMock(t),
		file: fs.NewFileMock(t),
		sys:  sys.NewSyscallMock(t),
	}
}

func Test_FileWriter_Write(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.fs.RenameMock.Set(func(_, newName string) error {
		if newName != "out.txt" {
			t.Fatal("expected out.txt as new name")
		}
		return nil
	})

	mocks.fs.OpenMock.Expect(".").Return(mocks.file, nil)
	mocks.file.SyncMock.Return(nil)
	mocks.file.CloseMock.Return(nil)
	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." && name != "/tmp" {
			t.Fatal("expected name to be '.' or '/tmp'")
		}
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.NoError(t, err)
}

func Test_FileWriter_Write_bad_Rename(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.fs.RenameMock.Set(func(_, newName string) error {
		if newName != "out.txt" {
			t.Fatal("expected new file to be out.txt")
		}
		return errors.New("rename failed")
	})

	mocks.fs.RemoveMock.Set(func(_ string) error {
		// the file is some random .tmp file
		return nil
	})

	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." && name != "/tmp" {
			t.Fatal("expected name to be '.' or '/tmp'")
		}
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "rename failed")
}

func Test_FileWriter_Write_bad_Open(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.fs.RenameMock.Set(func(_, newName string) error {
		if newName != "out.txt" {
			t.Fatal("expected new name to be out.txt")
		}
		return nil
	})
	mocks.fs.OpenMock.Expect(".").Return(nil, errors.New("open failed"))
	mocks.fs.RemoveMock.Set(func(_ string) error {
		return nil
	})
	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." && name != "/tmp" {
			t.Fatal("expected name to be '.' or '/tmp'")
		}
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "open failed")
}

func Test_FileWriter_Write_bad_stat_destination(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." {
			t.Fatal("expected name to be '.'")
		}
		return errors.New("stat failed")
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "stat failed")
}

func Test_FileWriter_Write_bad_stat_tmp_dir(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name == "." {
			return nil
		}

		if name == "/tmp" {
			return errors.New("stat failed")
		}

		t.Fatal("name expected to be '.' or '/tmp'")
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "stat failed")
}

func Test_FileWriter_Write_bad_Sync(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.fs.RenameMock.Set(func(_, old string) error {
		if old != "out.txt" {
			t.Fatal("expected new name to be 'out.txt'")
		}
		return nil
	})
	mocks.fs.OpenMock.Expect(".").Return(mocks.file, nil)
	mocks.fs.RemoveMock.Set(func(_ string) error {
		return nil
	})
	mocks.file.SyncMock.Return(errors.New("sync failed"))
	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." && name != "/tmp" {
			t.Fatal("expected name to be '.' or '/tmp'")
		}
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "sync failed")
}

func Test_FileWriter_Write_bad_Close(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.fs.RenameMock.Set(func(_, newName string) error {
		if newName != "out.txt" {
			t.Fatal("expected new name to be 'out.txt'")
		}
		return nil
	})
	mocks.fs.OpenMock.Expect(".").Return(mocks.file, nil)
	mocks.fs.RemoveMock.Set(func(_ string) error {
		return nil
	})
	mocks.file.SyncMock.Return(nil)
	mocks.file.CloseMock.Return(errors.New("close failed"))
	mocks.sys.StatMock.Set(func(name string, _ *syscall.Stat_t) error {
		if name != "." && name != "/tmp" {
			t.Fatal("expected name to be '.' or '/tmp'")
		}
		return nil
	})

	writer := NewFileWriter(Options{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mocks.fs,
		Sys:          mocks.sys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "close failed")
}
