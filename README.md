<h1 align="center">go-vss</h1>

<p  align="center">
 <a href="https://github.com/forensicanalysis/go-vss/actions"><img src="https://github.com/forensicanalysis/go-vss/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/go-vss"><img src="https://codecov.io/gh/forensicanalysis/go-vss/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/go-vss"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/go-vss" alt="report" /></a>
 <a href="https://pkg.go.dev/github.com/forensicanalysis/go-vss"><img src="https://img.shields.io/badge/go.dev-documentation-007d9c?logo=go&logoColor=white" alt="doc" /></a>
</p>


This projects adds bindings for [libvshadow](https://github.com/libyal/libvshadow) for Go for macos, linux and windows systems.
go-vss parses a raw VSS file into a [NTFS file system](https://github.com/forensicanalysis/fslib/tree/master/filesystem/ntfs).

# Usage

In general, there are different ways to load a raw vss filesystem:
1. Using a go reader object
2. Using a file descriptor
3. Using a path

Make sure, you have installed the package:
```bash
go get github.com/forensicanalysis/go-vss
```

## Using a go reader object
```go
package main

import (
	"fmt"
	"os"

	"github.com/forensicanalysis/go-vss"
)

func main() {
	reader, err := os.Open("vss.raw")

	stores, err := vss.NewByReader(reader)

	for _, fs := range stores {
		root, _ := fs.Open("/")
		filenames, _ := root.Readdirnames(0)
		fmt.Println(filenames)
	}
}
```

## Using a file descriptor
```go
package main

import (
	"fmt"
	"syscall"

	"github.com/forensicanalysis/go-vss"
)

func main() {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC, 755)

	stores, err := vss.NewByFd(fd)

	for _, fs := range stores {
		root, _ := fs.Open("/")
		filenames, _ := root.Readdirnames(0)
		fmt.Println(filenames)
	}
}
```

## Using a path
```go
package main

import (
	"fmt"

	"github.com/forensicanalysis/go-vss"
)

func main() {
	path := "vss.raw"

	stores, err := vss.NewByPath(path)

	for _, fs := range stores {
		root, _ := fs.Open("/")
		filenames, _ := root.Readdirnames(0)
		fmt.Println(filenames)
	}
}
```

# Updating libvshadow
This project builds upon the libvshadow created by Joachim Metz. The prebuild static libraries of libvshadow are placed in the subfolders ```libvshadow_prebuild\windows```, ```libvshadow_prebuild\darwin``` and ```libvshadow_prebuild\linux```.

The current version used in the project is: ```20191221```

If you need to update the library use the following instructions to build new static libraries:

1. Download the required ```libvshadow-[new-version].tar.gz``` from the [relases](https://github.com/libyal/libvshadow/releases).
2. Remove the current ```libvshadow-[old-version].tar.gz``` and place the new ```libvshadow-[new-version].tar.gz``` into the root directory of the repository.
3. Adjust the ```Makefile``` and ```.github/workflows```: Make sure to update the version numbers in order to use the updated ```libvshadow-[new-version].tar.gz```.
4. Commit and push the directory. Github actions creates artifacts containing the static libraries. Once all actions succeeded, you can download the static libvshadow libraries for each operating system: macos, linux and windows.
5. Extract the static libs (all of them are called libvshadow.a) to the correct destination:
    - libvshadow_windows -> ```libvshadow_prebuild\windows``` 
    - libvshadow_macos -> ```libvshadow_prebuild\macos```
    - libvshadow_linux -> ```libvshadow_prebuild\linux```
