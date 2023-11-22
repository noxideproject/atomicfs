package atomicfs

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
	"noxide.lol/go/atomicfs/fs"
	"noxide.lol/go/atomicfs/sys"
)

func TestFileWriter_WriteFile(t *testing.T) {
	tmpDir := t.TempDir()

	writer := New(Options{
		TmpDirectory: tmpDir,
		TmpExtension: ".tempfile",
		Mode:         0600,
		FS:           fs.New(),
		Sys:          sys.New(),
	})

	input := strings.NewReader("foobar")
	filePath := filepath.Join(tmpDir, "out.txt")
	err := writer.WriteFile(input, filePath)
	must.NoError(t, err)
}
