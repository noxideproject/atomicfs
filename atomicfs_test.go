package atomicfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/atomicfs/fs"
	"github.com/shoenig/atomicfs/fs/fstest"
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

	mockFS.On("Rename", mock.AnythingOfType("string"), "out.txt").Return(nil)
	mockFS.On("Open", ".").Return(mockFile, nil)
	mockFile.On("Sync").Return(nil)
	mockFile.On("Close").Return(nil)

	writer := NewFileWriter(FileWriterOptions{
		TmpDirectory: "/tmp",
		TmpExtension: ".temp",
		Mode:         0600,
		FS:           mockFS,
	})

	input := strings.NewReader("foobar")
	err := writer.Write(input, "out.txt")
	require.NoError(t, err)
}
