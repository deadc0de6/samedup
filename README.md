# samedup

*[samedup](https://github.com/deadc0de6/samedup) is a file duplicates finder*

It is safe and will **not** modify or remove any file on your filesystem.

Features:

* Ability to filter files by pattern
* Provides different hash methods (sha1, xxhash, crc32, md5)
* Allow to ignore specific patterns when searching for duplicates
* Different output formats (tree, csv, ...)
* Can create a editable shell script to remove/handle duplicates

---

**Table of Contents**

* [Installation](#installation)
* [Usage](#usage)

  * [Filter patterns](#filter-patterns)
  * [Ignore patterns](#ignore-patterns)
  * [Output formats](#output-formats)

* [Contribution](#contribution)
* [Thank you](#thank-you)

# Installation

```bash
## You need at least golang 1.20
$ go install -v github.com/deadc0de6/samedup/cmd/samedup@latest
$ samedup
```

Compilation (go 1.20 and above)
```bash
$ go mod tidy
$ make
$ ./bin/samedup --help
```

# Usage

Search for duplicates
```bash
# basic use
$ samedup dir1 dir2

# find duplicates only among .txt files
$ samedup --filter '.*\.txt' /usr

# ignore dotfiles
$ samedup --ignore '\..*' .
```

Note that `samedup` will **not** follow symlinks

## Filter patterns

Filter pattern is done using [re2 syntax](https://github.com/google/re2/wiki/Syntax).
The below example will only find duplicates among files ending with `.go` or `.md`
```bash
$ samedup --filter=".*\.go" --filter=".*\.md" dir1 dir2
```

## Ignore patterns

Ignore pattern is done using [re2 syntax](https://github.com/google/re2/wiki/Syntax).
The below example will ignore any file in a `.git` directory.
```bash
# compare but ignore any .git/ directory
$ samedup --ignore=".*\.git" dir1 dir2
```

## Output formats

`tree` output format
```bash
$ samedup --output=tree .git
└─┬sha1:ea50fe35ffe6ebcdee543fb3c0ed46c88c0bf150 (2)
  ├──/tmp/samedup/.git/logs/HEAD (size:1.6kB)
  └──/tmp/samedup/.git/logs/refs/heads/master (size:1.6kB)

└─┬sha1:36c81ec49edbac1050262bd69e9c019fbad2b902 (3)
  ├──/tmp/samedup/.git/ORIG_HEAD (size:41B)
  ├──/tmp/samedup/.git/refs/heads/master (size:41B)
  └──/tmp/samedup/.git/refs/remotes/origin/master (size:41B)
```

`csv` output format: `<path>,<checksum>,<mode>,<size>`
```bash
$ samedup --output=csv .git
/tmp/samedup/.git/ORIG_HEAD,sha1:36c81ec49edbac1050262bd69e9c019fbad2b902,-rw-r--r--,41B
/tmp/samedup/.git/refs/heads/master,sha1:36c81ec49edbac1050262bd69e9c019fbad2b902,-rw-r--r--,41B
/tmp/samedup/.git/refs/remotes/origin/master,sha1:36c81ec49edbac1050262bd69e9c019fbad2b902,-rw-r--r--,41B
/tmp/samedup/.git/logs/HEAD,sha1:ea50fe35ffe6ebcdee543fb3c0ed46c88c0bf150,-rw-r--r--,1.6kB
/tmp/samedup/.git/logs/refs/heads/master,sha1:ea50fe35ffe6ebcdee543fb3c0ed46c88c0bf150,-rw-r--r--,1.6kB
```

`stairs` output format:
```bash
$ samedup --output=stairs .git
sha1:ea50fe35ffe6ebcdee543fb3c0ed46c88c0bf150
  /tmp/samedup/.git/logs/HEAD
  /tmp/samedup/.git/logs/refs/heads/master
sha1:36c81ec49edbac1050262bd69e9c019fbad2b902
  /tmp/samedup/.git/ORIG_HEAD
  /tmp/samedup/.git/refs/heads/master
  /tmp/samedup/.git/refs/remotes/origin/master
```

`oneline` output format:
```bash
$ samedup --output=oneline .git
/tmp/samedup/.git/ORIG_HEAD /tmp/samedup/.git/refs/heads/master /tmp/samedup/.git/refs/remotes/origin/master
/tmp/samedup/.git/logs/HEAD /tmp/samedup/.git/logs/refs/heads/master
```

`nlines` output format:
```bash
$ samedup --output=nlines .git
/tmp/samedup/.git/logs/HEAD
/tmp/samedup/.git/logs/refs/heads/master

/tmp/samedup/.git/ORIG_HEAD
/tmp/samedup/.git/refs/heads/master
/tmp/samedup/.git/refs/remotes/origin/master
```

`script` output format:
```bash
$ samedup --output=script .git
#!/usr/bin/env bash
#
# total 2 duplicates, total wasted: 5.4kB
#
# 3 duplicates for "sha1:ac06fd0d9d50c025bd9c612369e1b889af44587b" - would free 82B
#rm -fv '/tmp/samedup/.git/ORIG_HEAD'
rm -fv '/tmp/samedup/.git/refs/heads/master'
rm -fv '/tmp/samedup/.git/refs/remotes/origin/master'

# 2 duplicates for "sha1:e228e22e2751aa32779e5c6c5775c244829d34eb" - would free 5.4kB
#rm -fv '/tmp/samedup/.git/logs/HEAD'
rm -fv '/tmp/samedup/.git/logs/refs/heads/master'
```

# Contribution

If you are having trouble installing or using `samedup`, open an issue.

If you want to contribute, feel free to do a PR.

The `test.sh` script handles the linting and runs the tests for the code.

# Thank you

If you like `samedup`, [buy me a coffee](https://ko-fi.com/deadc0de6).

# License

This project is licensed under the terms of the GPLv3 license.
