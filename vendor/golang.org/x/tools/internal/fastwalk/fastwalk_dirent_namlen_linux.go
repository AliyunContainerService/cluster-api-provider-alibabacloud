// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
//go:build linux && !appengine
// +build linux,!appengine
=======
// +build linux
// +build !appengine
>>>>>>> 79bfea2d (update vendor)

package fastwalk

import (
	"bytes"
	"syscall"
	"unsafe"
)

func direntNamlen(dirent *syscall.Dirent) uint64 {
	const fixedHdr = uint16(unsafe.Offsetof(syscall.Dirent{}.Name))
	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
<<<<<<< HEAD
	const nameBufLen = uint16(len(nameBuf))
	limit := dirent.Reclen - fixedHdr
	if limit > nameBufLen {
		limit = nameBufLen
	}
	nameLen := bytes.IndexByte(nameBuf[:limit], 0)
=======
	nameLen := bytes.IndexByte(nameBuf[:dirent.Reclen-fixedHdr], 0)
>>>>>>> 79bfea2d (update vendor)
	if nameLen < 0 {
		panic("failed to find terminating 0 byte in dirent")
	}
	return uint64(nameLen)
}
