// Package atomicfs provides tools for doing all-or-nothing atomic filesystem operations in Go.
package atomicfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"go.gophers.dev/pkgs/atomicfs/fs"
	"go.gophers.dev/pkgs/atomicfs/sys"
)

// A FileWriter is used to read content from a source and
// write it to a destination file, or die trying without producing
// corrupt data where the destination file should have been.
//
// The only design goal of a FileWriter is to be oriented around
// correctness - performance is not a consideration.
type FileWriter interface {
	// Write the contents of the io.Reader into a file at the
	// destination filepath given.
	Write(io.Reader, string) error
}

// Options are used to configure the behavior of a
// FileWriter when it is used to write a file.
type Options struct {
	// TmpDirectory is where tmp files are generated during
	// the process of writing a file in preparation for executing
	// an atomic rename. Because of this, TmpDirectory *MUST* be on
	// the same filesystem device as the destination file being
	// written, otherwise, errors will be returned during the Write.
	TmpDirectory string

	// TmpExtension is used to demarcate the tmp files generated
	// during the Write operation from other normal files. The
	// extension is appended to the end of the filename of the
	// destination file being written.
	TmpExtension string

	// Mode is the file mode of the destination file to be written.
	// If no Mode is provided, write only permissions for the user
	// are used (i.e. 0200).
	Mode os.FileMode

	// FS is the underlying filesystem implementation to be used
	// for writing files to disk. Typically this should be left
	// nil, as then fs.FileSystem is used, which in turn is
	// implemented using the file operations defined in the os
	// package.
	FS fs.FileSystem

	// Sys is the underlying syscall implementation to be used,
	// in this case only for stat-ing files. Typically this should
	// be left nil, as then sys.Syscall is used, which in turn is
	// implemented using the syscall implementation provided by
	// the go runtime.
	Sys sys.Syscall
}

// NewFileWriter creates a new FileWriter backed by the configuration
// settings in the provided Options. Creating a FileWriter always
// succeeds, replacing empty options with sane defaults.
func NewFileWriter(options Options) FileWriter {
	tmpExt := strings.TrimPrefix(options.TmpExtension, ".")
	if tmpExt == "" {
		tmpExt = "tmp"
	}

	tmpDir := options.TmpDirectory
	if tmpDir == "" {
		tmpDir = os.TempDir()
	}

	mode := options.Mode
	if mode == 0 {
		mode = os.FileMode(0200)
	}

	fileSystem := options.FS
	if fileSystem == nil {
		fileSystem = fs.New()
	}

	systemCalls := options.Sys
	if systemCalls == nil {
		systemCalls = sys.New()
	}

	return &fsFileWriter{
		tmpDir:   tmpDir,
		tmpExt:   tmpExt,
		fileMode: mode,
		fs:       fileSystem,
		sys:      systemCalls,
	}
}

type fsFileWriter struct {
	tmpDir   string
	tmpExt   string
	fileMode os.FileMode
	fs       fs.FileSystem
	sys      sys.Syscall
}

func (w *fsFileWriter) Write(source io.Reader, filePath string) error {
	fileDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	// 1) check the target directory is on the same filesystem device
	// as the configured tmp directory - otherwise atomic operations
	// are not possible
	if err := w.checkDevice(fileDir); err != nil {
		return err
	}

	// 2) write the content to a tmp file, also triggering
	// a complete flush (fsync)
	tmpPath, err := w.writeTmp(source, fileName)
	if err != nil {
		_ = w.fs.Remove(tmpPath)
		return err
	}

	// 3) since we know the tmp file and destination file exist on the
	// same device, a filesystem rename will be an atomic operation
	if err := w.rename(tmpPath, filePath); err != nil {
		_ = w.fs.Remove(tmpPath)
		return err
	}

	return nil
}

func (w *fsFileWriter) checkDevice(fileDir string) error {
	var stat syscall.Stat_t
	if err := w.sys.Stat(fileDir, &stat); err != nil {
		return errors.Wrapf(err, "atomicfs: unable to stat destination directory %s", fileDir)
	}
	fileDirDeviceID := stat.Dev

	if err := w.sys.Stat(w.tmpDir, &stat); err != nil {
		return errors.Wrapf(err, "atomicfs: unable to stat tmp directory %s", w.tmpDir)
	}
	tmpDirDeviceID := stat.Dev

	if fileDirDeviceID != tmpDirDeviceID {
		return errors.Errorf("atomicfs: tmp & destination directories not on same device")
	}

	return nil
}

func (w *fsFileWriter) rename(old, new string) error {
	if err := w.fs.Rename(old, new); err != nil {
		return errors.Wrapf(err, "atomicfs: unable to rename tmp file %s to %s", old, new)
	}

	return w.syncDir(new)
}

func (w *fsFileWriter) syncDir(name string) error {
	directory := filepath.Dir(name)
	f, err := w.fs.Open(directory)
	if err != nil {
		return errors.Wrapf(err, "atomicfs: unable to open directory %s for syncing", directory)
	}

	if err := f.Sync(); err != nil {
		return errors.Wrapf(err, "atomicfs: unable to sync directory %s", directory)
	}

	if err := f.Close(); err != nil {
		return errors.Wrapf(err, "atomicfs: unable to close directory %s after syncing", directory)
	}

	return nil
}

func (w *fsFileWriter) writeTmp(source io.Reader, name string) (string, error) {
	tmpName := fmt.Sprintf("%s.%s", name, w.tmpExt)

	f, err := ioutil.TempFile(w.tmpDir, tmpName)
	if err != nil {
		return "", errors.Wrapf(err, "atomicfs: unable to create tmp file in %s", w.tmpDir)
	}
	tmpPath := f.Name()

	if err := f.Chmod(w.fileMode); err != nil {
		_ = f.Close()
		return tmpPath, errors.Wrapf(err, "atomicfs: unable to chmod tmp file in %s", w.tmpDir)
	}

	if _, err := io.Copy(f, source); err != nil {
		_ = f.Close()
		return tmpPath, errors.Wrapf(err, "atomicfs: unable to copy source into tmp file in %s", w.tmpDir)
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		return tmpPath, errors.Wrapf(err, "atomicfs: unable to fsync tmp file in %s", w.tmpDir)
	}

	if err := f.Close(); err != nil {
		return tmpPath, errors.Wrapf(err, "atomicfs: unable to close tmp file in %s", w.tmpDir)
	}

	return tmpPath, nil
}
