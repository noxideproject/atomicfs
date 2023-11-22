package sys

import (
	"syscall"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_Syscall_Stat_file(t *testing.T) {
	sc := New()
	var stat syscall.Stat_t
	err := sc.Stat("/etc/hosts", &stat)
	must.NoError(t, err)

	t.Log("file dev:", stat.Dev)
}

func Test_Syscall_Stat_dir(t *testing.T) {
	sc := New()
	var stat syscall.Stat_t
	err := sc.Stat("/etc", &stat)
	must.NoError(t, err)

	t.Log("etc dev:", stat.Dev)
}
