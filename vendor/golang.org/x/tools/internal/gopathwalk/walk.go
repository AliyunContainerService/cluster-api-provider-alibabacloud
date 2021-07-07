// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gopathwalk is like filepath.Walk but specialized for finding Go
// packages, particularly in $GOPATH and $GOROOT.
package gopathwalk

import (
	"bufio"
	"bytes"
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"go/build"
	"golang.org/x/tools/internal/fastwalk"
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"time"

	"golang.org/x/tools/internal/fastwalk"
=======
>>>>>>> 79bfea2d (update vendor)
=======
	"time"

	"golang.org/x/tools/internal/fastwalk"
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	"time"

	"golang.org/x/tools/internal/fastwalk"
>>>>>>> 03397665 (update api)
)

// Options controls the behavior of a Walk call.
type Options struct {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	// If Logf is non-nil, debug logging is enabled through this function.
	Logf func(format string, args ...interface{})
	// Search module caches. Also disables legacy goimports ignore rules.
	ModulesEnabled bool
<<<<<<< HEAD
<<<<<<< HEAD
=======
	Debug          bool // Enable debug logging
	ModulesEnabled bool // Search module caches. Also disables legacy goimports ignore rules.
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
}

// RootType indicates the type of a Root.
type RootType int

const (
	RootUnknown RootType = iota
	RootGOROOT
	RootGOPATH
	RootCurrentModule
	RootModuleCache
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	RootOther
=======
>>>>>>> 79bfea2d (update vendor)
=======
	RootOther
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	RootOther
>>>>>>> 03397665 (update api)
)

// A Root is a starting point for a Walk.
type Root struct {
	Path string
	Type RootType
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// Walk walks Go source directories ($GOROOT, $GOPATH, etc) to find packages.
// For each package found, add will be called (concurrently) with the absolute
// paths of the containing source directory and the package directory.
// add will be called concurrently.
func Walk(roots []Root, add func(root Root, dir string), opts Options) {
	WalkSkip(roots, add, func(Root, string) bool { return false }, opts)
}

// WalkSkip walks Go source directories ($GOROOT, $GOPATH, etc) to find packages.
// For each package found, add will be called (concurrently) with the absolute
// paths of the containing source directory and the package directory.
// For each directory that will be scanned, skip will be called (concurrently)
// with the absolute paths of the containing source directory and the directory.
// If skip returns false on a directory it will be processed.
// add will be called concurrently.
// skip will be called concurrently.
func WalkSkip(roots []Root, add func(root Root, dir string), skip func(root Root, dir string) bool, opts Options) {
	for _, root := range roots {
		walkDir(root, add, skip, opts)
	}
}

// walkDir creates a walker and starts fastwalk with this walker.
func walkDir(root Root, add func(Root, string), skip func(root Root, dir string) bool, opts Options) {
	if _, err := os.Stat(root.Path); os.IsNotExist(err) {
		if opts.Logf != nil {
			opts.Logf("skipping nonexistent directory: %v", root.Path)
		}
		return
	}
	start := time.Now()
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanning %s", root.Path)
=======
// SrcDirsRoots returns the roots from build.Default.SrcDirs(). Not modules-compatible.
func SrcDirsRoots() []Root {
	var roots []Root
	roots = append(roots, Root{filepath.Join(build.Default.GOROOT, "src"), RootGOROOT})
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		roots = append(roots, Root{filepath.Join(p, "src"), RootGOPATH})
	}
	return roots
}

=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
// Walk walks Go source directories ($GOROOT, $GOPATH, etc) to find packages.
// For each package found, add will be called (concurrently) with the absolute
// paths of the containing source directory and the package directory.
// add will be called concurrently.
func Walk(roots []Root, add func(root Root, dir string), opts Options) {
	WalkSkip(roots, add, func(Root, string) bool { return false }, opts)
}

// WalkSkip walks Go source directories ($GOROOT, $GOPATH, etc) to find packages.
// For each package found, add will be called (concurrently) with the absolute
// paths of the containing source directory and the package directory.
// For each directory that will be scanned, skip will be called (concurrently)
// with the absolute paths of the containing source directory and the directory.
// If skip returns false on a directory it will be processed.
// add will be called concurrently.
// skip will be called concurrently.
func WalkSkip(roots []Root, add func(root Root, dir string), skip func(root Root, dir string) bool, opts Options) {
	for _, root := range roots {
		walkDir(root, add, skip, opts)
	}
}

// walkDir creates a walker and starts fastwalk with this walker.
func walkDir(root Root, add func(Root, string), skip func(root Root, dir string) bool, opts Options) {
	if _, err := os.Stat(root.Path); os.IsNotExist(err) {
		if opts.Logf != nil {
			opts.Logf("skipping nonexistent directory: %v", root.Path)
		}
		return
	}
<<<<<<< HEAD
<<<<<<< HEAD
	if opts.Debug {
		log.Printf("scanning %s", root.Path)
>>>>>>> 79bfea2d (update vendor)
=======
	start := time.Now()
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanning %s", root.Path)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	start := time.Now()
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanning %s", root.Path)
>>>>>>> 03397665 (update api)
	}
	w := &walker{
		root: root,
		add:  add,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		skip: skip,
=======
>>>>>>> 79bfea2d (update vendor)
=======
		skip: skip,
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
		skip: skip,
>>>>>>> 03397665 (update api)
		opts: opts,
	}
	w.init()
	if err := fastwalk.Walk(root.Path, w.walk); err != nil {
		log.Printf("gopathwalk: scanning directory %v: %v", root.Path, err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanned %s in %v", root.Path, time.Since(start))
=======
	if opts.Debug {
		log.Printf("scanned %s", root.Path)
>>>>>>> 79bfea2d (update vendor)
=======
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanned %s in %v", root.Path, time.Since(start))
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	if opts.Logf != nil {
		opts.Logf("gopathwalk: scanned %s in %v", root.Path, time.Since(start))
>>>>>>> 03397665 (update api)
	}
}

// walker is the callback for fastwalk.Walk.
type walker struct {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	root Root                    // The source directory to scan.
	add  func(Root, string)      // The callback that will be invoked for every possible Go package dir.
	skip func(Root, string) bool // The callback that will be invoked for every dir. dir is skipped if it returns true.
	opts Options                 // Options passed to Walk by the user.
<<<<<<< HEAD
<<<<<<< HEAD
=======
	root Root               // The source directory to scan.
	add  func(Root, string) // The callback that will be invoked for every possible Go package dir.
	opts Options            // Options passed to Walk by the user.
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)

	ignoredDirs []os.FileInfo // The ignored directories, loaded from .goimportsignore files.
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// init initializes the walker based on its Options
=======
// init initializes the walker based on its Options.
>>>>>>> 79bfea2d (update vendor)
=======
// init initializes the walker based on its Options
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// init initializes the walker based on its Options
>>>>>>> 03397665 (update api)
func (w *walker) init() {
	var ignoredPaths []string
	if w.root.Type == RootModuleCache {
		ignoredPaths = []string{"cache"}
	}
	if !w.opts.ModulesEnabled && w.root.Type == RootGOPATH {
		ignoredPaths = w.getIgnoredDirs(w.root.Path)
		ignoredPaths = append(ignoredPaths, "v", "mod")
	}

	for _, p := range ignoredPaths {
		full := filepath.Join(w.root.Path, p)
		if fi, err := os.Stat(full); err == nil {
			w.ignoredDirs = append(w.ignoredDirs, fi)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
			if w.opts.Logf != nil {
				w.opts.Logf("Directory added to ignore list: %s", full)
			}
		} else if w.opts.Logf != nil {
			w.opts.Logf("Error statting ignored directory: %v", err)
<<<<<<< HEAD
=======
			if w.opts.Debug {
				log.Printf("Directory added to ignore list: %s", full)
			}
		} else if w.opts.Debug {
			log.Printf("Error statting ignored directory: %v", err)
>>>>>>> 79bfea2d (update vendor)
=======
			if w.opts.Logf != nil {
				w.opts.Logf("Directory added to ignore list: %s", full)
			}
		} else if w.opts.Logf != nil {
			w.opts.Logf("Error statting ignored directory: %v", err)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
		}
	}
}

// getIgnoredDirs reads an optional config file at <path>/.goimportsignore
// of relative directories to ignore when scanning for go files.
// The provided path is one of the $GOPATH entries with "src" appended.
func (w *walker) getIgnoredDirs(path string) []string {
	file := filepath.Join(path, ".goimportsignore")
	slurp, err := ioutil.ReadFile(file)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
	if w.opts.Logf != nil {
		if err != nil {
			w.opts.Logf("%v", err)
		} else {
			w.opts.Logf("Read %s", file)
<<<<<<< HEAD
=======
	if w.opts.Debug {
=======
	if w.opts.Logf != nil {
>>>>>>> e879a141 (alibabacloud machine-api provider)
		if err != nil {
			w.opts.Logf("%v", err)
		} else {
<<<<<<< HEAD
			log.Printf("Read %s", file)
>>>>>>> 79bfea2d (update vendor)
=======
			w.opts.Logf("Read %s", file)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
		}
	}
	if err != nil {
		return nil
	}

	var ignoredDirs []string
	bs := bufio.NewScanner(bytes.NewReader(slurp))
	for bs.Scan() {
		line := strings.TrimSpace(bs.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignoredDirs = append(ignoredDirs, line)
	}
	return ignoredDirs
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// shouldSkipDir reports whether the file should be skipped or not.
func (w *walker) shouldSkipDir(fi os.FileInfo, dir string) bool {
=======
func (w *walker) shouldSkipDir(fi os.FileInfo) bool {
>>>>>>> 79bfea2d (update vendor)
=======
// shouldSkipDir reports whether the file should be skipped or not.
func (w *walker) shouldSkipDir(fi os.FileInfo, dir string) bool {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// shouldSkipDir reports whether the file should be skipped or not.
func (w *walker) shouldSkipDir(fi os.FileInfo, dir string) bool {
>>>>>>> 03397665 (update api)
	for _, ignoredDir := range w.ignoredDirs {
		if os.SameFile(fi, ignoredDir) {
			return true
		}
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	if w.skip != nil {
		// Check with the user specified callback.
		return w.skip(w.root, dir)
	}
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
	return false
}

// walk walks through the given path.
func (w *walker) walk(path string, typ os.FileMode) error {
	dir := filepath.Dir(path)
	if typ.IsRegular() {
		if dir == w.root.Path && (w.root.Type == RootGOROOT || w.root.Type == RootGOPATH) {
			// Doesn't make sense to have regular files
			// directly in your $GOPATH/src or $GOROOT/src.
			return fastwalk.ErrSkipFiles
<<<<<<< HEAD
=======
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return false
}

// walk walks through the given path.
func (w *walker) walk(path string, typ os.FileMode) error {
	dir := filepath.Dir(path)
	if typ.IsRegular() {
		if dir == w.root.Path && (w.root.Type == RootGOROOT || w.root.Type == RootGOPATH) {
			// Doesn't make sense to have regular files
			// directly in your $GOPATH/src or $GOROOT/src.
<<<<<<< HEAD
			return fastwalk.SkipFiles
>>>>>>> 79bfea2d (update vendor)
=======
			return fastwalk.ErrSkipFiles
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		w.add(w.root, dir)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		return fastwalk.ErrSkipFiles
=======
		return fastwalk.SkipFiles
>>>>>>> 79bfea2d (update vendor)
=======
		return fastwalk.ErrSkipFiles
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
		return fastwalk.ErrSkipFiles
>>>>>>> 03397665 (update api)
	}
	if typ == os.ModeDir {
		base := filepath.Base(path)
		if base == "" || base[0] == '.' || base[0] == '_' ||
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
			base == "testdata" ||
			(w.root.Type == RootGOROOT && w.opts.ModulesEnabled && base == "vendor") ||
			(!w.opts.ModulesEnabled && base == "node_modules") {
			return filepath.SkipDir
		}
		fi, err := os.Lstat(path)
		if err == nil && w.shouldSkipDir(fi, path) {
=======
			base == "testdata" || (!w.opts.ModulesEnabled && base == "node_modules") {
			return filepath.SkipDir
		}
		fi, err := os.Lstat(path)
		if err == nil && w.shouldSkipDir(fi) {
>>>>>>> 79bfea2d (update vendor)
=======
=======
>>>>>>> 03397665 (update api)
			base == "testdata" ||
			(w.root.Type == RootGOROOT && w.opts.ModulesEnabled && base == "vendor") ||
			(!w.opts.ModulesEnabled && base == "node_modules") {
			return filepath.SkipDir
		}
		fi, err := os.Lstat(path)
		if err == nil && w.shouldSkipDir(fi, path) {
<<<<<<< HEAD
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
			return filepath.SkipDir
		}
		return nil
	}
	if typ == os.ModeSymlink {
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".#") {
			// Emacs noise.
			return nil
		}
		fi, err := os.Lstat(path)
		if err != nil {
			// Just ignore it.
			return nil
		}
		if w.shouldTraverse(dir, fi) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
			return fastwalk.ErrTraverseLink
=======
			return fastwalk.TraverseLink
>>>>>>> 79bfea2d (update vendor)
=======
			return fastwalk.ErrTraverseLink
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
			return fastwalk.ErrTraverseLink
>>>>>>> 03397665 (update api)
		}
	}
	return nil
}

// shouldTraverse reports whether the symlink fi, found in dir,
// should be followed.  It makes sure symlinks were never visited
// before to avoid symlink loops.
func (w *walker) shouldTraverse(dir string, fi os.FileInfo) bool {
	path := filepath.Join(dir, fi.Name())
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}
	ts, err := os.Stat(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	if !ts.IsDir() {
		return false
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	if w.shouldSkipDir(ts, dir) {
=======
	if w.shouldSkipDir(ts) {
>>>>>>> 79bfea2d (update vendor)
=======
	if w.shouldSkipDir(ts, dir) {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	if w.shouldSkipDir(ts, dir) {
>>>>>>> 03397665 (update api)
		return false
	}
	// Check for symlink loops by statting each directory component
	// and seeing if any are the same file as ts.
	for {
		parent := filepath.Dir(path)
		if parent == path {
			// Made it to the root without seeing a cycle.
			// Use this symlink.
			return true
		}
		parentInfo, err := os.Stat(parent)
		if err != nil {
			return false
		}
		if os.SameFile(ts, parentInfo) {
			// Cycle. Don't traverse.
			return false
		}
		path = parent
	}

}
