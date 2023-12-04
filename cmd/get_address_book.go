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

var getAddressBookCmd = &cobra.Command{
	Use:   "get-address-book",
	Short: "Get a list of address book entries",
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

		request := &prime.GetAddressBookRequest{
			PortfolioId: client.Credentials.PortfolioId,
			Symbol:      utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			Search:      utils.GetFlagStringValue(cmd, utils.SearchFlag),
			Pagination:  pagination,
		}
		response, err := client.GetAddressBook(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating portfolio allocations: %v", err)
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
	rootCmd.AddCommand(getAddressBookCmd)

	getAddressBookCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Currency symbol for filtering address book entries")
	getAddressBookCmd.Flags().StringP(utils.SearchFlag, "e", "", "Search term for filtering address book entries")
	getAddressBookCmd.Flags().StringP(utils.CursorFlag, "c", "", "Cursor for pagination")
	getAddressBookCmd.Flags().StringP(utils.LimitFlag, "l", "", "Limit for pagination")
	getAddressBookCmd.Flags().StringP(utils.SortDirectionFlag, "d", "", "Sort direction for pagination")
}
