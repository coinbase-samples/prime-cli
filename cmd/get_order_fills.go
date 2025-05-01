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

var getOrderFillsCmd = &cobra.Command{
	Use:   "get-order-fills",
	Short: "Get fills from a given Order ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := orders.NewOrdersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		orderId := utils.GetFlagStringValue(cmd, utils.OrderIdFlag)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := getOrderFills(svc, portfolioId, orderId, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Fills); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func getOrderFills(
	svc orders.OrdersService,
	portfolioId,
	orderId string,
	pagination *model.PaginationParams,
) (*orders.ListOrderFillsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &orders.ListOrderFillsRequest{
		PortfolioId: portfolioId,
		OrderId:     orderId,
		Pagination:  pagination,
	}

	response, err := svc.ListOrderFills(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot get order fills: %w", err)
	}

	return response, err
}

func init() {
	rootCmd.AddCommand(getOrderFillsCmd)

	getOrderFillsCmd.Flags().StringP(utils.OrderIdFlag, "i", "", "Order ID (Required)")

	utils.AddPortfolioIdFlag(getOrderFillsCmd)
	utils.AddPaginationFlags(getOrderFillsCmd, true)

	getOrderFillsCmd.MarkFlagRequired(utils.OrderIdFlag)
}
