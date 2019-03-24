# go-diff [![Build Status](https://travis-ci.org/sourcegraph/go-diff.svg?branch=master)](https://travis-ci.org/sourcegraph/go-diff) [![GoDoc](https://godoc.org/github.com/sourcegraph/go-diff/diff?status.svg)](https://godoc.org/github.com/sourcegraph/go-diff/diff)

Diff parser and printer for Go.

Installing
----------

```bash
go get -u github.com/sourcegraph/go-diff/diff
```

Usage
-----

It doesn't actually compute a diff. It only reads in (and prints out, given a Go struct representation) unified diff output, such as the following. The corresponding data structure in Go is the `diff.FileDiff` struct.

```diff
--- oldname	2009-10-11 15:12:20.000000000 -0700
+++ newname	2009-10-11 15:12:30.000000000 -0700
@@ -1,3 +1,9 @@ Section Header
+This is an important
+notice! It should
+therefore be located at
+the beginning of this
+document!
+
 This part of the
 document has stayed the
 same from version to
@@ -5,16 +11,10 @@
 be shown if it doesn't
 change.  Otherwise, that
 would not be helping to
-compress the size of the
-changes.
-
-This paragraph contains
-text that is outdated.
-It will be deleted in the
-near future.
+compress anything.
```
