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

var listActivitiesCmd = &cobra.Command{
	Use:   "list-activities",
	Short: "Lists activities meeting filter criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		symbols, err := cmd.Flags().GetStringSlice(utils.SymbolsFlag)
		if err != nil {
			return err
		}

		categories, err := cmd.Flags().GetStringSlice(utils.CategoriesFlag)
		if err != nil {
			return err
		}

		statuses, err := cmd.Flags().GetStringSlice(utils.StatusesFlag)
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

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.ListActivitiesRequest{
			PortfolioId: portfolioId,
			Symbols:     symbols,
			Categories:  categories,
			Statuses:    statuses,
			Start:       start,
			End:         end,
			Pagination:  pagination,
		}

		response, err := client.ListActivities(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list activities: %w", err)
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
	rootCmd.AddCommand(listActivitiesCmd)

	listActivitiesCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listActivitiesCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listActivitiesCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listActivitiesCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")
	listActivitiesCmd.Flags().StringSliceP(utils.CategoriesFlag, "t", []string{}, "List of categories")
	listActivitiesCmd.Flags().StringSliceP(utils.StatusesFlag, "u", []string{}, "List of statuses")
	listActivitiesCmd.Flags().StringP(utils.StartFlag, "r", "", "Start time in RFC3339 format")
	listActivitiesCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")
	listActivitiesCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	listActivitiesCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	rootCmd.AddCommand(listActivitiesCmd)
}
