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

package assets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/assets"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/spf13/cobra"
)

var listAssetsCmd = &cobra.Command{
	Use:   "list",
	Short: "List assets for associated entity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		svc := assets.NewAssetsService(client)

		vals, err := listAssets(svc, entityId)
		if err != nil {
			return err
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {

				limit := paginationParams.Limit

				sortItemsByName(vals, paginationParams.SortDirection)

				var toPrint []*model.Asset
				next := ""

				if int32(len(vals)) >= limit {

					toPrint = vals[:limit]
					vals = vals[limit:]
					next = "true"

				} else {
					toPrint = vals
					vals = nil
				}

				if err := utils.PrintJsonDocs(cmd, toPrint); err != nil {
					return nil, err
				}

				return &model.Pagination{NextCursor: next}, nil
			},
		)
	},
}

func listAssets(svc assets.AssetsService, entityId string) ([]*model.Asset, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &assets.ListAssetsRequest{
		EntityId: entityId,
	}

	response, err := svc.ListAssets(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list assets: %w", err)
	}

	return response.Assets, nil
}

func sortItemsByName(items []*model.Asset, direction string) {
	sort.Slice(items, func(i, j int) bool {
		if strings.ToLower(direction) == "asc" {
			return items[i].Name < items[j].Name
		}
		return items[i].Name > items[j].Name
	})
}

func init() {

	utils.AddPaginationFlags(listAssetsCmd, true)
	utils.AddEntityIdFlag(listAssetsCmd)
	Cmd.AddCommand(listAssetsCmd)
}
