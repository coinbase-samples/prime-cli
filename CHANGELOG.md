# Changelog

## [0.5.0] - 2026-JUN-24

### Added

- Bumped `prime-sdk-go` to v0.9.0 (module path: `github.com/coinbase/prime-sdk-go`)
- New financing commands: `get-cross-margin-risk-parameters`, `get-cross-margin-prime-overview`, `get-market-data`, `update-funding-settings`
- Example script for `orders create-preview`

### Fixed

- `orders create-preview` now calls `CreateOrderPreview` instead of `CreateOrder`, so previews no longer submit live orders

## [0.4.1] - 2026-MAY-01

### Added

- Bumped `prime-sdk-go` to v0.6.3

## [0.4.0] - 2026-APR-28

### Added

- Bumped `prime-sdk-go` to v0.6.2
- New top-level command groups: `advanced-transfers`, `futures`, `positions`
- 36 new commands across 11 service areas (see README for full list)

### Fixed

- `financing create-locate` correctly marks the `date` flag as required
- `wallets create-deposit-address` no longer references unrelated wallet-create flags

## [0.3.0] - 2025-MAY-05

### Added

- New version command
- Pagination support for list commands
  - Added `--all` and `--interactive` command-line arguments
  - Removed the `--cursor` argument
- **Breaking change:** Commands are now grouped by function
  - E.g., `primectl create-order` is now `primectl orders create`
- **Breaking change:** Removed all command-line argument abbreviations
