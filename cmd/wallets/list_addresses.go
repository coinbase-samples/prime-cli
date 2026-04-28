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

package wallets

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/wallets"

	"github.com/spf13/cobra"
)

var listWalletAddressesCmd = &cobra.Command{
	Use:   "list-addresses",
	Short: "Lists addresses for a wallet",
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

		walletId := utils.GetFlagStringValue(cmd, utils.WalletIdFlag)
		networkId := utils.GetFlagStringValue(cmd, utils.NetworkIdFlag)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listWalletAddresses(svc, portfolioId, walletId, networkId, paginationParams)
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

func listWalletAddresses(
	svc wallets.WalletsService,
	portfolioId,
	walletId,
	networkId string,
	pagination *model.PaginationParams,
) (*wallets.ListWalletAddressesResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &wallets.ListWalletAddressesRequest{
		PortfolioId: portfolioId,
		WalletId:    walletId,
		NetworkId:   networkId,
		Pagination:  pagination,
	}

	response, err := svc.ListWalletAddresses(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list wallet addresses: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listWalletAddressesCmd)

	listWalletAddressesCmd.Flags().String(utils.NetworkIdFlag, "", "Filter by network ID")

	utils.AddPortfolioIdFlag(listWalletAddressesCmd)
	utils.AddWalletIdFlag(listWalletAddressesCmd)
	utils.AddPaginationFlags(listWalletAddressesCmd, true)
}
