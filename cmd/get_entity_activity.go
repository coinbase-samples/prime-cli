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

var getEntityActivityCmd = &cobra.Command{
	Use:   "get-entity-activity",
	Short: "Get activity information using Activity ID only",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		activitiesService := activities.NewActivitiesService(client)

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &activities.GetEntityActivityRequest{
			ActivityId: utils.GetFlagStringValue(cmd, utils.ActivityIdFlag),
		}

		response, err := activitiesService.GetEntityActivity(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get entity activity: %w", err)
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
	rootCmd.AddCommand(getEntityActivityCmd)

	getEntityActivityCmd.Flags().StringP(utils.ActivityIdFlag, "i", "", "Activity ID (Required)")
	getEntityActivityCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getEntityActivityCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getEntityActivityCmd.MarkFlagRequired(utils.GenericIdFlag)
}
