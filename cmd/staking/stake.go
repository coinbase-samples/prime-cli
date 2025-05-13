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

package staking

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	primeStaking "github.com/coinbase-samples/prime-sdk-go/staking"
	"github.com/spf13/cobra"
)

var createStakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Creates a request to stake or delegate funds to a validator",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := primeStaking.NewStakingService(client)

		idempotencyKey := utils.GetFlagStringValue(cmd, utils.IdempotencyKeyFlag)
		if idempotencyKey == "" {
			idempotencyKey = utils.NewUuidStr()
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		request := &primeStaking.CreateStakeRequest{
			PortfolioId:    portfolioId,
			WalletId:       utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
			IdempotencyKey: idempotencyKey,
		}

		amount := utils.GetFlagStringValue(cmd, utils.AmountFlag)
		if len(amount) > 0 {
			request.Inputs = primeStaking.CreateStakeInputs{Amount: amount}
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.CreateStake(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create staking request: %w", err)
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
	Cmd.AddCommand(createStakeCmd)
	utils.AddPortfolioIdFlag(createStakeCmd)
	utils.AddWalletIdFlag(createStakeCmd)
	utils.AddIdempotencyKeyFlag(createStakeCmd)

	createStakeCmd.Flags().String(utils.AmountFlag, "", "Optional amount to stake. If omitted, the wallet will stake or unstake the maximum amount available")
}
