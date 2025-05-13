#!/bin/zsh

PRODUCT_ID="ETH-USD"
BASE_QUANTITY="0.005"
LIMIT_PRICE="2700"

# Execute the quote request
QUOTE_ID=$(./../primectl orders create-quote --product-id $PRODUCT_ID --side BUY --base-quantity $BASE_QUANTITY --limit-price $LIMIT_PRICE | jq -r '.quote_id')

# Accept the quote request
ORDER_ID=$(./../primectl orders accept-quote --product-id $PRODUCT_ID --side BUY --quote-id $QUOTE_ID | jq -r '.order_id')

echo "RFQ executed - order id: $ORDER_ID\n"

