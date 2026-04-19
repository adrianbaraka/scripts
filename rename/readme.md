# rename

`rename` is a specialized CLI tool designed to batch rename media files into a standardized format, making them compatible with media servers like Plex, Jellyfin.

---

## Code Structure

The project uses a clean, command-based structure built with the **Cobra** library:

* **`main.go`**: The entry point that launches the CLI application.
* **`cmd/`**: Contains the command logic.
    * **`root.go`**: Defines global persistent flags such as color, quiet mode, and verbosity.
    * **`folder.go`**: The core logic for the `folder` command, handling directory traversal and the renaming regex/logic based on user-provided season and name flags.
    * **`util.go`**: Shared utility functions for file system operations and output formatting.
* **`go.mod` / `go.sum`**: Defines Go dependencies and versioning.

---

## Usage

```bash
# Rename all files in a folder to "The Boys S01EXX"
rename folder ./downloads/TheBoys --name "The Boys" --season 01
```

### Flags
* `-n, --name`: The formal title of the series.
* `-s, --season`: The season number (e.g., `01`, `02`).