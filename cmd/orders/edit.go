/**
 * Copyright 2026-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package orders

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/orders"

	"github.com/spf13/cobra"
)

var editOrderCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edits an existing order",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("cannot get client from environment: %w", err)
		}

		ordersService := orders.NewOrdersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &orders.EditOrderRequest{
			PortfolioId:   portfolioId,
			OrderId:       utils.GetFlagStringValue(cmd, utils.OrderIdFlag),
			ClientOrderId: utils.GetFlagStringValue(cmd, utils.ClientOrderIdFlag),
			BaseQuantity:  utils.GetFlagStringValue(cmd, utils.NewBaseQuantityFlag),
			QuoteValue:    utils.GetFlagStringValue(cmd, utils.NewQuoteValueFlag),
			LimitPrice:    utils.GetFlagStringValue(cmd, utils.NewLimitPriceFlag),
		}

		response, err := ordersService.EditOrder(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot edit order: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)
		return nil
	},
}

func init() {
	Cmd.AddCommand(editOrderCmd)

	editOrderCmd.Flags().String(utils.OrderIdFlag, "", "Order ID (Required)")
	if err := editOrderCmd.MarkFlagRequired(utils.OrderIdFlag); err != nil {
		return
	}

	editOrderCmd.Flags().String(utils.ClientOrderIdFlag, "", "Updated client order ID")
	editOrderCmd.Flags().String(utils.NewBaseQuantityFlag, "", "Updated order size in base asset units")
	editOrderCmd.Flags().String(utils.NewQuoteValueFlag, "", "Updated order size in quote asset units")
	editOrderCmd.Flags().String(utils.NewLimitPriceFlag, "", "Updated limit price")

	utils.AddPortfolioIdFlag(editOrderCmd)
}
