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
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/addressbook"

	"github.com/spf13/cobra"
)

var getAddressBookCmd = &cobra.Command{
	Use:   "get-address-book",
	Short: "Get a list of address book entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		addressBookService := addressbook.NewAddressBookService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &addressbook.GetAddressBookRequest{
			PortfolioId: portfolioId,
			Symbol:      utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			Search:      utils.GetFlagStringValue(cmd, utils.SearchFlag),
			Pagination:  pagination,
		}
		response, err := addressBookService.GetAddressBook(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get address book: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)
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
	getAddressBookCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

}
