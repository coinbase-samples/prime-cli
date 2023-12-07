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

var getOrderFillsCmd = &cobra.Command{
	Use:   "get-order-fills",
	Short: "Get fills from a given Order ID",
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

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		request := &prime.ListOrderFillsRequest{
			PortfolioId: portfolioId,
			OrderId:     utils.GetFlagStringValue(cmd, utils.OrderIdFlag),
			Pagination:  pagination,
		}

		response, err := client.ListOrderFills(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get order fills: %w", err)
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
	rootCmd.AddCommand(getOrderFillsCmd)

	getOrderFillsCmd.Flags().StringP(utils.OrderIdFlag, "i", "", "Order ID (Required)")
	getOrderFillsCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	getOrderFillsCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	getOrderFillsCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	getOrderFillsCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getOrderFillsCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getOrderFillsCmd.MarkFlagRequired(utils.OrderIdFlag)
}
