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
)

var listActivitiesCmd = &cobra.Command{
	Use:   "list-activities",
	Short: "Lists activities meeting filter criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

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

		request := &prime.ListActivitiesRequest{
			PortfolioId: client.Credentials.PortfolioId,
			Symbols:     symbols,
			Categories:  categories,
			Statuses:    statuses,
			Start:       start,
			End:         end,
			Pagination:  pagination,
		}

		response, err := client.ListActivities(ctx, request)
		if err != nil {
			return fmt.Errorf("error listing activities: %v", err)
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
	rootCmd.AddCommand(listActivitiesCmd)

	listActivitiesCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listActivitiesCmd.Flags().StringP(utils.LimitFlag, "l", "", "Pagination limit")
	listActivitiesCmd.Flags().StringP(utils.SortDirectionFlag, "d", "", "Sort direction")
	listActivitiesCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")
	listActivitiesCmd.Flags().StringSliceP(utils.CategoriesFlag, "t", []string{}, "List of categories")
	listActivitiesCmd.Flags().StringSliceP(utils.StatusesFlag, "u", []string{}, "List of statuses")
	listActivitiesCmd.Flags().StringP(utils.StartFlag, "r", "", "Start time in RFC3339 format")
	listActivitiesCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")

	rootCmd.AddCommand(listActivitiesCmd)

}