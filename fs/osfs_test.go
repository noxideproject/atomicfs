package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_osFS_Create(t *testing.T) {
	fs := New()

	f, err := fs.Create("foo.txt")
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	err = os.Remove(f.Name())
	require.NoError(t, err)
}

func Test_osFS_Mkdir(t *testing.T) {
	fs := New()

	err := fs.Mkdir("foo", 0700)
	require.NoError(t, err)

	err = os.RemoveAll("foo")
	require.NoError(t, err)
}

func Test_osFS_MkdirAll(t *testing.T) {
	fs := New()

	err := fs.MkdirAll("foo/bar/baz", 0700)
	require.NoError(t, err)

	err = os.RemoveAll("foo")
	require.NoError(t, err)
}

func Test_osFS_Open(t *testing.T) {
	fs := New()

	f, err := fs.Open("/etc/hosts")
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)
}

func Test_osFS_OpenFile(t *testing.T) {
	fs := New()

	f, err := fs.OpenFile("/etc/hosts", os.O_RDONLY, 0600)
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)
}

func Test_osFS_Remove(t *testing.T) {
	fs := New()

	f, err := fs.Create("foo.txt")
	require.NoError(t, err)

	err = fs.Remove(f.Name())
	require.NoError(t, err)
}

func Test_osFS_RemoveAll(t *testing.T) {
	fs := New()

	err := fs.MkdirAll("foo/bar/baz/", 0700)
	require.NoError(t, err)

	err = fs.RemoveAll("foo")
	require.NoError(t, err)
}

func Test_osFS_Rename(t *testing.T) {
	fs := New()

	f, err := fs.Create("foo.txt")
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	err = fs.Rename("foo.txt", "bar.txt")
	require.NoError(t, err)

	err = os.Remove("bar.txt")
	require.NoError(t, err)
}

func Test_osFS_Stat(t *testing.T) {
	fs := New()

	info, err := fs.Stat("/etc/hosts")
	require.NoError(t, err)

	name := info.Name()
	require.Equal(t, "hosts", name)
}
