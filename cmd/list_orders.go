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
	"github.com/spf13/cobra"
	"strings"
)

var listOrdersCmd = &cobra.Command{
	Use:   "list-orders",
	Short: "Lists orders meeting filter criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		productIds, err := cmd.Flags().GetStringSlice(utils.ProductIdsFlag)
		if err != nil {
			return err
		}

		statuses, err := cmd.Flags().GetStringSlice(utils.OrderStatusesFlag)
		if err != nil {
			return err
		}

		orderType, err := cmd.Flags().GetString(utils.OrderTypeFlag)
		if err != nil {
			return err
		}

		if strings.ToUpper(orderType) == utils.OrderStatusOpen {
			return fmt.Errorf("invalid order type: 'OPEN' cannot be used")
		}

		orderSide, err := cmd.Flags().GetString(utils.OrderSideFlag)
		if err != nil {
			return err
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		startStr, err := cmd.Flags().GetString(utils.StartFlag)
		if err != nil {
			return err
		}

		endStr, err := cmd.Flags().GetString(utils.EndFlag)
		if err != nil {
			return err
		}

		start, end, err := utils.ParseDateRange(startStr, endStr)
		if err != nil {
			return err
		}

		request := &prime.ListOrdersRequest{
			PortfolioId: client.Credentials.PortfolioId,
			Statuses:    statuses,
			ProductIds:  productIds,
			Type:        orderType,
			OtherSide:   orderSide,
			Start:       start,
			End:         end,
			Pagination:  pagination,
		}

		response, err := client.ListOrders(ctx, request)
		if err != nil {
			return fmt.Errorf("error listing orders: %v", err)
		}

		jsonResponse, err := json.MarshalIndent(response, "", utils.JsonIndent)
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %v", err)
		}
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listOrdersCmd)

	listOrdersCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listOrdersCmd.Flags().StringP(utils.LimitFlag, "l", "", "Pagination limit")
	listOrdersCmd.Flags().StringP(utils.SortDirectionFlag, "d", "", "Sort direction")
	listOrdersCmd.Flags().StringSliceP(utils.OrderStatusesFlag, "r", []string{}, "List of statuses")
	listOrdersCmd.Flags().StringSliceP(utils.ProductIdsFlag, "p", []string{}, "List of product IDs")
	listOrdersCmd.Flags().StringP(utils.OrderTypeFlag, "t", "", "Type of orders")
	listOrdersCmd.Flags().StringP(utils.OrderSideFlag, "o", "", "Side of orders")
	listOrdersCmd.Flags().StringP(utils.StartFlag, "s", "", "Start time in RFC3339 format")
	listOrdersCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")

	listOrdersCmd.MarkFlagRequired(utils.StartFlag)
}
