// Copyright (c) NOXIDE.LOL
// SPDX-License-Identifier: BSD-3-Clause

package sys

import "syscall"

type Syscall interface {
	Stat(string, *syscall.Stat_t) error
}

type linuxSyscall struct{}

var _ Syscall = (*linuxSyscall)(nil)

func New() Syscall {
	return &linuxSyscall{}
}

func (s *linuxSyscall) Stat(path string, stat *syscall.Stat_t) error {
	return syscall.Stat(path, stat)
}
