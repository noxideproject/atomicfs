package atomicfs

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/atomicfs/fs"
	"github.com/shoenig/atomicfs/fs/fstest"
	"github.com/shoenig/atomicfs/sys"
	"github.com/shoenig/atomicfs/sys/systest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

	writer := NewFileWriter(FileWriterOptions{
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
	fs   *fstest.FileSystem
	file *fstest.File
	sys  *systest.Syscall
}

func (m *mocks) assertions(t *testing.T) {
	m.fs.AssertExpectations(t)
	m.file.AssertExpectations(t)
	m.sys.AssertExpectations(t)
}

func newMocks() *mocks {
	return &mocks{
		fs:   &fstest.FileSystem{},
		file: &fstest.File{},
		sys:  &systest.Syscall{},
	}
}

func Test_FileWriter_Write(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.fs.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil).Once()
	mocks.fs.On("Open", ".").Return(mocks.file, nil).Once()

	mocks.file.On("Sync").Return(nil).Once()
	mocks.file.On("Close").Return(nil).Once()

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.fs.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(
		errors.New("rename failed"),
	).Once()

	mocks.fs.On("Remove", mock.AnythingOfType("string")).Return(nil).Once() // ignored

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.fs.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil).Once()
	mocks.fs.On("Open", ".").Return(nil, errors.New(
		"open failed",
	)).Once()

	mocks.fs.On("Remove", mock.AnythingOfType("string")).Return(nil).Once() // ignored

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(
		errors.New("stat failed"),
	).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(
		errors.New("stat failed"),
	).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.fs.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil).Once()
	mocks.fs.On("Open", ".").Return(mocks.file, nil).Once()

	mocks.fs.On("Remove", mock.AnythingOfType("string")).Return(nil).Once() // ignored

	mocks.file.On("Sync").Return(
		errors.New("sync failed"),
	).Once()
	// mocks.file.On("Close").Return(nil).Once()

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
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
	mocks := newMocks()
	defer mocks.assertions(t)

	mocks.fs.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil).Once()
	mocks.fs.On("Open", ".").Return(mocks.file, nil).Once()

	mocks.fs.On("Remove", mock.AnythingOfType("string")).Return(nil).Once() // ignored

	mocks.file.On("Sync").Return(nil).Once()
	mocks.file.On("Close").Return(
		errors.New("close failed"),
	).Once()

	mocks.sys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mocks.sys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
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
