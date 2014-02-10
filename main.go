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
	"sort"
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

// getSourceImports runs "go list" with the specified path,
// and then gathers the imports for each listed package.
func getSourceImports(path string) (map[string][]string, error) {
	packages, err := listPackages(path)
	if err != nil {
		return nil, err
	}
	imports := make(map[string][]string)
	for _, path := range packages {
		pkg, err := build.Import(path, ".", 0)
		if err != nil {
			return nil, err
		}
		imports[path] = pkg.Imports
		if testImports {
			imports[path] = append(imports[path], pkg.TestImports...)
			imports[path] = append(imports[path], pkg.XTestImports...)
		}
		sort.Strings(imports[path])
	}
	return imports, nil
}

func Main() error {
	imports, err := getSourceImports(flag.Args()[0])
	if err != nil {
		return err
	}
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
	for path, imports := range imports {
		if targets[path] {
			fmt.Println(path)
		} else {
			for _, imp := range imports {
				if targets[imp] {
					fmt.Println(path)
					break
				}
			}
		}
	}
	return nil
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
