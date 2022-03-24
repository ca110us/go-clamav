# go-clamav
[![GoDoc](https://pkg.go.dev/badge/github.com/ca110us/go-clamav?status.svg)](https://pkg.go.dev/github.com/ca110us/go-clamav?tab=doc)

go-clamav is go wrapper for [libclamav](https://docs.clamav.net/manual/Development/libclamav.html)

## Environment
### Ubuntu

```bash
apt-get update && apt-get install -y \
  `# install tools` \
  gcc make pkg-config python3 python3-pip python3-pytest valgrind \
  `# install clamav dependencies` \
  check libbz2-dev libcurl4-openssl-dev libjson-c-dev libmilter-dev \
  libncurses5-dev libpcre2-dev libssl-dev libxml2-dev zlib1g-dev

  python3 -m pip install --user cmake / apt-get install cmake
```

Download the source from the clamav [downloads page](https://www.clamav.net/downloads)

```bash
tar xzf clamav-[ver].tar.gz
cd clamav-[ver]

mkdir build && cd build

cmake ..
cmake --build .
ctest
sudo cmake --build . --target install
```

For other Linux distributions, see [clamav documentation](https://docs.clamav.net/manual/Installing/Installing-from-source-Unix.html)

## Quick Start
```bash
$ cd example && cat main.go
```

```go
package main

import (
	"fmt"

	clamav "github.com/ca110us/go-clamav"
)

func main() {
	// new clamav instance
	c := new(clamav.Clamav)
	err := c.Init(clamav.SCAN_OPTIONS{
		General:   0,
		Parse:     clamav.CL_SCAN_PARSE_ARCHIVE | clamav.CL_SCAN_PARSE_ELF,
		Heuristic: 0,
		Mail:      0,
		Dev:       0,
	})

	if err != nil {
		panic(err)
	}

	// free clamav memory
	defer c.Free()

	// load db
	signo, err := c.LoadDB("./db", uint(clamav.CL_DB_DIRECTORY))
	if err != nil {
		panic(err)
	}
	fmt.Println("db load succeed:", signo)

	// compile engine
	err = c.CompileEngine()
	if err != nil {
		panic(err)
	}

	c.EngineSetNum(clamav.CL_ENGINE_MAX_SCANSIZE, 1024*1024*40)
	c.EngineSetNum(clamav.CL_ENGINE_MAX_SCANTIME, 9000)
	// fmt.Println(c.EngineGetNum(clamav.CL_ENGINE_MAX_SCANSIZE))

	// scan
	scanned, virusName, ret := c.ScanFile("./test_file/nmap")
	fmt.Println(scanned, virusName, ret)
}
```

```bash
$ go run main.go

db load succeed: 9263
209 YARA.Unix_Packer_UpxDetail.UNOFFICIAL Virus(es) detected
```

## Reference
[mirtchovski/clamav](https://github.com/mirtchovski/clamav)

*This project was written because `mirtchovski/clamav` no longer supports the new version `clamav`*