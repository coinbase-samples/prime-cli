# Changelog

## [0.3.0] - 2025-MAY-05

### Added

- New version command
- Pagination support for list commands
  - Added `--all` and `--interactive` command-line arguments
  - Removed the `--cursor` argument
- **Breaking change:** Commands are now grouped by function
  - E.g., `primectl create-order` is now `primectl orders create`
- **Breaking change:** Removed all command-line argument abbreviations
