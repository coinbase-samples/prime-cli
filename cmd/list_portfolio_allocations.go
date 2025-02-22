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
	"github.com/coinbase-samples/prime-sdk-go/allocations"
	"github.com/spf13/cobra"
)

var listPortfolioAllocationsCmd = &cobra.Command{
	Use:   "list-portfolio-allocations",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		allocationsService := allocations.NewAllocationsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		productIds, err := cmd.Flags().GetStringSlice(utils.ProductIdsFlag)
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

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &allocations.ListPortfolioAllocationsRequest{
			PortfolioId: portfolioId,
			ProductIds:  productIds,
			Side:        utils.GetFlagStringValue(cmd, utils.OrderSideFlag),
			Start:       start,
			End:         end,
			Pagination:  pagination,
		}

		response, err := allocationsService.ListPortfolioAllocations(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list allocations: %w", err)
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
	rootCmd.AddCommand(listPortfolioAllocationsCmd)

	listPortfolioAllocationsCmd.Flags().StringSliceP(utils.ProductIdsFlag, "p", []string{}, "List of product IDs")
	listPortfolioAllocationsCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listPortfolioAllocationsCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listPortfolioAllocationsCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listPortfolioAllocationsCmd.Flags().StringP(utils.StartFlag, "s", "", "Start time in RFC3339 format (Required)")
	listPortfolioAllocationsCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")
	listPortfolioAllocationsCmd.Flags().StringP(utils.OrderSideFlag, "o", "", "Side of orders")
	listPortfolioAllocationsCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	listPortfolioAllocationsCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	listPortfolioAllocationsCmd.MarkFlagRequired(utils.StartFlag)
}
