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

package portfolios

import (
	"fmt"
	"sort"
	"strings"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/portfolios"
	"github.com/spf13/cobra"
)

var listPortfoliosCmd = &cobra.Command{
	Use:   "list",
	Short: "List portfolios associated with API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("cannot get client from environment: %w", err)
		}

		vals, err := listPortfolios(portfolios.NewPortfoliosService(client))
		if err != nil {
			return err
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {

				limit := paginationParams.Limit

				sortItemsByName(vals, paginationParams.SortDirection)

				var toPrint []*model.Portfolio
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

func listPortfolios(svc portfolios.PortfoliosService) ([]*model.Portfolio, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListPortfolios(ctx, &portfolios.ListPortfoliosRequest{})
	if err != nil {
		return nil, fmt.Errorf("cannot list portfolios: %w", err)
	}

	return response.Portfolios, nil
}

func sortItemsByName(items []*model.Portfolio, direction string) {
	sort.Slice(items, func(i, j int) bool {
		if strings.ToLower(direction) == "asc" {
			return items[i].Name < items[j].Name
		}
		return items[i].Name > items[j].Name
	})
}

func init() {
	utils.AddPaginationFlags(listPortfoliosCmd, true)
	Cmd.AddCommand(listPortfoliosCmd)
}
