// Copyright 2016 The plugs Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package plugs provides an API to compile and open plugins.
//
// A simple usage looks like the following:
//
//  p, err := plugs.Open("github.com/sbinet/plugs/testdata/p1")
//  if err != nil {
//      panic(err)
//  }
//  v, err := p.Lookup("V")
package plugs

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"
)

const (
	plugExt = ".so"
)

var mgr = struct {
	sync.RWMutex
	plugins map[string]*plugin.Plugin
}{
	plugins: make(map[string]*plugin.Plugin),
}

// Open compiles and opens a Go plugin.
// If a package path has already been successfully compiled and opened,
// then the existing *plugin.Plugin is returned.
//
// It is safe for concurrent use by multiple goroutines.
func Open(pkg string) (*plugin.Plugin, error) {
	mgr.Lock()
	defer mgr.Unlock()
	if p, ok := mgr.plugins[pkg]; ok {
		return p, nil
	}

	return buildPlugin(pkg)
}

func buildPlugin(path string) (*plugin.Plugin, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	bpkg, err := build.Import(path, cwd, 0)
	if err != nil {
		return nil, err
	}

	odir := filepath.Join(bpkg.PkgTargetRoot, bpkg.ImportPath)
	err = os.MkdirAll(odir, 0755)
	if err != nil {
		return nil, err
	}

	fname := filepath.Join(odir, "_go_plugin"+plugExt)

	cmd := exec.Command(
		"go", "build", "-buildmode=plugin", "-o", fname,
	)
	cmd.Dir = bpkg.Dir

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("plugs: error compiling plugin: %v", err)
	}

	p, err := plugin.Open(fname)
	if err != nil {
		return nil, err
	}
	mgr.plugins[path] = p

	return p, nil
}
