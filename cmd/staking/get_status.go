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

package staking

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	primeStaking "github.com/coinbase-samples/prime-sdk-go/staking"
	"github.com/spf13/cobra"
)

var getStakingStatusCmd = &cobra.Command{
	Use:   "get-status",
	Short: "Gets staking status for a wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := primeStaking.NewStakingService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		request := &primeStaking.GetStakingStatusRequest{
			PortfolioId: portfolioId,
			WalletId:    utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.GetStakingStatus(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get staking status: %w", err)
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
	Cmd.AddCommand(getStakingStatusCmd)
	utils.AddPortfolioIdFlag(getStakingStatusCmd)
	utils.AddWalletIdFlag(getStakingStatusCmd)
}
