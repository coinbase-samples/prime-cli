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

package users

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/users"

	"github.com/spf13/cobra"
)

var listEntityUsersCmd = &cobra.Command{
	Use:   "list-entity",
	Short: "List users for associated entity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := users.NewUsersService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listEntityUsers(svc, entityId, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Users); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listEntityUsers(
	svc users.UsersService,
	entityId string,
	pagination *model.PaginationParams,
) (*users.ListEntityUsersResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &users.ListEntityUsersRequest{
		EntityId:   entityId,
		Pagination: pagination,
	}

	response, err := svc.ListEntityUsers(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list users: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listEntityUsersCmd)

	utils.AddEntityIdFlag(listEntityUsersCmd)
	utils.AddPaginationFlags(listEntityUsersCmd, true)
}
