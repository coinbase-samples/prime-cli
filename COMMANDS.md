# Prime CLI Commands

A copy/paste-friendly reference for every `primectl` command in v0.4.0. Each command is shown as a runnable bash snippet that uses environment variables for the IDs you'll most often substitute.

> Anything marked `<...>` is a placeholder you should replace before running.

## Setup

Before running any of the commands below, export your Prime API credentials and the IDs you'll use most often. The CLI reads `PRIME_CREDENTIALS` from the environment, and most ID flags fall back to the `portfolioId` / `entityId` in those credentials when the corresponding flag is omitted — but exporting `PORTFOLIO_ID`, `ENTITY_ID`, and `WALLET_ID` lets you reuse the same shell variables across commands.

```bash
export PRIME_CREDENTIALS='{
  "accessKey":"ACCESSKEY_HERE",
  "passphrase":"PASSPHRASE_HERE",
  "signingKey":"SIGNINGKEY_HERE",
  "portfolioId":"PORTFOLIOID_HERE",
  "svcAccountId":"SVCACCOUNTID_HERE",
  "entityId":"ENTITYID_HERE"
}'

export PORTFOLIO_ID="<your-portfolio-uuid>"
export ENTITY_ID="<your-entity-uuid>"
export WALLET_ID="<a-wallet-uuid>"

# Optional: extend default 7s request timeout
export primeCliTimeout=30

# Build once
go build -o primectl
```

Tip: append `--format` to any command for pretty-printed JSON output.

---

## Top-level

```bash
./primectl --help
./primectl version
```

## activities

```bash
./primectl activities list --portfolio-id "$PORTFOLIO_ID"
./primectl activities list --portfolio-id "$PORTFOLIO_ID" --symbols ETH,BTC --categories TRADE --statuses COMPLETED --start 2026-01-01T00:00:00Z --end 2026-04-28T00:00:00Z --all
./primectl activities list-entity --entity-id "$ENTITY_ID" --all
./primectl activities get --portfolio-id "$PORTFOLIO_ID" --id <activity-id>
./primectl activities get-entity --portfolio-id "$PORTFOLIO_ID" --id <activity-id>
```

## address-book

```bash
./primectl address-book list --portfolio-id "$PORTFOLIO_ID"
./primectl address-book list --portfolio-id "$PORTFOLIO_ID" --symbol ETH --search 0x
./primectl address-book create \
  --portfolio-id "$PORTFOLIO_ID" \
  --address 0xabc123... \
  --symbol ETH \
  --name "My ETH wallet" \
  --account-identifier "memo-or-tag"
```

## advanced-transfers

```bash
./primectl advanced-transfers list --portfolio-id "$PORTFOLIO_ID" --all
./primectl advanced-transfers list --portfolio-id "$PORTFOLIO_ID" --states ADVANCED_TRANSFER_STATE_PENDING --transfer-type ADVANCED_TRANSFER_TYPE_BLIND_MATCH

./primectl advanced-transfers create \
  --portfolio-id "$PORTFOLIO_ID" \
  --transfer-type ADVANCED_TRANSFER_TYPE_BLIND_MATCH \
  --amount 1.0 \
  --currency ETH \
  --source-type WALLET \
  --source-value "$WALLET_ID" \
  --target-type COUNTERPARTY_ID \
  --target-value <counterparty-id> \
  --reference-id <ref-id> \
  --trade-date 2026-04-28 \
  --settlement-date 2026-04-30 \
  --settlement-time 16:00:00Z

./primectl advanced-transfers cancel \
  --portfolio-id "$PORTFOLIO_ID" \
  --advanced-transfer-id <advanced-transfer-id>

./primectl advanced-transfers list-transactions \
  --portfolio-id "$PORTFOLIO_ID" \
  --advanced-transfer-id <advanced-transfer-id>
```

## allocations

```bash
./primectl allocations list --portfolio-id "$PORTFOLIO_ID" --start 2026-01-01T00:00:00Z --all
./primectl allocations get --portfolio-id "$PORTFOLIO_ID" --allocation-id <allocation-id>
./primectl allocations get-net --portfolio-id "$PORTFOLIO_ID" --allocation-id <netting-id>

./primectl allocations create \
  --allocation-id <client-allocation-id> \
  --source-portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --order-ids <order-uuid-1> --order-ids <order-uuid-2> \
  --allocation-legs '[{"allocation_leg_id":"leg-1","destination_portfolio_id":"<dest-portfolio>","amount":"1.0"}]' \
  --size-type BASE \
  --remainder-destination-portfolio-id <remainder-portfolio>

./primectl allocations create-net \
  --source-portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --order-ids <order-uuid-1> \
  --allocation-legs '[{"allocation_leg_id":"leg-1","destination_portfolio_id":"<dest>","amount":"1.0"}]' \
  --size-type BASE \
  --remainder-destination-portfolio-id <remainder-portfolio>
```

## assets

```bash
./primectl assets list --entity-id "$ENTITY_ID"
```

## balances

```bash
./primectl balances list --portfolio-id "$PORTFOLIO_ID"
./primectl balances list --portfolio-id "$PORTFOLIO_ID" --symbols ETH --type TRADING_BALANCES
./primectl balances list-entity --entity-id "$ENTITY_ID" --all
./primectl balances get-wallet --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"
./primectl balances list-onchain --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"
```

## commission

```bash
./primectl commission get --portfolio-id "$PORTFOLIO_ID"
```

## financing

```bash
./primectl financing get-buying-power --portfolio-id "$PORTFOLIO_ID" --base-currency ETH --quote-currency USD
./primectl financing get-credit-info --portfolio-id "$PORTFOLIO_ID"
./primectl financing get-cross-margin-overview --entity-id "$ENTITY_ID"
./primectl financing get-entity-locate-availabilities --entity-id "$ENTITY_ID"
./primectl financing get-margin-info --entity-id "$ENTITY_ID"
./primectl financing get-pricing-fees --entity-id "$ENTITY_ID"
./primectl financing get-withdrawal-power --portfolio-id "$PORTFOLIO_ID" --symbol USD

./primectl financing list-financing-eligible-assets
./primectl financing list-locates --portfolio-id "$PORTFOLIO_ID" --date 2026-04-28
./primectl financing list-interest-accruals --entity-id "$ENTITY_ID" --portfolio-id "$PORTFOLIO_ID" --start-date 2026-04-01 --end-date 2026-04-28
./primectl financing list-portfolio-interest-accruals --portfolio-id "$PORTFOLIO_ID" --start-date 2026-04-01 --end-date 2026-04-28
./primectl financing list-margin-call-summaries --entity-id "$ENTITY_ID" --start-date 2026-04-01 --end-date 2026-04-28
./primectl financing list-margin-conversions --entity-id "$ENTITY_ID" --portfolio-id "$PORTFOLIO_ID" --start-date 2026-04-01 --end-date 2026-04-28

./primectl financing create-locate \
  --portfolio-id "$PORTFOLIO_ID" \
  --symbol ETH \
  --amount 10 \
  --date 2026-04-29
```

## futures

All futures commands accept `--entity-id`. If omitted, the value falls back to the `entityId` in `PRIME_CREDENTIALS`.

```bash
./primectl futures get-balance              --entity-id "$ENTITY_ID"
./primectl futures get-positions            --entity-id "$ENTITY_ID"
./primectl futures get-margin-call-details  --entity-id "$ENTITY_ID"
./primectl futures get-risk-limits          --entity-id "$ENTITY_ID"
./primectl futures get-settings             --entity-id "$ENTITY_ID"
./primectl futures list-sweeps              --entity-id "$ENTITY_ID"

./primectl futures set-auto-sweep   --entity-id "$ENTITY_ID" --auto-sweep-enabled=true
./primectl futures schedule-sweep   --entity-id "$ENTITY_ID" --amount 1000 --currency USD
./primectl futures cancel-sweep     --entity-id "$ENTITY_ID"
./primectl futures set-settings     --entity-id "$ENTITY_ID" --target-derivatives-excess 0.10
```

## invoices

```bash
./primectl invoices list --all
./primectl invoices list --states INVOICE_STATE_PAID --billing-year 2026 --billing-month 4
```

## onchain-address-book

```bash
./primectl onchain-address-book list-groups --portfolio-id "$PORTFOLIO_ID"

./primectl onchain-address-book create-group-entry \
  --portfolio-id "$PORTFOLIO_ID" \
  --id <address-group-id> \
  --address 0xabc123... \
  --network-type ethereum-mainnet \
  --name "Counterparty A"

./primectl onchain-address-book update-group-entry \
  --portfolio-id "$PORTFOLIO_ID" \
  --id <address-group-id> \
  --address 0xdef456... \
  --network-type ethereum-mainnet \
  --name "Counterparty A (renamed)"

./primectl onchain-address-book delete-group-entry \
  --portfolio-id "$PORTFOLIO_ID" \
  --id <address-group-id>
```

## orders

```bash
./primectl orders list                 --portfolio-id "$PORTFOLIO_ID" --side BUY --start 2026-04-01T00:00:00Z --all
./primectl orders list-open            --portfolio-id "$PORTFOLIO_ID"
./primectl orders list-fills           --portfolio-id "$PORTFOLIO_ID" --order-id <order-id>
./primectl orders list-portfolio-fills --portfolio-id "$PORTFOLIO_ID" --start 2026-04-01T00:00:00Z

./primectl orders get          --portfolio-id "$PORTFOLIO_ID" --order-id <order-id>
./primectl orders cancel       --portfolio-id "$PORTFOLIO_ID" --order-id <order-id>
./primectl orders edit-history --portfolio-id "$PORTFOLIO_ID" --order-id <order-id>

./primectl orders create \
  --portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --side BUY \
  --type MARKET \
  --base-quantity 0.001

./primectl orders create-preview \
  --portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --side BUY \
  --type LIMIT \
  --base-quantity 0.01 \
  --limit-price 2000 \
  --time-in-force GOOD_UNTIL_CANCELLED

./primectl orders edit \
  --portfolio-id "$PORTFOLIO_ID" \
  --order-id <order-id> \
  --new-base-quantity 0.02 \
  --new-limit-price 2050

./primectl orders create-quote \
  --portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --side BUY \
  --base-quantity 0.5 \
  --limit-price 2000 \
  --settle-currency USD

./primectl orders accept-quote \
  --portfolio-id "$PORTFOLIO_ID" \
  --quote-id <quote-id>
```

## payment-methods

```bash
./primectl payment-methods list --entity-id "$ENTITY_ID"
./primectl payment-methods get  --payment-method-id <payment-method-id>
```

## portfolios

```bash
./primectl portfolios list
./primectl portfolios get              --portfolio-id "$PORTFOLIO_ID"
./primectl portfolios get-credit       --portfolio-id "$PORTFOLIO_ID"
./primectl portfolios get-counterparty --portfolio-id "$PORTFOLIO_ID"
```

## positions

```bash
./primectl positions list           --entity-id "$ENTITY_ID" --all
./primectl positions list-aggregate --entity-id "$ENTITY_ID" --all
```

## products

```bash
./primectl products list --portfolio-id "$PORTFOLIO_ID" --all
./primectl products get-candles \
  --portfolio-id "$PORTFOLIO_ID" \
  --product-id ETH-USD \
  --start 2026-04-27T00:00:00Z \
  --end   2026-04-28T00:00:00Z \
  --granularity ONE_HOUR
```

## staking

```bash
./primectl staking stake                    --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"
./primectl staking unstake                  --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID" --amount 1.0
./primectl staking claim-rewards            --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"
./primectl staking get-status               --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"
./primectl staking preview-unstake          --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID" --amount 1.0
./primectl staking get-unstaking-status     --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"

./primectl staking portfolio-stake-initiate --portfolio-id "$PORTFOLIO_ID" --symbol ETH --amount 1.0
./primectl staking portfolio-unstake        --portfolio-id "$PORTFOLIO_ID" --symbol ETH --amount 1.0

./primectl staking query-validators \
  --portfolio-id "$PORTFOLIO_ID" \
  --transaction-ids <txn-id-1>,<txn-id-2>
```

## transactions

```bash
./primectl transactions list                 --portfolio-id "$PORTFOLIO_ID" --all
./primectl transactions list-wallet          --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID" --all
./primectl transactions get                  --transaction-id <transaction-id>

./primectl transactions create-transfer \
  --portfolio-id "$PORTFOLIO_ID" \
  --source-wallet-id "$WALLET_ID" \
  --destination-wallet-id <dest-wallet-id> \
  --symbol ETH \
  --amount 1.0

./primectl transactions create-withdrawal \
  --portfolio-id "$PORTFOLIO_ID" \
  --source-wallet-id "$WALLET_ID" \
  --symbol ETH \
  --amount 1.0 \
  --destination-type DESTINATION_BLOCKCHAIN \
  --blockchain-address 0xabc123...

./primectl transactions create-conversion \
  --portfolio-id "$PORTFOLIO_ID" \
  --source-wallet-id <usd-wallet-id> \
  --source-symbol USD \
  --destination-wallet-id <usdc-wallet-id> \
  --destination-symbol USDC \
  --amount 100

./primectl transactions create-onchain \
  --portfolio-id "$PORTFOLIO_ID" \
  --wallet-id "$WALLET_ID" \
  --raw-unsigned-transaction 0x... \
  --chain-id 1 \
  --url https://mainnet.infura.io/v3/...

./primectl transactions submit-deposit-travel-rule-data \
  --portfolio-id "$PORTFOLIO_ID" \
  --transaction-id <transaction-id> \
  --is-self=false \
  --originator '{"name":"Alice","address":"123 Main St"}' \
  --beneficiary '{"name":"Bob","address":"456 Oak Ave"}'

./primectl transactions get-travel-rule-data \
  --portfolio-id "$PORTFOLIO_ID" \
  --transaction-id <transaction-id>
```

## users

```bash
./primectl users list        --portfolio-id "$PORTFOLIO_ID" --all
./primectl users list-entity --entity-id "$ENTITY_ID" --all
```

## wallets

```bash
./primectl wallets list --portfolio-id "$PORTFOLIO_ID" --all
./primectl wallets list --portfolio-id "$PORTFOLIO_ID" --type TRADING --symbols ETH,USD

./primectl wallets get --portfolio-id "$PORTFOLIO_ID" --wallet-id "$WALLET_ID"

./primectl wallets create \
  --portfolio-id "$PORTFOLIO_ID" \
  --name "My new ETH wallet" \
  --type VAULT \
  --symbol ETH

./primectl wallets create \
  --portfolio-id "$PORTFOLIO_ID" \
  --name "My onchain wallet" \
  --type ONCHAIN \
  --network-family NETWORK_FAMILY_EVM \
  --network-id base \
  --network-type mainnet

./primectl wallets create-deposit-address \
  --portfolio-id "$PORTFOLIO_ID" \
  --wallet-id "$WALLET_ID" \
  --network-id ethereum-mainnet

./primectl wallets get-deposit-instructions \
  --portfolio-id "$PORTFOLIO_ID" \
  --wallet-id "$WALLET_ID" \
  --deposit-type CRYPTO

./primectl wallets list-addresses \
  --portfolio-id "$PORTFOLIO_ID" \
  --wallet-id "$WALLET_ID" \
  --all
```

---

## Common flags reference

- `--format` — pretty-print JSON output (root-level flag, works on every command).
- `--portfolio-id` — overrides the `portfolioId` from `PRIME_CREDENTIALS`.
- `--entity-id` — overrides the `entityId` from `PRIME_CREDENTIALS`.
- `--idempotency-key` — supply your own UUID for retry-safe writes; auto-generated if blank.
- `--client-order-id` / `--client-quote-id` — your own client-side ID; auto-generated if blank.
- Pagination (lists that return `Pagination`): `--limit`, `--sort-direction`, `--all` (drain all pages), `--interactive` (page on key-press).
- Time ranges: `--start` / `--end` use RFC3339 (e.g. `2026-04-28T00:00:00Z`). The financing date filters use `--start-date` / `--end-date`.

## Tips

- `--all` is the easiest way to drain a paginated list into a single JSON-per-line stream you can pipe into `jq`.
- Because most ID-bearing flags read from `PRIME_CREDENTIALS` when omitted, you can drop `--portfolio-id` / `--entity-id` if your credentials JSON already has the right values for the call.
- For shell pipelines, prefer `./primectl <cmd> | jq .` over `--format` so you keep one document per line for streaming/`jq -c` use.
