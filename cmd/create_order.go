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
	"encoding/json"
	"fmt"
	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var createOrderCmd = &cobra.Command{
	Use:   "create-order",
	Short: "Submit an order.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		clientOrderId := uuid.New().String()

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		order := &prime.Order{
			PortfolioId:   client.Credentials.PortfolioId,
			Side:          utils.GetFlagStringValue(cmd, utils.SideFlag),
			Type:          utils.GetFlagStringValue(cmd, utils.TypeFlag),
			ClientOrderId: clientOrderId,
			ProductId:     utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			BaseQuantity:  utils.GetFlagStringValue(cmd, utils.BaseQuantityFlag),
			QuoteValue:    utils.GetFlagStringValue(cmd, utils.QuoteValueFlag),
			LimitPrice:    utils.GetFlagStringValue(cmd, utils.LimitPriceFlag),
			StartTime:     utils.GetFlagStringValue(cmd, utils.StartTimeFlag),
			ExpiryTime:    utils.GetFlagStringValue(cmd, utils.ExpiryTimeFlag),
			TimeInForce:   utils.GetFlagStringValue(cmd, utils.TimeInForceFlag),
		}

		request := &prime.CreateOrderRequest{
			Order: order,
		}

		response, err := client.CreateOrder(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}

		jsonResponse, err := json.MarshalIndent(response, "", utils.JsonIndent)
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}
		fmt.Println(string(jsonResponse))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createOrderCmd)

	createOrderCmd.Flags().StringP(utils.SideFlag, "s", "", "Order side (Required)")
	createOrderCmd.Flags().StringP(utils.ProductIdFlag, "i", "", "ID of the product (Required)")
	createOrderCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of the order (Required)")

	createOrderCmd.Flags().StringP(utils.BaseQuantityFlag, "b", "", "Order size in base asset units")
	createOrderCmd.Flags().StringP(utils.QuoteValueFlag, "q", "", "Order size in quote asset units")

	createOrderCmd.Flags().StringP(utils.TimeInForceFlag, "f", "", "Determine order fill strategy")

	createOrderCmd.Flags().StringP(utils.LimitPriceFlag, "l", "", "Limit price for the order")
	createOrderCmd.Flags().StringP(utils.StartTimeFlag, "", "", "The start time of the order in UTC (TWAP only)")
	createOrderCmd.Flags().StringP(utils.ExpiryTimeFlag, "", "", "The expiry time of the order in UTC (TWAP and limit GTD only)")

	createOrderCmd.MarkFlagRequired(utils.SideFlag)
	createOrderCmd.MarkFlagRequired(utils.ProductIdFlag)
	createOrderCmd.MarkFlagRequired(utils.TypeFlag)

	createOrderCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
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
