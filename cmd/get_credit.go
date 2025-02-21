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
	"github.com/coinbase-samples/prime-sdk-go/portfolios"

	"github.com/spf13/cobra"
)

var getCreditCmd = &cobra.Command{
	Use:   "get-credit",
	Short: "Get portfolio credit information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfoliosService := portfolios.NewPortfoliosService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &portfolios.GetPortfolioCreditRequest{
			Id: portfolioId,
		}

		response, err := portfoliosService.GetPortfolioCredit(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get portfolio credit: %w", err)
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
	rootCmd.AddCommand(getCreditCmd)

	getCreditCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getCreditCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

}
