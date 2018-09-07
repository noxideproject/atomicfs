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

func Test_FileWriter_Write(t *testing.T) {
	mockFS := &fstest.FileSystem{}
	defer mockFS.AssertExpectations(t)

	mockFile := &fstest.File{}
	defer mockFile.AssertExpectations(t)

	mockSys := &systest.Syscall{}
	defer mockSys.AssertExpectations(t)

	mockFS.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil).Once()
	mockFS.On("Open", ".").Return(mockFile, nil).Once()

	mockFile.On("Sync").Return(nil).Once()
	mockFile.On("Close").Return(nil).Once()

	mockSys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mockSys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mockFS,
		Sys:          mockSys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.NoError(t, err)
}

func Test_FileWriter_Write_bad_Rename(t *testing.T) {
	mockFS := &fstest.FileSystem{}
	defer mockFS.AssertExpectations(t)

	mockFile := &fstest.File{}
	defer mockFile.AssertExpectations(t)

	mockSys := &systest.Syscall{}
	defer mockSys.AssertExpectations(t)

	mockFS.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(
		errors.New("rename failed"),
	).Once()

	mockFS.On("Remove", mock.AnythingOfType("string")).Return(nil).Once() // ignored

	mockSys.On("Stat", ".", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()
	mockSys.On("Stat", "/tmp", mock.AnythingOfType("*syscall.Stat_t")).Return(nil).Once()

	writer := NewFileWriter(FileWriterOptions{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mockFS,
		Sys:          mockSys,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "rename failed")
}
