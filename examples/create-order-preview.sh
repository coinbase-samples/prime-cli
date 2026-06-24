#!/usr/bin/env bash
# Preview an order without submitting it to the exchange.
set -euo pipefail

./primectl orders create-preview \
  --portfolio-id "${PORTFOLIO_ID:?set PORTFOLIO_ID}" \
  --product-id ETH-USD \
  --side BUY \
  --type LIMIT \
  --base-quantity 0.01 \
  --limit-price 2000 \
  --time-in-force GOOD_UNTIL_CANCELLED
