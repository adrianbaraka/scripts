# mkvaudio

A CLI tool to recursively scan and downmix multi-channel `.mkv` audio tracks into stereo, optimizing your media library's storage footprint.

- Automatically detects and downmixes audio streams to stereo if they aren't already.
- Results in an approximate **10–20% file size reduction** depending on the original audio codec.


## Example Usage

```bash
# Preview changes without modifying files
mkvaudio process folder "/videofiles" --dry-run

# Downmix to stereo, keep a backup of original files, and enable verbose output
mkvaudio process folder "/videofiles" --backup -v

# Downmix to stereo, overwrite original files, and log errors only
mkvaudio process folder "/videofiles" -q

```

## TODOs

* Show the actual file/folder sizes before and after in the final summary output.
* Add a `-h, --human-readable` flag to format file sizes (e.g., `4.2 GB` instead of raw bytes).
* Add support to differentiate between binary (`KiB`, `MiB`, `GiB`) and decimal (`KB`, `MB`, `GB`) units.
