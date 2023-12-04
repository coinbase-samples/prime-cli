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

var listEntityUsersCmd = &cobra.Command{
	Use:   "list-entity-users",
	Short: "List users for associated entity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		request := &prime.ListEntityUsersRequest{
			EntityId:   client.Credentials.EntityId,
			Pagination: pagination,
		}

		response, err := client.ListEntityUsers(ctx, request)
		if err != nil {
			return fmt.Errorf("error listing users: %v", err)
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
	rootCmd.AddCommand(listEntityUsersCmd)

	listEntityUsersCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listEntityUsersCmd.Flags().StringP(utils.LimitFlag, "l", "", "Pagination limit")
	listEntityUsersCmd.Flags().StringP(utils.SortDirectionFlag, "d", "", "Sort direction")

}
