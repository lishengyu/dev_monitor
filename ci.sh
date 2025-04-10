#!/bin/bash

VERSION_FILE=VERSION.txt

function build() {
    echo "Building...."
    
    COMPILE_TIME=$(date "+%Y%m%d%H%M%S")
    VER_NEW=$(cat ${VERSION_FILE})
    bash ./build.sh
    mv ./build/release.tar.gz ./build/bsyslog_${VER_NEW}_${COMPILE_TIME}.tar.gz
    sleep 1
    echo "BUild done"
}

build