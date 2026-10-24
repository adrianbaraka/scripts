#!/usr/bin/env bash

#set env variable for ps1
echo export IS_REMOTE="True" >> "$HOME"/.bashrc

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
echo "PATH=$PATH:$SCRIPT_DIR/scripts" >> "$HOME"/.bashrc

#bashrc
while read -r line; do
    if [ "$line" = "# [Home PC]" ]; then
        break
    else
        echo "$line" >> "$HOME"/.bashrc
    fi
done < ./bash_func




