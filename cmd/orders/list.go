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
	"strings"
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/orders"
	"github.com/spf13/cobra"
)

var listOrdersCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists orders meeting filter criteria.",
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

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listOrders(svc, portfolioId, productIds, statuses, orderType, orderSide, start, end, paginationParams)
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

func listOrders(
	svc orders.OrdersService,
	portfolioId string,
	productIds,
	statuses []string,
	orderType,
	orderSide string,
	start,
	end time.Time,
	pagination *model.PaginationParams,
) (*orders.ListOrdersResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &orders.ListOrdersRequest{
		PortfolioId: portfolioId,
		Statuses:    statuses,
		ProductIds:  productIds,
		Type:        orderType,
		OrderSide:   orderSide,
		Start:       start,
		End:         end,
		Pagination:  pagination,
	}

	response, err := svc.ListOrders(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list orders: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listOrdersCmd)

	listOrdersCmd.Flags().StringSliceP(utils.OrderStatusesFlag, "", []string{}, "List of statuses")
	listOrdersCmd.Flags().StringSliceP(utils.ProductIdsFlag, "", []string{}, "List of product IDs")
	listOrdersCmd.Flags().StringP(utils.OrderTypeFlag, "", "", "Type of orders")
	listOrdersCmd.Flags().StringP(utils.OrderSideFlag, "", "", "Side of orders")

	utils.AddStartEndFlags(listOrdersCmd)
	listOrdersCmd.MarkFlagRequired(utils.StartFlag)

	utils.AddPortfolioIdFlag(listOrdersCmd)
	utils.AddPaginationFlags(listOrdersCmd, true)

}
