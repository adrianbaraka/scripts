#!/usr/bin/env bash

cvt 1920 1080

xrandr --newmode "1920x1080_60.00"  173.00  1920 2048 2248 2576  1080 1083 1088 1120 -hsync +vsync

xrandr --addmode DP-2 "1920x1080_60.00"


xrandr --output DP-2 --mode "1920x1080_60.00"

