#!/usr/bin/env bash

# exit if something fails
set -e

# get the dir containing the script
if (( $# != 1 )); then 
    echo "Enter the name of the script to install eg mkvclean" >&2
    exit 1
fi

app=$1
BashCompDir="$HOME/.local/share/bash-completion/completions"
destDir="$HOME/.local/bin"

# system bash completions /etc/bash_completion.d/
# system bin 

(
    cd "$app"
    #install mkvclean to users path and add shell completions
    #build the app
    go mod tidy

    go build -o "$app" main.go
    echo "Compiled $app"

    # shell completions
    case "$SHELL" in
        *"bash"*)
            
            mkdir -p "$BashCompDir"
            ./"$app" completion bash > "${app}_bash_comp"
            mv "${app}_bash_comp" "$BashCompDir/$app"
            echo "Generated and copied the bash completion script to $BashCompDir"
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

    mkdir -p "$destDir"
    mv -v "$app" "$destDir/$app"

    echo "Installation complete. Ensure $destDir is in your PATH."
)