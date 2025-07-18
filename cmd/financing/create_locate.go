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

var createLocateCmd = &cobra.Command{
	Use:   "create-locate",
	Short: "Create a new locate for a portfolio and assset",
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

		amount, err := cmd.Flags().GetString("amount")
		if err != nil {
			return err
		}

		locateDate, err := cmd.Flags().GetString("date")
		if err != nil {
			return err
		}

		request := &prime.CreateLocateRequest{
			PortfolioId: portfolioId,
			Symbol:      symbol,
			Amount:      amount,
			LocateDate:  locateDate,
		}

		response, err := createLocate(svc, request)
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

func createLocate(
	svc prime.FinancingService,
	req *prime.CreateLocateRequest,
) (*prime.CreateLocateResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.CreateLocate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot create new locate: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(createLocateCmd)

	createLocateCmd.Flags().String("symbol", "", "The symbol for the asset")
	createLocateCmd.MarkFlagRequired("symbol")

	createLocateCmd.Flags().String("amount", "", "The locate amount")
	createLocateCmd.MarkFlagRequired("amount")

	createLocateCmd.Flags().String("date", "", "The target date of the locate (YYYY-MM-DD)")
	createLocateCmd.MarkFlagRequired("locate-date")

	utils.AddPortfolioIdFlag(createLocateCmd)
}
