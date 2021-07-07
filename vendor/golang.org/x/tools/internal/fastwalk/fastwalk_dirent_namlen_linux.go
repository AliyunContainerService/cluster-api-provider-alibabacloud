// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
//go:build linux && !appengine
// +build linux,!appengine
=======
// +build linux
// +build !appengine
>>>>>>> 79bfea2d (update vendor)
=======
//go:build linux && !appengine
// +build linux,!appengine
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
//go:build linux && !appengine
// +build linux,!appengine
>>>>>>> 03397665 (update api)

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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	const nameBufLen = uint16(len(nameBuf))
	limit := dirent.Reclen - fixedHdr
	if limit > nameBufLen {
		limit = nameBufLen
	}
	nameLen := bytes.IndexByte(nameBuf[:limit], 0)
<<<<<<< HEAD
<<<<<<< HEAD
=======
	nameLen := bytes.IndexByte(nameBuf[:dirent.Reclen-fixedHdr], 0)
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	if nameLen < 0 {
		panic("failed to find terminating 0 byte in dirent")
	}
	return uint64(nameLen)
}
