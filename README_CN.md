# go-clamav

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

  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
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

其他 Linux 发行版参照 [clamav documentation](https://docs.clamav.net/manual/Installing/Installing-from-source-Unix.html)

## 快速开始
参考 `example` 目录

## 参考
[mirtchovski/clamav](https://github.com/mirtchovski/clamav)

*因为 `mirtchovski/clamav` 不再支持新版本 `clamav`，所以写了该项目*
