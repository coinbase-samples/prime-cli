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

package financing

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	prime "github.com/coinbase-samples/prime-sdk-go/financing"
	"github.com/spf13/cobra"
)

var getWithdrawalPowerCmd = &cobra.Command{
	Use:   "get-withdrawal-power",
	Short: "Get withdrawal power information for a portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		symbol, err := cmd.Flags().GetString("symbol")
		if err != nil {
			return err
		}

		request := &prime.GetWithdrawalPowerRequest{
			PortfolioId: portfolioId,
			Symbol:      symbol,
		}

		response, err := getWithdrawalPower(svc, request)
		if err != nil {
			return err
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func getWithdrawalPower(
	svc prime.FinancingService,
	req *prime.GetWithdrawalPowerRequest,
) (*prime.GetWithdrawalPowerResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.GetWithdrawalPower(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get withdrawal power: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(getWithdrawalPowerCmd)

	getWithdrawalPowerCmd.Flags().String("symbol", "", "The symbol for the asset")
	getWithdrawalPowerCmd.MarkFlagRequired("symbol")

	utils.AddPortfolioIdFlag(getWithdrawalPowerCmd)
}
