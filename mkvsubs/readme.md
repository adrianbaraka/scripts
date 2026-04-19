# mkvsubs

`mkvsubs` is a CLI tool designed to manage and refine subtitle tracks within Matroska (`.mkv`) containers. It leverages `mkvmerge` and `mkvextract` to extract, offset, merge, and clean subtitle data.

## Features
* **Subtitle Extraction & Merging:** Extract internal tracks or merge external subtitle files into existing MKV containers.
* **Timing Adjustments:** Apply millisecond-level delays (positive or negative) to subtitle synchronization.
* **Smart Cleaning:** Option to delete image-based subtitles (VobSub) or strip extra subtitle tracks to keep files lean.
* **Safety First:** Includes a `--dry-run` mode to preview changes and a `--backup` flag to preserve original files.
* **Automation:** Use `--merge-scan` to automatically find and attach local subtitle files matching the video filename.

---

## Code Structure

The project is organized into modular packages to separate CLI handling from core subtitle logic:

* **`main.go`**: The application entry point.
* **`cmd/`**: CLI command definitions using Cobra.
    * **`root.go`**: Sets up global flags (backup, verbosity, binary paths).
    * **`process.go`**: High-level orchestration of the processing pipeline.
    * **`file.go` / `folder.go`**: Handlers for individual file targets or recursive directory scanning.
    * **`util.go`**: Helper functions for CLI output and execution.
* **`subs/`**: Core library for subtitle manipulation.
    * **`types.go` / `caseTypes.go`**: Data models for subtitle tracks and processing configurations.
    * **`subs.go`**: The primary logic for interacting with MKVToolNix binaries.
    * **`*_test.go`**: Test suites using JSON-based test data in `testdata/`.

---

## Usage

```bash
# Apply a 300ms delay to the first subtitle track
mkvsubs process file movie.mkv --subtitle-number 1 --delay 300

# Automatically find and merge matching external .srt files in a folder
mkvsubs process folder ./tv-shows --merge-scan --delete-image-subs
```

---

## Todo

- **Enhanced Delay Support:** Add the delay option to both external and internal srt files.
- **Localization:** Add a language option for the new subtitle track to ensure proper metadata tagging.