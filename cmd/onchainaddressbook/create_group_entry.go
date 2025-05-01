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
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/onchainaddressbook"
	"github.com/spf13/cobra"
)

var createOnchainAddressBookEntryCmd = &cobra.Command{
	Use:   "create-group-entry",
	Short: "Create an onchain address book entry",
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

		networkTypeMap := map[string]model.OnchainNetworkType{
			"NETWORK_TYPE_EVM":         model.OnchainNetworkTypeEvm,
			"NETWORK_TYPE_SOLANA":      model.OnchainNetworkTypeSolana,
			"NETWORK_TYPE_UNSPECIFIED": model.OnchainNetworkTypeUnspecified,
		}

		networkTypeStr := utils.GetFlagStringValue(cmd, utils.NetworkTypeFlag)
		networkType, ok := networkTypeMap[networkTypeStr]
		if !ok {
			return fmt.Errorf("invalid network type: %s", networkTypeStr)
		}

		address := utils.GetFlagStringValue(cmd, utils.AddressFlag)
		if address == "" {
			return fmt.Errorf("address is required")
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &onchainaddressbook.CreateOnchainAddressBookEntryRequest{
			PortfolioId: portfolioId,
			AddressGroup: &model.OnchainAddressGroup{
				Id:          utils.GetFlagStringValue(cmd, utils.GenericIdFlag),
				Name:        utils.GetFlagStringValue(cmd, utils.NameFlag),
				NetworkType: networkType,
				Addresses: []*model.OnchainAddress{
					{
						Address: address,
						Name:    utils.GetFlagStringValue(cmd, utils.NameFlag),
					},
				},
			},
		}

		response, err := onchainService.CreateOnchainAddressBookEntry(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create onchain address book entry: %w", err)
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
	Cmd.AddCommand(createOnchainAddressBookEntryCmd)

	utils.AddPortfolioIdFlag(createOnchainAddressBookEntryCmd)
	createOnchainAddressBookEntryCmd.Flags().StringP(utils.GenericIdFlag, "i", "", "Address group ID (Required)")
	createOnchainAddressBookEntryCmd.Flags().StringP(utils.AddressFlag, "a", "", "Address (Required)")
	createOnchainAddressBookEntryCmd.Flags().StringP(utils.NetworkTypeFlag, "t", "", "Network type (Required)")
	createOnchainAddressBookEntryCmd.Flags().StringP(utils.NameFlag, "n", "", "Name for the address group")

	createOnchainAddressBookEntryCmd.MarkFlagRequired(utils.GenericIdFlag)
	createOnchainAddressBookEntryCmd.MarkFlagRequired(utils.AddressFlag)
	createOnchainAddressBookEntryCmd.MarkFlagRequired(utils.NetworkTypeFlag)
}
