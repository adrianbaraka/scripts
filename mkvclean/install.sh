#!/usr/bin/env sh

# exit if something fails
set -e

# install mkvclean to users path and add shell completions
# build the app
go mod tidy

go build -o mkvclean main.go

# shell completions
case "$SHELL" in
    *"bash"*)
    compDir="/etc/bash_completion.d"
        if [ -d "$compDir" ]; then
            ./mkvclean completion bash > mkvclean_bash_comp
             sudo mv mkvclean_bash_comp "$COMP_DIR/mkvclean"
        fi
    ;;
    *"zsh"*)
        echo "zsh completions not yet implemented"
    ;;
    *"fish"*)
        echo "fish completions not yet implemented"
    ;;
    *)
        echo "Unsupported shell: $SHELL" 
    ;;
esac

# install in the users bin 
# TODO systemwide install

destDir="$HOME/.local/bin"
mkdir -p "$destDir"
mv -v mkvclean "$destDir/mkvclean"

echo "Installation complete. Ensure $destDir is in your PATH."
