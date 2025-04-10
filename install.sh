#!/bin/bash

SERVICE_PATH=/lib/systemd/system

function installBaudit() {
    if [ -e "/home/baudit" ]; then
        systemctl stop baudit
        rm -rf /home/baudit_bak
        mv /home/baudit /home/baudit_bak
    fi
    
    echo "install baudit..."
    cp -rf ./baudit/config/baudit.service ${SERVICE_PATH}
        systemctl enable baudit
    cp -rf ./baudit /home/
    chmod +x /home/baudit/go-audit
}

function installBsyslog() {
    if [ -e "/home/bsyslog" ]; then
        systemctl stop bsyslog
        rm -rf /home/bsyslog_bak
        mv /home/bsyslog /home/bsyslog_bak
    fi
    
    echo "install bsyslog..."
    cp -rf ./bsyslog/config/bsyslog.service ${SERVICE_PATH}
    systemctl enable bsyslog
    cp -rf ./bsyslog /home/
    chmod +x /home/bsyslog/bsyslog
}

function InstallDemo() {
echo "install demo..."
    cp -rf demo /home/bsyslog/
}

function install() {
    installBaudit
    installBsyslog
    InstallDemo
}

install