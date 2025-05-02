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

		ordersService := orders.NewOrdersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &orders.ListOpenOrdersRequest{
			PortfolioId: portfolioId,
			ProductId:   utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
		}

		response, err := ordersService.ListOpenOrders(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list open orders: %w", err)
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
	Cmd.AddCommand(listOpenOrdersCmd)

	listOpenOrdersCmd.Flags().String(utils.ProductIdFlag, "", "ID of the product")
	utils.AddPortfolioIdFlag(listOpenOrdersCmd)

	listOpenOrdersCmd.MarkFlagRequired(utils.ProductIdFlag)
}
