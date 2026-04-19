# scripts

A collection of various utility scripts mostly for a homelab.


## Installation Script

The repository includes a helper bash script to automate the building and installation of the various Go utilities. This script handles compilation, binary placement, and shell completion setup.

### What it does
* **Tidies Dependencies:** Runs `go mod tidy` to ensure all requirements are met.
* **Compiles:** Builds the Go source into a binary named after the project.
* **Sets up Autocompletion:** Automatically generates and installs Bash completion scripts to `~/.local/share/bash-completion/completions`.
* **Installs to PATH:** Moves the compiled binary to `~/.local/bin/`.

---

### Usage

To install a specific utility, run the script from the root of the repository and pass the directory name as an argument

Current binaries: 
- mkvclean
- mkvsubs
- rename

```bash
./install.sh mkvclean
```

#### Requirements
1. **Go:** Ensure the Go toolchain is installed.
2. **PATH:** Ensure `~/.local/bin` is in your system `$PATH`. You can add it by adding this line to your `.bashrc` or `.zshrc`:
   ```bash
   export PATH="$HOME/.local/bin:$PATH"
   ```
3. **Bash Completion:** To use the generated completions immediately, source your bashrc or restart your terminal:
   ```bash
   source ~/.bashrc
   ```

> [!NOTE]  
> Currently, autocompletion setup is only implemented for **Bash**. Zsh and Fish support are planned for future updates.