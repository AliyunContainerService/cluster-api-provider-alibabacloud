// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris) && gc && !ppc64le && !ppc64
// +build darwin dragonfly freebsd linux netbsd openbsd solaris
<<<<<<< HEAD
<<<<<<< HEAD
// +build gc
// +build !ppc64le
// +build !ppc64
=======
// +build !gccgo
>>>>>>> 79bfea2d (update vendor)
=======
// +build gc
// +build !ppc64le
// +build !ppc64
>>>>>>> e879a141 (alibabacloud machine-api provider)

package unix

import "syscall"

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
