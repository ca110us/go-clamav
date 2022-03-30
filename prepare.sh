#!/bin/bash

# make static lib for clamav deps
apt-get update
apt-get install -y python3-pip
apt-get install -y build-essential clang llvm
apt-get install -y libstdc++6 libstdc++-6-dev
apt-get install -y flex bison python3-dev pkg-config ninja-build
python3 -m pip install mussels

wget https://cmake.org/files/v3.21/cmake-3.21.5-linux-x86_64.tar.gz
tar -zxvf cmake-3.21.5-linux-x86_64.tar.gz -C /usr/local/
ln -s /usr/local/cmake-3.21.5-linux-x86_64/bin/cmake /usr/bin/cmake
rm -rf cmake-3.21.5-linux-x86_64.tar.gz

git clone --depth 1 https://github.com/ca110us/clamav-mussels-cookbook.git
cd clamav-mussels-cookbook
rm -rf  mussels/* &> /dev/null
mkdir mussels &> /dev/null
msl build libclamav_deps -t host-static -w mussels/work -i mussels/install
cd -

# make get clamav source code
git clone https://github.com/Cisco-Talos/clamav.git
cd clamav
git checkout clamav-0.104.0
cd -


# libclamav
cd clamav

rm -rf  ./build/* &> /dev/null
mkdir build &> /dev/null
cd -
export CLAMAV_DEPENDENCIES="$(pwd)/clamav-mussels-cookbook/mussels/install/"

cd clamav/build

cmake .. -G Ninja                                                      \
    -DCMAKE_BUILD_TYPE="Release"                                       \
    -DJSONC_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include/json-c"          \
    -DJSONC_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libjson-c.a"             \
    -DBZIP2_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include"                 \
    -DBZIP2_LIBRARY_RELEASE="$CLAMAV_DEPENDENCIES/lib/libbz2_static.a" \
    -DOPENSSL_ROOT_DIR="$CLAMAV_DEPENDENCIES"                          \
    -DOPENSSL_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include"               \
    -DOPENSSL_CRYPTO_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libcrypto.a"    \
    -DOPENSSL_SSL_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libssl.a"          \
    -DLIBXML2_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include/libxml2"       \
    -DLIBXML2_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libxml2.a"             \
    -DPCRE2_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include"                 \
    -DPCRE2_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libpcre2-8.a"            \
    -DZLIB_INCLUDE_DIR="$CLAMAV_DEPENDENCIES/include"                  \
    -DZLIB_LIBRARY="$CLAMAV_DEPENDENCIES/lib/libz.a"                   \
    -DENABLE_JSON_SHARED=OFF                                           \
    -DENABLE_STATIC_LIB=ON                                             \
    -DENABLE_SYSTEMD=OFF                                               \
    -DENABLE_TESTS=OFF                                                 \
    -DENABLE_LIBCLAMAV_ONLY=ON                                         \
    -DENABLE_UNRAR=ON                                                  \
    -DENABLE_SHARED_LIB=OFF                                            \
    -DDATABASE_DIRECTORY=/var/lib/clamav                               \
    -DCMAKE_INSTALL_PREFIX=install

cmake --build .

cd -

rm -rf ./lib/*
mkdir lib &> /dev/null
cp clamav/build/libclamav/libclamav_static.a ./lib
cp clamav/build/libclammspack/libclammspack_static.a ./lib
cp clamav/build/libclamunrar/libclamunrar_static.a ./lib
cp clamav/build/libclamunrar_iface/libclamunrar_iface_static.a ./lib

cp "$CLAMAV_DEPENDENCIES/lib/libbz2_static.a" ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libjson-c.a" ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libcrypto.a" ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libssl.a" ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libxml2.a"  ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libpcre2-8.a"  ./lib
cp "$CLAMAV_DEPENDENCIES/lib/libz.a" ./lib

rm -rf ./include/*
mkdir include &> /dev/null
cp clamav/build/*.h ./include
cp clamav/libclamav/clamav.h ./include

cp ./include/* clamav-mussels-cookbook/mussels/install/include/
cp ./lib/* clamav-mussels-cookbook/mussels/install/lib/