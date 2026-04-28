/**
 * Copyright 2026-present Coinbase Global, Inc.
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

package positions

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/positions"

	"github.com/spf13/cobra"
)

var listAggregatePositionsCmd = &cobra.Command{
	Use:   "list-aggregate",
	Short: "Lists aggregate positions for an entity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := positions.NewPositionsService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listAggregatePositions(svc, entityId, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Positions); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listAggregatePositions(
	svc positions.PositionsService,
	entityId string,
	pagination *model.PaginationParams,
) (*positions.ListAggregateEntityPositionsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &positions.ListAggregateEntityPositionsRequest{
		EntityId:   entityId,
		Pagination: pagination,
	}

	response, err := svc.ListAggregateEntityPositions(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list aggregate positions: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listAggregatePositionsCmd)

	utils.AddEntityIdFlag(listAggregatePositionsCmd)
	utils.AddPaginationFlags(listAggregatePositionsCmd, true)
}
