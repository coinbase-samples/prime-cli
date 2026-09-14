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

package futures

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/futures"

	"github.com/spf13/cobra"
)

var getDerivativePositionsCmd = &cobra.Command{
	Use:   "get-derivative-positions",
	Short: "Gets active derivative positions for a portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := futures.NewFuturesService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get portfolio ID: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &futures.GetDerivativePositionsRequest{
			PortfolioId: portfolioId,
			ProductId:   utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
		}

		response, err := svc.GetDerivativePositions(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get derivative positions: %w", err)
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
	Cmd.AddCommand(getDerivativePositionsCmd)

	utils.AddPortfolioIdFlag(getDerivativePositionsCmd)
	getDerivativePositionsCmd.Flags().String(utils.ProductIdFlag, "", "Optional product ID to filter positions")
}
