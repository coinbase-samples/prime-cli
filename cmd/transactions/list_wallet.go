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
 * limitations under the License.  */

package transactions

import (
	"fmt"
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/transactions"

	"github.com/spf13/cobra"
)

var listWalletTransactionsCmd = &cobra.Command{
	Use:   "list-wallet",
	Short: "Lists transaction for a given wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := transactions.NewTransactionsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		walletId, err := cmd.Flags().GetString(utils.WalletIdFlag)
		if err != nil {
			return err
		}

		types, err := cmd.Flags().GetStringSlice(utils.TypesFlag)
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
				response, err := listWalletTransactions(svc, portfolioId, walletId, types, start, end, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Transactions); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listWalletTransactions(
	svc transactions.TransactionsService,
	portfolioId,
	walletId string,
	types []string,
	start,
	end time.Time,
	pagination *model.PaginationParams,
) (*transactions.ListWalletTransactionsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &transactions.ListWalletTransactionsRequest{
		PortfolioId: portfolioId,
		WalletId:    walletId,
		Types:       types,
		Start:       start,
		End:         end,
		Pagination:  pagination,
	}

	response, err := svc.ListWalletTransactions(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list transactions: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listWalletTransactionsCmd)

	listWalletTransactionsCmd.Flags().StringSlice(utils.TypesFlag, []string{}, "Types of transactions")
	listWalletTransactionsCmd.Flags().String(utils.SymbolsFlag, "", "Asset symbols")

	utils.AddWalletIdFlag(listWalletTransactionsCmd)
	utils.AddPortfolioIdFlag(listWalletTransactionsCmd)
	utils.AddPaginationFlags(listWalletTransactionsCmd, true)
	utils.AddStartEndFlags(listWalletTransactionsCmd)

}
