rdep
====

Reverse dependency lister for Go

Building
========

go get github.com/axw/rdep

Using
=====

Specify the source package (list) as the first argument;
subsequent arguments (targets) are the packages which
we are testing for dependence on.

e.g. to list all packages under the current working
directory that import os or runtime, run:

    rdep ./... os runtime

The target packages are always included in the output.

