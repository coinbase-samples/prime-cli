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

var listWalletTransactionsCmd = &cobra.Command{
	Use:   "list-wallet-transactions",
	Short: "Lists transaction for a given wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		types, err := cmd.Flags().GetStringSlice(utils.TypesFlag)
		if err != nil {
			return err
		}

		startStr, err := cmd.Flags().GetString(utils.StartFlag)
		if err != nil {
			return err
		}

		endStr, err := cmd.Flags().GetString(utils.EndFlag)
		if err != nil {
			return err
		}

		start, end, err := utils.ParseDateRange(startStr, endStr)
		if err != nil {
			return err
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		request := &prime.ListWalletTransactionsRequest{
			PortfolioId: client.Credentials.PortfolioId,
			WalletId:    utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
			Types:       types,
			Start:       start,
			End:         end,
			Pagination:  pagination,
		}

		response, err := client.ListWalletTransactions(ctx, request)
		if err != nil {
			return fmt.Errorf("error listing transactions: %w", err)
		}

		jsonResponse, err := json.MarshalIndent(response, "", utils.JsonIndent)
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}
		fmt.Println(string(jsonResponse))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listWalletTransactionsCmd)

	listWalletTransactionsCmd.Flags().StringP(utils.WalletIdFlag, "i", "", "ID for given wallet")
	listWalletTransactionsCmd.Flags().StringSliceP(utils.TypesFlag, "t", []string{}, "Types of transactions")
	listWalletTransactionsCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listWalletTransactionsCmd.Flags().StringP(utils.LimitFlag, "l", "", "Pagination limit")
	listWalletTransactionsCmd.Flags().StringP(utils.SortDirectionFlag, "d", "", "Sort direction")
	listWalletTransactionsCmd.Flags().StringP(utils.StartFlag, "s", "", "Start time in RFC3339 format (Required)")
	listWalletTransactionsCmd.Flags().StringP(utils.EndFlag, "e", "", "End time in RFC3339 format")
	listWalletTransactionsCmd.Flags().StringP(utils.SymbolsFlag, "y", "", "Asset symbols")

	listWalletTransactionsCmd.MarkFlagRequired(utils.WalletIdFlag)

}
