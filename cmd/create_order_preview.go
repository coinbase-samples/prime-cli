/**
 * Copyright 2023-present Coinbase Global, Inc.
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

package cmd

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/orders"
	"github.com/spf13/cobra"
)

var createOrderPreviewCmd = &cobra.Command{
	Use:   "create-order-preview",
	Short: "Preview an order before submitting.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		ordersService := orders.NewOrdersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		order := &model.Order{
			PortfolioId:  portfolioId,
			Side:         utils.GetFlagStringValue(cmd, utils.SideFlag),
			Type:         utils.GetFlagStringValue(cmd, utils.TypeFlag),
			ProductId:    utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			BaseQuantity: utils.GetFlagStringValue(cmd, utils.BaseQuantityFlag),
			QuoteValue:   utils.GetFlagStringValue(cmd, utils.QuoteValueFlag),
			LimitPrice:   utils.GetFlagStringValue(cmd, utils.LimitPriceFlag),
			StartTime:    utils.GetFlagStringValue(cmd, utils.StartTimeFlag),
			ExpiryTime:   utils.GetFlagStringValue(cmd, utils.ExpiryTimeFlag),
			TimeInForce:  utils.GetFlagStringValue(cmd, utils.TimeInForceFlag),
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &orders.CreateOrderRequest{
			Order: order,
		}

		response, err := ordersService.CreateOrder(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create order preview: %w", err)

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
	rootCmd.AddCommand(createOrderPreviewCmd)

	createOrderPreviewCmd.Flags().StringP(utils.SideFlag, "s", "", "Order side (Required)")
	createOrderPreviewCmd.Flags().StringP(utils.ProductIdFlag, "i", "", "ID of the product (Required)")
	createOrderPreviewCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of the order (Required)")
	createOrderPreviewCmd.Flags().StringP(utils.BaseQuantityFlag, "b", "", "Order size in base asset units")
	createOrderPreviewCmd.Flags().StringP(utils.QuoteValueFlag, "q", "", "Order size in quote asset units")
	createOrderPreviewCmd.Flags().StringP(utils.TimeInForceFlag, "f", "", "Determine order fill strategy")
	createOrderPreviewCmd.Flags().StringP(utils.LimitPriceFlag, "l", "", "Limit price for the order")
	createOrderPreviewCmd.Flags().StringP(utils.StartTimeFlag, "", "", "Start time of the order in UTC (TWAP only)")
	createOrderPreviewCmd.Flags().StringP(utils.ExpiryTimeFlag, "", "", "Expiry time of the order in UTC (TWAP and limit GTDT only)")
	createOrderPreviewCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	createOrderPreviewCmd.MarkFlagRequired(utils.SideFlag)
	createOrderPreviewCmd.MarkFlagRequired(utils.ProductIdFlag)
	createOrderPreviewCmd.MarkFlagRequired(utils.TypeFlag)

	createOrderPreviewCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateSide(cmd); err != nil {
			return err
		}
		if err := utils.ValidateOrderTypeAndLimitPrice(cmd); err != nil {
			return err
		}
		if err := utils.ValidateTimeInForce(cmd); err != nil {
			return err
		}
		if err := utils.ValidateQuantities(cmd); err != nil {
			return err
		}
		return nil
	}
}
