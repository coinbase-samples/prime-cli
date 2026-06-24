# Prime CLI

## Overview

The Coinbase Prime command-line interface (CLI) to simplify programmatic interaction with [Coinbase Prime's](https://prime.coinbase.com/) [REST APIs](https://docs.cloud.coinbase.com/prime/reference).

## License

The Prime CLI is free and open source and released under the [Apache License, Version 2.0](LICENSE.txt).

The application and code are only available for demonstration purposes.

## Installation

Install the Prime CLI binary on your computer.

### MacOS

Ensure that you have [Homebrew](https://brew.sh/) installed.

Install the Coinbase Samples Homebrew tap:

```
brew tap coinbase-samples/homebrew-tap
```

Next, install the Prime CLI Homebrew formula:

```
brew install prime-cli
```

### Other Platforms

To install the Prime CLI on other platforms, clone this repository, build, and then add the binary to the [path](https://en.wikipedia.org/wiki/PATH_(variable)).

```
git clone git@github.com:coinbase-samples/prime-cli.git
cd prime-cli
go build -o primectl
```

## Configuration

Once you have the CLI installed, configure your environment.

Set an environment variable in your shell called `PRIME_CREDENTIALS` with your API and portfolio information.

Coinbase Prime API credentials can be created in the Prime web console under Settings -> APIs. Entity ID can be retrieved by calling [Get Portfolio](https://docs.cloud.coinbase.com/prime/reference/primerestapi_getportfolio). If you are not configured yet to call this endpoint, you may proceed without including the `entityId` key and value for the time being, but certain endpoints such as List Invoices and List Assets require it.

`PRIME_CREDENTIALS` should match the following format:
```
export PRIME_CREDENTIALS='{
"accessKey":"ACCESSKEY_HERE",
"passphrase":"PASSPHRASE_HERE",
"signingKey":"SIGNINGKEY_HERE",
"portfolioId":"PORTFOLIOID_HERE",
"svcAccountId":"SVCACCOUNTID_HERE",
"entityId":"ENTITYID_HERE"
}'
```

You may also pass an environment variable called `primeCliTimeout` which will override the default request timeout of 7 seconds. This value should be an integer in seconds.

## Usage

Build the application binary and specify an output name, e.g. `primectl`:

```
go build -o primectl
```

To ensure your project's dependencies are up-to-date, run:
```
go mod tidy
```

To verify that your application is installed correctly and accessible from any location, run the following command. It will include all available requests:

```
./primectl
```

Finally, to run commands for each endpoint, use the following format to test each endpoint. Please note that many endpoints require flags, which are detailed with the `--help` flag.

```
./primectl portfolios list
```

```
./primectl orders create --help
```

```
./primectl orders create-preview -b 0.001 -i ETH-USD -s BUY -t MARKET
```

As of v0.5.0, the CLI covers the full surface area of [prime-sdk-go](https://github.com/coinbase/prime-sdk-go) v0.9.0, including the `advanced-transfers`, `futures`, and `positions` command groups.

## MCP Server

The Prime CLI can run as a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server, exposing Coinbase Prime operations as tools for AI assistants such as Claude Desktop, Cursor, and other MCP-compatible clients.

### Starting the MCP server

```
primectl mcp
```

The server communicates over stdio (stdin/stdout) using the MCP protocol. MCP clients spawn the process and communicate with it directly — no separate port or network configuration is needed.

### Configuring Claude Desktop

Add the following to your Claude Desktop configuration file (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "coinbase-prime": {
      "command": "primectl",
      "args": ["mcp"],
      "env": {
        "PRIME_CREDENTIALS": "{\"accessKey\":\"...\",\"passphrase\":\"...\",\"signingKey\":\"...\",\"portfolioId\":\"...\",\"svcAccountId\":\"...\",\"entityId\":\"...\"}"
      }
    }
  }
}
```

If `PRIME_CREDENTIALS` is already set in your shell environment, you can omit the `env` block entirely.

### Available tools

The MCP server exposes 98 tools across the Coinbase Prime API:

| Tool | Description |
|---|---|
| **Activities** | |
| `list_activities` | List portfolio activities; categories: ORDER, TRANSACTION, ACCOUNT, ALLOCATION, LENDING |
| `get_activity` | Get a specific activity by ID |
| `list_entity_activities` | List activities for an entity across all portfolios |
| `get_entity_activity` | Get a specific entity-level activity by ID |
| **Address Book** | |
| `list_address_book` | List address book entries for a portfolio |
| `create_address_book_entry` | Create a new address book entry |
| **Advanced Transfers** | |
| `list_advanced_transfers` | List advanced transfers for a portfolio |
| `create_advanced_transfer` | Create a new advanced transfer |
| `cancel_advanced_transfer` | Cancel an advanced transfer |
| `list_advanced_transfer_transactions` | List transactions for an advanced transfer |
| **Allocations** | |
| `list_allocations` | List historical allocations for a portfolio |
| `get_allocation` | Get an allocation by ID |
| `create_allocation` | Create a portfolio allocation |
| `get_net_allocation` | Get a net allocation by netting ID |
| `create_net_allocation` | Create a net portfolio allocation |
| **Assets** | |
| `list_assets` | List all supported assets for the entity |
| **Balances** | |
| `list_portfolio_balances` | List asset balances for a portfolio; type: TRADING_BALANCES, VAULT_BALANCES, TOTAL_BALANCES, PRIME_CUSTODY_BALANCES, or UNIFIED_TOTAL_BALANCES |
| `get_wallet_balance` | Get the balance for a specific wallet |
| `list_onchain_balances` | List onchain balances for a wallet |
| `list_entity_balances` | List balances at the entity level; type: TRADING_BALANCES, VAULT_BALANCES, TOTAL_BALANCES, PRIME_CUSTODY_BALANCES, or UNIFIED_TOTAL_BALANCES |
| **Commission** | |
| `get_commission` | Get commission rates and fee tiers |
| **Financing** | |
| `get_buying_power` | Get buying power for a portfolio |
| `get_portfolio_credit_info` | Get post-trade credit information for a portfolio |
| `get_cross_margin_overview` | Get cross-margin overview for an entity |
| `get_entity_locate_availabilities` | Get locate availabilities for an entity |
| `get_margin_information` | Get margin information for an entity |
| `get_pricing_fees` | Get trade finance tiered pricing fees for an entity |
| `get_withdrawal_power` | Get withdrawal power for a portfolio |
| `create_locate` | Create a new locate for a portfolio and asset |
| `list_financing_eligible_assets` | List assets eligible for financing |
| `list_interest_accruals` | List interest accruals for an entity |
| `list_locates` | List existing locates for a portfolio |
| `list_margin_call_summaries` | List margin call summaries for an entity |
| `list_margin_conversions` | List margin conversions for a portfolio (deprecated) |
| `list_portfolio_interest_accruals` | List interest accruals for a portfolio |
| **Futures (FCM)** | |
| `get_fcm_balance` | Get FCM balance summary for an entity |
| `get_fcm_positions` | Get FCM futures positions for an entity |
| `get_fcm_settings` | Get FCM settings for an entity |
| `get_fcm_margin_call_details` | Get FCM margin call details for an entity |
| `get_fcm_risk_limits` | Get FCM risk limits for an entity |
| `set_fcm_settings` | Update FCM settings for an entity |
| `set_fcm_auto_sweep` | Enable or disable FCM auto sweep |
| `list_fcm_sweeps` | List futures sweeps for an entity |
| `schedule_fcm_sweep` | Schedule a futures sweep |
| `cancel_fcm_sweep` | Cancel a scheduled futures sweep |
| **Invoices** | |
| `list_invoices` | List invoices for an entity |
| **Onchain Address Book** | |
| `list_onchain_address_groups` | List onchain address book groups for a portfolio |
| `create_onchain_address_group` | Create an onchain address book group entry |
| `update_onchain_address_group` | Update an onchain address book group entry |
| `delete_onchain_address_group` | Delete an onchain address book group entry |
| **Orders** | |
| `list_orders` | List orders with optional filters; statuses: OPEN, FILLED, CANCELLED, EXPIRED, FAILED, PENDING |
| `list_open_orders` | List currently open orders; filter by product, type, or side |
| `get_order` | Get details of a specific order |
| `list_order_fills` | List fills (executions) for a specific order |
| `get_order_edit_history` | Get the edit history for an order |
| `list_portfolio_fills` | List all fills for a portfolio |
| `preview_order` | Preview an order before submitting; types: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, PEG |
| `create_order` | Submit a new order; types: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, PEG; TIF: GTC, GTD, IOC, FOK |
| `cancel_order` | Attempt to cancel an open order |
| `edit_order` | Edit an existing open order's quantity or price |
| `create_quote` | Create a quote request |
| `accept_quote` | Accept a quote and create an order |
| **Payment Methods** | |
| `list_payment_methods` | List payment methods for an entity |
| `get_payment_method` | Get payment method details by ID |
| **Portfolios** | |
| `list_portfolios` | List all portfolios associated with the API key |
| `get_portfolio` | Get details of a specific portfolio |
| `get_portfolio_credit` | Get post-trade credit information for a portfolio |
| `get_portfolio_counterparty` | Get counterparty ID for a portfolio |
| **Positions** | |
| `list_positions` | List positions for an entity |
| `list_aggregate_positions` | List aggregate positions for an entity |
| **Products** | |
| `list_products` | List all tradeable products (use to discover valid `product_id` values) |
| `get_product_candles` | Get candlestick data for a product; granularities: ONE_MINUTE, FIVE_MINUTES, FIFTEEN_MINUTES, THIRTY_MINUTES, ONE_HOUR, TWO_HOURS, FOUR_HOURS, SIX_HOURS, ONE_DAY |
| **Staking** | |
| `stake` | Create a stake or delegate request for a wallet |
| `unstake` | Create an unstake request for a wallet |
| `get_staking_status` | Get staking status for a wallet |
| `claim_staking_rewards` | Claim staking rewards for a wallet |
| `get_unstaking_status` | Get the status of an unstake operation |
| `preview_unstake` | Preview an unstake operation for a wallet |
| `portfolio_stake_initiate` | Initiate a portfolio-level stake request |
| `portfolio_unstake` | Initiate a portfolio-level unstake request |
| `query_validators` | Query transaction validators for a portfolio |
| **Transactions** | |
| `list_portfolio_transactions` | List portfolio transactions with optional filters |
| `get_transaction` | Get details of a specific transaction |
| `list_wallet_transactions` | List transactions for a specific wallet |
| `create_transfer` | Create an internal transfer between wallets |
| `create_withdrawal` | Create an external withdrawal |
| `create_conversion` | Convert between fiat and stablecoins |
| `create_onchain_transaction` | Create an onchain transaction |
| `get_travel_rule_data` | Get travel rule data for a transaction |
| `submit_deposit_travel_rule_data` | Submit travel rule data for a deposit transaction |
| **Users** | |
| `list_portfolio_users` | List users associated with a portfolio |
| `list_entity_users` | List users for an entity |
| **Wallets** | |
| `list_wallets` | List wallets for a portfolio by type: VAULT, ONCHAIN, TRADING, or QC |
| `get_wallet` | Get details of a specific wallet |
| `create_wallet` | Create a new wallet (VAULT, ONCHAIN, TRADING, or QC) |
| `create_deposit_address` | Create a new deposit address for a wallet |
| `get_wallet_deposit_instructions` | Get deposit instructions for a wallet |
| `list_wallet_addresses` | List addresses for a wallet |

Most tools that require a `portfolio_id` or `entity_id` will fall back to the values in `PRIME_CREDENTIALS` if those arguments are not provided, so many calls work without any parameters.

### Testing with MCP Inspector

You can inspect and test the server interactively using the [MCP Inspector](https://github.com/modelcontextprotocol/inspector):

```
npx @modelcontextprotocol/inspector primectl mcp
```
## Releasing

The Prime CLI is distributed via the [`coinbase-samples/homebrew-tap`](https://github.com/coinbase-samples/homebrew-tap) Homebrew tap. Cutting a new release is a two-repo process: tag the source here, then bump the formula in the tap.

### 1. Bump the version in this repo

1. Update the version string in `cmd/version.go` (`primectlVersion`) to the new semver, e.g. `0.5.0`.
2. Add a new entry at the top of `CHANGELOG.md` following the existing `Added` / `Fixed` format.
3. If user-facing commands changed, refresh `README.md` and `COMMANDS.md`.
4. Open a PR, get it reviewed, and merge to `main`.

### 2. Tag and create a GitHub release

From `main` after the bump is merged:

```
git checkout main && git pull
git tag vX.Y.Z
git push origin vX.Y.Z

gh release create vX.Y.Z \
  --repo coinbase-samples/prime-cli \
  --title "vX.Y.Z" \
  --notes "See CHANGELOG.md for details."
```

Tags must be prefixed with `v` (e.g. `v0.5.0`) — the Homebrew formula's `url` resolves the tag by that exact name.

### 3. Compute the source tarball SHA256

Homebrew pins the GitHub source tarball by hash:

```
VERSION=X.Y.Z
curl -sL "https://github.com/coinbase-samples/prime-cli/archive/refs/tags/v${VERSION}.tar.gz" \
  | shasum -a 256
```

Copy the resulting 64-character hex digest.

### 4. Update the Homebrew formula

In a clone of [`coinbase-samples/homebrew-tap`](https://github.com/coinbase-samples/homebrew-tap), edit `Formula/prime-cli.rb` and update three lines:

```
url "https://github.com/coinbase-samples/prime-cli/archive/refs/tags/vX.Y.Z.tar.gz"
sha256 "<sha256 from the previous step>"
version "X.Y.Z"
```

Leave the `depends_on`, `install`, and `test` blocks untouched — the formula builds from source with `go build`.

### 5. Validate the formula locally

From the tap clone:

```
brew install --build-from-source ./Formula/prime-cli.rb
brew test prime-cli
brew audit --strict --online prime-cli
primectl version   # should print {"version":"X.Y.Z"}
```

`brew audit` will catch a mismatched `sha256`, a broken URL, or formula style issues before review.

### 6. Open the PR against the tap

```
gh pr create \
  --repo coinbase-samples/homebrew-tap \
  --title "prime-cli X.Y.Z" \
  --body "Bumps prime-cli formula to vX.Y.Z. See https://github.com/coinbase-samples/prime-cli/releases/tag/vX.Y.Z"
```

Once merged to `main` on the tap, end users pick up the new version with:

```
brew update
brew upgrade prime-cli
```

