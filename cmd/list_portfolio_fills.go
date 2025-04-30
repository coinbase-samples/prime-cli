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

		ordersService := orders.NewOrdersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		startDateStr, err := cmd.Flags().GetString(utils.StartFlag)
		if err != nil {
			return err
		}

		endDateStr, err := cmd.Flags().GetString(utils.EndFlag)
		if err != nil {
			return err
		}

		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return fmt.Errorf("invalid start date format: %w", err)
		}

		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return fmt.Errorf("invalid end date format: %w", err)
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &orders.ListPortfolioFillsRequest{
			PortfolioId: portfolioId,
			Start:       startDate,
			End:         endDate,
			Pagination:  pagination,
		}

		response, err := ordersService.ListPortfolioFills(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list portfolio fills: %w", err)
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
	rootCmd.AddCommand(listPortfolioFillsCmd)

	listPortfolioFillsCmd.Flags().StringP(utils.StartFlag, "s", "", "Start date (Required)")
	listPortfolioFillsCmd.Flags().StringP(utils.EndFlag, "e", "", "End date (Required)")
	listPortfolioFillsCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listPortfolioFillsCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listPortfolioFillsCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listPortfolioFillsCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	listPortfolioFillsCmd.MarkFlagRequired(utils.StartFlag)
}
