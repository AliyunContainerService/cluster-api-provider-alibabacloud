// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
//go:build go1.11
=======
>>>>>>> 79bfea2d (update vendor)
=======
//go:build go1.11
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
//go:build go1.11
>>>>>>> 03397665 (update api)
// +build go1.11

package gcimporter

import "go/types"

func newInterface(methods []*types.Func, embeddeds []types.Type) *types.Interface {
	return types.NewInterfaceType(methods, embeddeds)
}
