# go-diff

Diff parser and printer for Go.

It doesn't actually compute a diff. It only reads in (and prints out,
given a Go struct representation) unified diff output, such as:

```diff
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
