#!/bin/bash

VERSION=$1
WORKING_DIR=$(pwd)
RELEASE_PATH=/root/tool/release

if [[ -z "$VERSION" ]]; then
    VERSION="1.0.0.1"
    echo "ENTRYPOINT NOT FOUND"
else
    echo "using version ${VERSION}"
fi

function build(){
    cd $1
    /root/tool/buildwin.sh
    /root/tool/buildunix.sh
    cd ${WORKING_DIR}
}



function create_hash(){
    
    
    echo "checksum version ${VERSION}"
    cd $1

    CHECKSUM=/tmp/checksum.txt

    rm -f ${CHECKSUM}
    rm -f checksum.txt

    echo "sha256:" >> ${CHECKSUM}
    echo "------------------------------------" >> ${CHECKSUM}
    shasum * >> ${CHECKSUM}
    echo "------------------------------------" >> ${CHECKSUM}
    echo "openssl-sha256:" >> ${CHECKSUM}
    echo "------------------------------------" >> ${CHECKSUM}
    sha256sum * >> ${CHECKSUM}
    cat ${CHECKSUM}
    mv ${CHECKSUM} checksum.txt


    cd ${WORKING_DIR}

    # cd osin-compress
    # echo "sha256: `shasum osin-win-1.38.19.84.zip`" >> checksums.txt
    # echo "openssl-sha256: `sha256sum osin-win-1.38.19.84.zip`" >> checksums.txt
    # echo "sha256: `shasum osin-win-not_strip-1.38.19.84.zip`" >> checksums.txt
    # echo "openssl-sha256: `sha256sum osin-win-not_strip-1.38.19.84.zip`" >> checksums.txt
    # cat checksums.txt
    # cd ..
}

function create_win_release() {
    NAME="Windows_v${VERSION}.zip"

    cd $1
    zip -r ${RELEASE_PATH}/${NAME} .
    cd ${WORKING_DIR}
}

function create_unix_release() {
    NAME="x86_64-linux-gnu_v${VERSION}.tar.gz"
    cd $1
    tar -cvzf ${RELEASE_PATH}/${NAME} *
    cd ${WORKING_DIR}
}


echo "building coin"
build /root/coin

echo "resetting folder ${RELEASE_PATH}"
rm -rf  ${RELEASE_PATH}
mkdir -p ${RELEASE_PATH}

echo "calculate hash"
create_hash /root/tool/unixdist
create_hash /root/tool/windist

create_win_release /root/tool/windist
create_unix_release /root/tool/unixdist

create_hash ${RELEASE_PATH}


