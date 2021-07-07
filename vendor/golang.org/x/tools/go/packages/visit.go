<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
package packages

import (
	"fmt"
	"os"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"sort"
=======
>>>>>>> 79bfea2d (update vendor)
=======
	"sort"
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	"sort"
>>>>>>> 03397665 (update api)
)

// Visit visits all the packages in the import graph whose roots are
// pkgs, calling the optional pre function the first time each package
// is encountered (preorder), and the optional post function after a
// package's dependencies have been visited (postorder).
// The boolean result of pre(pkg) determines whether
// the imports of package pkg are visited.
func Visit(pkgs []*Package, pre func(*Package) bool, post func(*Package)) {
	seen := make(map[*Package]bool)
	var visit func(*Package)
	visit = func(pkg *Package) {
		if !seen[pkg] {
			seen[pkg] = true

			if pre == nil || pre(pkg) {
				paths := make([]string, 0, len(pkg.Imports))
				for path := range pkg.Imports {
					paths = append(paths, path)
				}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
				sort.Strings(paths) // Imports is a map, this makes visit stable
=======
>>>>>>> 79bfea2d (update vendor)
=======
				sort.Strings(paths) // Imports is a map, this makes visit stable
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
				sort.Strings(paths) // Imports is a map, this makes visit stable
>>>>>>> 03397665 (update api)
				for _, path := range paths {
					visit(pkg.Imports[path])
				}
			}

			if post != nil {
				post(pkg)
			}
		}
	}
	for _, pkg := range pkgs {
		visit(pkg)
	}
}

// PrintErrors prints to os.Stderr the accumulated errors of all
// packages in the import graph rooted at pkgs, dependencies first.
// PrintErrors returns the number of errors printed.
func PrintErrors(pkgs []*Package) int {
	var n int
	Visit(pkgs, nil, func(pkg *Package) {
		for _, err := range pkg.Errors {
			fmt.Fprintln(os.Stderr, err)
			n++
		}
	})
	return n
}
