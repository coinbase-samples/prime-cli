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

var cancelOrderCmd = &cobra.Command{
	Use:   "cancel-order",
	Short: "Attempt to cancel an open order.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("cannot get client from environment: %w", err)
		}

		orderId, err := cmd.Flags().GetString(utils.OrderIdFlag)
		if err != nil {
			return fmt.Errorf("cannot cancel order: %w", err)
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CancelOrderRequest{
			PortfolioId: portfolioId,
			OrderId:     orderId,
		}

		response, err := client.CancelOrder(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot cancel order: %w", err)
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
	rootCmd.AddCommand(cancelOrderCmd)

	cancelOrderCmd.Flags().StringP(utils.OrderIdFlag, "i", "", "ID of the order to cancel (Required)")
	err := cancelOrderCmd.MarkFlagRequired(utils.OrderIdFlag)
	if err != nil {
		return
	}

	cancelOrderCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateUUIDFlag(cmd, utils.OrderIdFlag); err != nil {
			return err
		}
		return nil
	}

	cancelOrderCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
	cancelOrderCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
}
