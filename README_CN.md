# go-clamav
[![GoDoc](https://pkg.go.dev/badge/github.com/ca110us/go-clamav?status.svg)](https://pkg.go.dev/github.com/ca110us/go-clamav?tab=doc)

go-clamav 是 go 语言对 [libclamav](https://docs.clamav.net/manual/Development/libclamav.html) 的封装

## 环境
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

从 clamav 官方下载源码 [downloads page](https://www.clamav.net/downloads)

```bash
tar xzf clamav-[ver].tar.gz
cd clamav-[ver]

mkdir build && cd build

cmake ..
cmake --build .
ctest
sudo cmake --build . --target install
```

其他 Linux 发行版参照 [clamav documentation](https://docs.clamav.net/manual/Installing/Installing-from-source-Unix.html)

## 快速开始
### 动态链接
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
$ CGO_LDFLAGS="-L/usr/local/lib -lclamav" go run main.go

db load succeed: 9263
209 YARA.Unix_Packer_UpxDetail.UNOFFICIAL Virus(es) detected
```

如果找不到 `libclamav.so` 文件，尝试如下：

```bash
$ LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib CGO_LDFLAGS="-L/usr/local/lib -lclamav" go run main.go

db load succeed: 9263
209 YARA.Unix_Packer_UpxDetail.UNOFFICIAL Virus(es) detected
```

### 静态编译
```bash
$ sudo bash ./prepare.sh
$ SRCDIR=$(pwd)
$ export CGO_CFLAGS="-g -Wall -I${SRCDIR}/clamav-mussels-cookbook/mussels/install/include"
$ export CGO_LDFLAGS="-L${SRCDIR}/clamav-mussels-cookbook/mussels/install/lib -lclamav_static -lbz2_static -lclammspack_static -lclamunrar_iface_static -lclamunrar_static -lcrypto -ljson-c -lpcre2-8 -lpcre2-posix -lssl -lxml2 -lz -lm -ldl -lstdc++"
$ CGO_ENABLED=1 go build --ldflags '--extldflags "-static -fpic"' main.go
```

## 参考
[mirtchovski/clamav](https://github.com/mirtchovski/clamav)

*因为 `mirtchovski/clamav` 不再支持新版本 `clamav`，所以写了该项目*
