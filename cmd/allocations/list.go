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

package allocations

import (
	"fmt"
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/allocations"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/spf13/cobra"
)

var listPortfolioAllocationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the portfolio allocations",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := allocations.NewAllocationsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		productIds, err := cmd.Flags().GetStringSlice(utils.ProductIdsFlag)
		if err != nil {
			return err
		}

		start, end, err := utils.GetStartEndFlagsAsTime(cmd)
		if err != nil {
			return err
		}

		side := utils.GetFlagStringValue(cmd, utils.OrderSideFlag)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listPortfolioAllocations(svc, portfolioId, productIds, side, start, end, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Allocations); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)

	},
}

func listPortfolioAllocations(
	svc allocations.AllocationsService,
	portfolioId string,
	productIds []string,
	side string,
	start,
	end time.Time,
	pagination *model.PaginationParams,
) (*allocations.ListPortfolioAllocationsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &allocations.ListPortfolioAllocationsRequest{
		PortfolioId: portfolioId,
		ProductIds:  productIds,
		Side:        side,
		Start:       start,
		End:         end,
		Pagination:  pagination,
	}

	response, err := svc.ListPortfolioAllocations(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list allocations: %w", err)
	}

	return response, nil

}

func init() {
	Cmd.AddCommand(listPortfolioAllocationsCmd)

	listPortfolioAllocationsCmd.Flags().StringSlice(utils.ProductIdsFlag, []string{}, "List of product IDs")
	listPortfolioAllocationsCmd.Flags().String(utils.OrderSideFlag, "", "Side of orders")

	utils.AddPortfolioIdFlag(listPortfolioAllocationsCmd)
	utils.AddPaginationFlags(listPortfolioAllocationsCmd, true)
	utils.AddStartEndFlags(listPortfolioAllocationsCmd)

	listPortfolioAllocationsCmd.MarkFlagRequired(utils.StartFlag)
}
