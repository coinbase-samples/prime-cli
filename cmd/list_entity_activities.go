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

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/activities"
	"github.com/spf13/cobra"
)

var listEntityActivitiesCmd = &cobra.Command{
	Use:   "list-entity-activities",
	Short: "Lists entityactivities meeting filter criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		activitiesService := activities.NewActivitiesService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		activityLevel, err := cmd.Flags().GetString(utils.ActivityLevelFlag)
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

		request := &activities.ListEntityActivitiesRequest{
			EntityId:      entityId,
			ActivityLevel: activityLevel,
			Symbols:       symbols,
			Categories:    categories,
			Statuses:      statuses,
			StartTime:     start,
			EndTime:       end,
			Pagination:    pagination,
		}

		response, err := activitiesService.ListEntityActivities(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list entity activities: %w", err)
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
	rootCmd.AddCommand(listEntityActivitiesCmd)

	listEntityActivitiesCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listEntityActivitiesCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listEntityActivitiesCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listEntityActivitiesCmd.Flags().StringP(utils.ActivityLevelFlag, "a", "", "Activity level")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.CategoriesFlag, "t", []string{}, "List of categories")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.StatusesFlag, "u", []string{}, "List of statuses")
	listEntityActivitiesCmd.Flags().StringP(utils.StartFlag, "r", "", "Start time in RFC3339 format")
	listEntityActivitiesCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")
	listEntityActivitiesCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
	listEntityActivitiesCmd.Flags().StringP(utils.EntityIdFlag, "", "", "Entity ID. Uses environment variable if blank")
}
