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

package orders

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/orders"

	"github.com/spf13/cobra"
)

var listOpenOrdersCmd = &cobra.Command{
	Use:   "list-open",
	Short: "Lists open orders meeting filter criteria.",
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

		productIds, err := cmd.Flags().GetStringSlice(utils.ProductIdsFlag)
		if err != nil {
			return err
		}

		orderType, err := cmd.Flags().GetString(utils.OrderTypeFlag)
		if err != nil {
			return err
		}

		orderSide, err := cmd.Flags().GetString(utils.OrderSideFlag)
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

		request := &orders.ListOpenOrdersRequest{
			PortfolioId: portfolioId,
			ProductIds:  productIds,
			OrderType:   orderType,
			OrderSide:   orderSide,
			Start:       start,
			End:         end,
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {

				request.Pagination = paginationParams

				response, err := listOpenOrders(svc, request)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Orders); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listOpenOrders(
	svc orders.OrdersService,
	req *orders.ListOpenOrdersRequest,
) (*orders.ListOpenOrdersResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListOpenOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list open orders: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listOpenOrdersCmd)

	utils.AddPortfolioIdFlag(listOpenOrdersCmd)
	utils.AddProductIdsFlag(listOpenOrdersCmd)
	utils.AddPaginationFlags(listOpenOrdersCmd, false)
	utils.AddSortDirectionFlag(listOpenOrdersCmd)
	utils.AddStartEndFlags(listOpenOrdersCmd)

	utils.AddOrderSideFlag(listOpenOrdersCmd)
	utils.AddOrderTypeFlag(listOpenOrdersCmd)

	// listOpenOrdersCmd.MarkFlagRequired(utils.ProductIdsFlag)
}
