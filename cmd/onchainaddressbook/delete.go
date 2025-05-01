/**
 * Copyright 2025-present Coinbase Global, Inc.
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

package onchainaddressbook

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/onchainaddressbook"
	"github.com/spf13/cobra"
)

var deleteOnchainAddressBookEntryCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an onchain address book entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		onchainService := onchainaddressbook.NewOnchainAddressBookService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		addressGroupId := utils.GetFlagStringValue(cmd, utils.GenericIdFlag)
		if addressGroupId == "" {
			return fmt.Errorf("address group ID is required")
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &onchainaddressbook.DeleteOnchainAddressBookEntryRequest{
			PortfolioId:    portfolioId,
			AddressGroupId: addressGroupId,
		}

		response, err := onchainService.DeleteOnchainAddressBookEntry(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot delete onchain address book entry: %w", err)
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
	Cmd.AddCommand(deleteOnchainAddressBookEntryCmd)

	deleteOnchainAddressBookEntryCmd.Flags().StringP(utils.PortfolioIdFlag, "p", "", "Portfolio ID. Uses environment variable if blank")
	deleteOnchainAddressBookEntryCmd.Flags().StringP(utils.GenericIdFlag, "i", "", "Address group ID (Required)")

	deleteOnchainAddressBookEntryCmd.MarkFlagRequired(utils.GenericIdFlag)
}
