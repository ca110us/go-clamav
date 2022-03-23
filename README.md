# go-clamav

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

For other Linux distributions, see [clamav documentation](https://docs.clamav.net/manual/Installing/Installing-from-source-Unix.html)

## Quick Start
Refer to the `example` directory

## Reference
[mirtchovski/clamav](https://github.com/mirtchovski/clamav)

*This project was written because `mirtchovski/clamav` no longer supports the new version `clamav`*