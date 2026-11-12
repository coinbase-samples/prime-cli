# Changelog

## [0.7.0] - 2026-NOV-06

### Added

- Bumped `prime-sdk-go` to v0.13.0
- New `webhooks` command group (beta): `get-event-types`, `list-subscriptions`, `create-subscription`, `get-subscription`, `update-subscription`, `delete-subscription`, `rotate-hmac-secret`, `send-simulation`, `get-delivery-logs`, `get-delivery-attempts`, `get-subscription-stats`, `get-owner-stats`, `replay-event`
- MCP tools for every webhook endpoint (13 new tools, 126 total): `get_webhook_event_types`, `list_webhook_subscriptions`, `create_webhook_subscription`, `get_webhook_subscription`, `update_webhook_subscription`, `delete_webhook_subscription`, `rotate_webhook_hmac_secret`, `send_webhook_simulation`, `get_webhook_delivery_logs`, `get_webhook_delivery_attempts`, `get_webhook_subscription_stats`, `get_webhook_owner_stats`, `replay_webhook_event`
- `orders create-quote --is-buy-exact` and the matching `is_buy_exact` option on the `create_quote` MCP tool
- `advanced-transfers create --idempotency-key` and the `create_advanced_transfer` MCP `idempotency_key` are now also sent as the request-level idempotency key, so a repeated key returns the original transfer

### Changed

- `webhooks update-subscription` and `update_webhook_subscription` build the update mask from the fields you provide and fetch the current subscription version automatically when `--expected-version` is omitted
- Order, margin, FCM, derivative position, and transfer responses include the new SDK fields (SPREAD order and option preview fields, margin health and thresholds, `funding_pnl`, `coinbase_fees`) in their JSON output

## [0.6.0] - 2026-SEP-09

### Added

- Bumped `prime-sdk-go` to v0.12.0
- New financing commands: `get-conversion-fees`, `get-cross-margin-liquidation`, `list-cross-margin-liquidations`, `list-trade-finance-obligations`, `get-entity-rewards-rate`, `get-portfolio-rewards-rate`
- New futures commands: `get-equity`, `get-derivatives-currency-summary`, `get-derivative-positions`
- MCP tools for the new endpoints, plus backfill for `get_market_data`, `get_cross_margin_prime_overview`, `get_cross_margin_risk_parameters`, `update_funding_settings`
- Optional flags: `orders edit --offset/--wig-level`, `staking portfolio-unstake --validator-provider`, `staking stake/unstake --metadata-external-id`, `transactions create-withdrawal --counterparty`
- API error output now includes `trace_id` and the spec subcode/code description by default

## [0.5.1] - 2026-JUL-24

### Added

- Bumped `prime-sdk-go` to v0.9.1 (module path: `github.com/coinbase/prime-sdk-go`)

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
