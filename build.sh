#!/bin/bash

function beforeBuild() {
    rm -rf ./build
    mkdir -p ./build/release
}

function buildbsyslog() {
    mkdir -p ./build/release/bsyslog/config
    mkdir -p ./build/release/bsyslog/ca
    cp -rf ./src/bsyslog* ././build/release/bsyslog
    cp -rf ./conf/config.yaml ./build/release/bsyslog/config
    cp -rf ./conf/bsyslog.service ./build/release/bsyslog/config
    cp -rf ./conf/ca/* ./build/release/bsyslog/ca
}

function buildbaudit() {
    mkdir -p ./build/release/baudit/config
    mkdir -p ./build/release/baudit/
    cp -rf ./dependency/* ./build/release/baudit/
}

function builddemo() {
    mkdir -p ./build/release/demo/ca
    cp -rf ./src/syslog_server ./build/release/demo/
    cp -rf ./demo/ca/* ./build/release/demo/ca
}

function afterBuild() {
    cp ./install.sh ./build/release
    cd ./build/
    tar zcvf release.tar.gz ./release
    rm -rf ./release
    echo "tar release.tar.gz"
}

function build() {
    beforeBuild
    make clean; make -j
    sleep 1
    buildbsyslog
    buildbaudit
    builddemo
    afterBuild
}

build