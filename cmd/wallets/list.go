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

package wallets

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/wallets"

	"github.com/spf13/cobra"
)

var listWalletsCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets that meet filter criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := wallets.NewWalletsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		walletType := utils.GetFlagStringValue(cmd, utils.TypeFlag)

		symbols, err := cmd.Flags().GetStringSlice(utils.SymbolsFlag)
		if err != nil {
			return fmt.Errorf("cannot get symbols slice: %w", err)
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listWallets(svc, portfolioId, walletType, symbols, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Wallets); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listWallets(
	svc wallets.WalletsService,
	portfolioId,
	walletType string,
	symbols []string,
	pagination *model.PaginationParams,
) (*wallets.ListWalletsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &wallets.ListWalletsRequest{
		PortfolioId: portfolioId,
		Type:        walletType,
		Symbols:     symbols,
		Pagination:  pagination,
	}

	response, err := svc.ListWallets(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list wallets: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listWalletsCmd)

	listWalletsCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of balance (Required)")
	listWalletsCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")

	utils.AddPortfolioIdFlag(listWalletsCmd)
	utils.AddPaginationFlags(listWalletsCmd, true)

	listWalletsCmd.MarkFlagRequired(utils.TypeFlag)
}
