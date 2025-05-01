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
	"github.com/spf13/cobra"
)

var createAddressBookEntryCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an address address book entry.",
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

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &addressbook.CreateAddressBookEntryRequest{
			PortfolioId:       portfolioId,
			Address:           utils.GetFlagStringValue(cmd, utils.AddressFlag),
			Symbol:            utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			Name:              utils.GetFlagStringValue(cmd, utils.NameFlag),
			AccountIdentifier: utils.GetFlagStringValue(cmd, utils.AccountIdFlag),
		}

		response, err := addressBookService.CreateAddressBookEntry(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create address book entry: %w", err)
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
	Cmd.AddCommand(createAddressBookEntryCmd)

	createAddressBookEntryCmd.Flags().StringP(utils.AddressFlag, "a", "", "The address to add to the address book (Required)")
	createAddressBookEntryCmd.Flags().StringP(utils.SymbolFlag, "s", "", "The currency symbol (Required)")
	createAddressBookEntryCmd.Flags().StringP(utils.NameFlag, "n", "", "Name for the address book entry (Required)")
	createAddressBookEntryCmd.Flags().StringP(utils.AccountIdFlag, "i", "", "Account identifier for the address")
	utils.AddPortfolioIdFlag(createAddressBookEntryCmd)

	createAddressBookEntryCmd.MarkFlagRequired(utils.AddressFlag)
	createAddressBookEntryCmd.MarkFlagRequired(utils.SymbolFlag)
	createAddressBookEntryCmd.MarkFlagRequired(utils.NameFlag)
}
