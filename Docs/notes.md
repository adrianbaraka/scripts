# notes

## todos

- Add a doctor command and validate existence of mkvtools used before use.
- ReadMe's for each of the tools.
- Refine the installation script.

## some tips when using cobra

- The rootCmd should have a persistentPreRun to ensure things like the logger are initialized. This applies to any parent command that executes.
