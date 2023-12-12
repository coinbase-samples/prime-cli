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
	"github.com/coinbase-samples/prime-sdk-go"

	"github.com/spf13/cobra"
)

var listPortfolioBalancesCmd = &cobra.Command{
	Use:   "list-portfolio-balances",
	Short: "Lists balances that meet filter criteria",
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		symbols, err := cmd.Flags().GetStringSlice(utils.SymbolsFlag)
		if err != nil {
			return fmt.Errorf("cannot get symbols: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.ListPortfolioBalancesRequest{
			PortfolioId: portfolioId,
			Type:        utils.GetFlagStringValue(cmd, utils.TypeFlag),
			Symbols:     symbols,
		}

		response, err := client.ListPortfolioBalances(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list portfolio balances: %w", err)
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
	rootCmd.AddCommand(listPortfolioBalancesCmd)

	listPortfolioBalancesCmd.Flags().StringArrayP(utils.TypeFlag, "t", []string{}, "Type of balance")
	listPortfolioBalancesCmd.Flags().StringSliceP(utils.SymbolsFlag, "s", []string{}, "List of symbols")
	listPortfolioBalancesCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	listPortfolioBalancesCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
}
