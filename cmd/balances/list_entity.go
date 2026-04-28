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

package balances

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/balances"
	"github.com/coinbase-samples/prime-sdk-go/model"

	"github.com/spf13/cobra"
)

const aggregationTypeFlag = "aggregation-type"

var listEntityBalancesCmd = &cobra.Command{
	Use:   "list-entity",
	Short: "Lists entity-level balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := balances.NewBalancesService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		symbols, err := cmd.Flags().GetStringSlice(utils.SymbolsFlag)
		if err != nil {
			return fmt.Errorf("cannot get symbols: %w", err)
		}

		aggregationType := utils.GetFlagStringValue(cmd, aggregationTypeFlag)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listEntityBalances(svc, entityId, symbols, aggregationType, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Balances); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listEntityBalances(
	svc balances.BalancesService,
	entityId string,
	symbols []string,
	aggregationType string,
	pagination *model.PaginationParams,
) (*balances.ListEntityBalancesResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &balances.ListEntityBalancesRequest{
		EntityId:        entityId,
		Symbols:         symbols,
		AggregationType: model.AggregationType(aggregationType),
		Pagination:      pagination,
	}

	response, err := svc.ListEntityBalances(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list entity balances: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listEntityBalancesCmd)

	listEntityBalancesCmd.Flags().StringSlice(utils.SymbolsFlag, []string{}, "List of symbols")
	listEntityBalancesCmd.Flags().String(aggregationTypeFlag, "", "Balance aggregation type (e.g. TRADING_BALANCES, VAULT_BALANCES, TOTAL_BALANCES)")

	utils.AddEntityIdFlag(listEntityBalancesCmd)
	utils.AddPaginationFlags(listEntityBalancesCmd, true)
}
