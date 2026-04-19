# mkvclean

`mkvclean` is a Go-based CLI utility designed to sanitize Matroska (`.mkv`) files by stripping unnecessary metadata and attachments using `mkvpropedit`.

## Features
* **Remove Attachments:** Deletes all embedded fonts and images.
* **Strip Titles:** Clears metadata titles from video, audio, and subtitle tracks.
* **Clean Global Data:** Removes global tags and track statistics.
* **Batch Processing:** Supports individual file cleaning or recursive folder scanning.

---

## Code Structure

The project follows a standard Go CLI layout using the **Cobra** library:

* **`main.go`**: The entry point that initializes and executes the root command.
* **`cmd/`**: Contains the command logic and flag definitions.
    * **`root.go`**: Defines the base persistent flags (verbosity, color, pathing).
    * **`file.go`**: Implements logic for processing specific file arguments.
    * **`folder.go`**: Implements recursive directory walking to find and process MKV files.
    * **`utils.go`**: Contains shared helper functions for executing `mkvpropedit` and handling output formatting.
* **`go.mod` / `go.sum`**: Manages project dependencies.

---

## Usage Summary

```bash
# Clean specific files
mkvclean file movie1.mkv movie2.mkv

# Clean a directory recursively
mkvclean folder ./library --language eng
```

*Note: Requires `mkvpropedit` (part of MKVToolNix) to be installed in your `$PATH`.*