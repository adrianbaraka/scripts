#!/usr/bin/env bash

krfb-virtualmonitor --resolution 1920x1080 --name sunshine-vm --password password --port 5905 &

sleep 3

kscreen-doctor output.Virtual-sunshine-vm.addCustomMode.1920.1080.60000.full
kscreen-doctor output.Virtual-sunshine-vm.mode.1920x1080@60