package sys

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Syscall_Stat_file(t *testing.T) {
	sc := New()
	var stat syscall.Stat_t
	err := sc.Stat("/etc/hosts", &stat)
	require.NoError(t, err)

	t.Log("file dev:", stat.Dev)
}

func Test_Syscall_Stat_dir(t *testing.T) {
	sc := New()
	var stat syscall.Stat_t
	err := sc.Stat("/etc", &stat)
	require.NoError(t, err)

	t.Log("etc dev:", stat.Dev)
}
