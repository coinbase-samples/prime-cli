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

package activities

import (
	"fmt"
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/activities"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/spf13/cobra"
)

var listEntityActivitiesCmd = &cobra.Command{
	Use:   "list-entity",
	Short: "Lists entity activities meeting filter criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := activities.NewActivitiesService(client)

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

		start, end, err := utils.GetStartEndFlagsAsTime(cmd)
		if err != nil {
			return err
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listEntityActivities(svc, entityId, activityLevel, symbols, categories, statuses, start, end, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Activities); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listEntityActivities(
	svc activities.ActivitiesService,
	entityId,
	activityLevel string,
	symbols,
	categories,
	statuses []string,
	start,
	end time.Time,
	pagination *model.PaginationParams,
) (*activities.ListEntityActivitiesResponse, error) {
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

	response, err := svc.ListEntityActivities(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list entity activities: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listEntityActivitiesCmd)

	listEntityActivitiesCmd.Flags().StringP(utils.ActivityLevelFlag, "a", "", "Activity level")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.CategoriesFlag, "t", []string{}, "List of categories")
	listEntityActivitiesCmd.Flags().StringSliceP(utils.StatusesFlag, "u", []string{}, "List of statuses")

	utils.AddEntityIdFlag(listEntityActivitiesCmd)
	utils.AddStartEndFlags(listEntityActivitiesCmd)
	utils.AddPaginationFlags(listEntityActivitiesCmd, true)
}
