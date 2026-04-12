#!/usr/bin/env bash

XORG_CONF="/usr/share/X11/xorg.conf.d/20-intel.conf"

if [[ ! -e "$XORG_CONF" ]]; then
    echo "Config not found. Installing $XORG_CONF..."
    sudo cp 20-intel.conf "$XORG_CONF"
    echo "Success. Please reboot, then run this script again."
    exit 0
fi

source res.conf || exit 1

if xrandr | grep -q "VIRTUAL1"; then
    xrandr --addmode VIRTUAL1 "$screen1"
    xrandr --output VIRTUAL1 --mode "$screen1" --right-of eDP1
    echo "Display VIRTUAL1 started with resolution $screen1"
else
    echo "Error: VIRTUAL1 not detected by xrandr. Is the Intel driver loaded?"
    exit 1
fi

xrandr --query | grep "VIRTUAL1"
