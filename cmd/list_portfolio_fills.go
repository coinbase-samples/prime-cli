/**
 * Copyright 2025-present Coinbase Global, Inc.
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
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/orders"
	"github.com/spf13/cobra"
)

var listPortfolioFillsCmd = &cobra.Command{
	Use:   "list-portfolio-fills",
	Short: "Get fills from a given portfolio ID",
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

		start, end, err := utils.GetStartEndFlagsAsTime(cmd)
		if err != nil {
			return err
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listPortfolioFills(svc, portfolioId, start, end, paginationParams)
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

func listPortfolioFills(
	svc orders.OrdersService,
	portfolioId string,
	start,
	end time.Time,
	pagination *model.PaginationParams,
) (*orders.ListPortfolioFillsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &orders.ListPortfolioFillsRequest{
		PortfolioId: portfolioId,
		Start:       start,
		End:         end,
		Pagination:  pagination,
	}

	response, err := svc.ListPortfolioFills(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list portfolio fills: %w", err)
	}

	return response, nil
}

func init() {
	rootCmd.AddCommand(listPortfolioFillsCmd)

	utils.AddPortfolioIdFlag(listPortfolioFillsCmd)
	utils.AddPaginationFlags(listPortfolioFillsCmd, true)
	utils.AddStartEndFlags(listPortfolioFillsCmd)

	listPortfolioFillsCmd.MarkFlagRequired(utils.StartFlag)
}
