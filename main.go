// Copyright 2014 Andrew Wilkins <axwalk@gmail.com>
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"strings"
)

var (
	testImports bool
)

func init() {
	flag.BoolVar(&testImports, "tests", false, "enable inclusion of tests when gathering imports")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s <source> <target> [target...]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// listPackages returns the result of running "go list"
// with the specified path.
func listPackages(path string) ([]string, error) {
	cmd := exec.Command("go", "list", path)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(out)), nil
}

// getPackages imports and returns a build.Package for each listed package.
func getPackages(paths []string) (map[string]*build.Package, error) {
	packages := make(map[string]*build.Package)
	for _, path := range paths {
		pkg, err := build.Import(path, ".", 0)
		if err != nil {
			return nil, err
		}
		packages[path] = pkg
	}
	return packages, nil
}

func Main() error {
	targets := make(map[string]bool)
	for _, arg := range flag.Args()[1:] {
		packages, err := listPackages(arg)
		if err != nil {
			return err
		}
		for _, path := range packages {
			targets[path] = true
		}
	}
	paths, err := listPackages(flag.Args()[0])
	if err != nil {
		return err
	}
	packages, err := getPackages(paths)
	if err != nil {
		return err
	}
	for path, _ := range packages {
		if imports(path, packages, targets, testImports) {
			fmt.Println(path)
		}
	}
	return nil
}

// imports returns true if path imports any
// of the packages in "any", transitively.
func imports(path string, packages map[string]*build.Package, any map[string]bool, testImports bool) (res bool) {
	if any[path] {
		return true
	}
	pkg, _ := packages[path]
	if pkg == nil {
		return false
	}
	if testImports {
		for _, imp := range pkg.TestImports {
			if any[imp] {
				return true
			}
		}
		for _, imp := range pkg.XTestImports {
			if any[imp] {
				return true
			}
		}
	}
	for _, imp := range pkg.Imports {
		if imports(imp, packages, any, false) {
			any[path] = true
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
