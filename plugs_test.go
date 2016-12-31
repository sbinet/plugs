// Copyright 2016 The plugs Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugs

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"
)

func TestOpen(t *testing.T) {
	const pname = "github.com/sbinet/plugs/testdata/p1"

	gopath := filepath.SplitList(build.Default.GOPATH)[0]
	pdir := filepath.Join(
		gopath,
		build.Default.GOOS+"_"+build.Default.GOARCH,
		pname,
	)
	os.Remove(filepath.Join(pdir, "_go_plugin."+plugExt))

	p, err := Open(pname)
	if err != nil {
		t.Fatal(err)
	}
	v, err := p.Lookup("V")
	if err != nil {
		t.Errorf("error looking up V: %v\n", err)
	}
	f, err := p.Lookup("F")
	if err != nil {
		t.Fatalf("error looking up F: %v\n", err)
	}
	*v.(*int) = 7
	if got, want := f.(func() string)(), "Hello, number 7"; got != want {
		t.Errorf("got %q, want %q\n", got, want)
	}
}
