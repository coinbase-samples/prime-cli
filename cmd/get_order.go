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
	"github.com/coinbase-samples/prime-sdk-go"

	"github.com/spf13/cobra"
)

var getOrderCmd = &cobra.Command{
	Use:   "get-order",
	Short: "Get an order from an order ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		portfolioId := utils.GetFlagStringValue(cmd, utils.PortfolioIdFlag)
		if portfolioId == "" {
			portfolioId = client.Credentials.PortfolioId
		}

		request := &prime.GetOrderRequest{
			PortfolioId: portfolioId,
			OrderId:     utils.GetFlagStringValue(cmd, utils.OrderIdFlag),
		}

		response, err := client.GetOrder(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get order: %w", err)
		}

		shouldFormat, err := utils.CheckFormatFlag(cmd)
		if err != nil {
			return err
		}

		jsonResponse, err := utils.MarshalJSON(response, shouldFormat)
		if err != nil {
			return fmt.Errorf("cannot marshal response to JSON: %w", err)
		}
		fmt.Println(string(jsonResponse))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOrderCmd)

	getOrderCmd.Flags().StringP(utils.OrderIdFlag, "i", "", "Order ID (Required)")
	getOrderCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getOrderCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getOrderCmd.MarkFlagRequired(utils.OrderIdFlag)
}
