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

package addressbook

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/addressbook"
	"github.com/coinbase-samples/prime-sdk-go/model"

	"github.com/spf13/cobra"
)

var listAddressBookCmd = &cobra.Command{
	Use:   "list",
	Short: "List the address book entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := addressbook.NewAddressBookService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		symbol := utils.GetFlagStringValue(cmd, utils.SymbolFlag)
		search := utils.GetFlagStringValue(cmd, utils.SearchFlag)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := getAddressBookEntries(svc, portfolioId, symbol, search, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Addresses); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func getAddressBookEntries(
	svc addressbook.AddressBookService,
	portfolioId,
	symbol,
	search string,
	pagination *model.PaginationParams,
) (*addressbook.GetAddressBookResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &addressbook.GetAddressBookRequest{
		PortfolioId: portfolioId,
		Symbol:      symbol,
		Search:      search,
		Pagination:  pagination,
	}
	response, err := svc.GetAddressBook(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list address book: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listAddressBookCmd)
	listAddressBookCmd.Flags().String(utils.SymbolFlag, "", "Currency symbol for filtering address book entries")
	listAddressBookCmd.Flags().String(utils.SearchFlag, "", "Search term for filtering address book entries")
	utils.AddPortfolioIdFlag(listAddressBookCmd)
	utils.AddPaginationFlags(listAddressBookCmd, true)
}
