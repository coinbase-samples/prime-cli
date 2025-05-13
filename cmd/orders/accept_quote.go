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

package orders

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/orders"
	"github.com/spf13/cobra"
)

var acceptQuoteCmd = &cobra.Command{
	Use:   "accept-quote",
	Short: "Accept a quote request",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := orders.NewOrdersService(client)

		clientOrderId := utils.GetFlagStringValue(cmd, utils.ClientOrderIdFlag)
		if clientOrderId == "" {
			clientOrderId = utils.NewUuidStr()
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		request := &orders.AcceptQuoteRequest{
			PortfolioId:   portfolioId,
			ProductId:     utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			QuoteId:       utils.GetFlagStringValue(cmd, utils.QuoteIdFlag),
			ClientOrderId: clientOrderId,
			Side:          utils.GetFlagStringValue(cmd, utils.SideFlag),
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.AcceptQuote(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot accept quote: %w", err)
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
	Cmd.AddCommand(acceptQuoteCmd)
	utils.AddPortfolioIdFlag(acceptQuoteCmd)
	utils.AddProductIdFlag(acceptQuoteCmd)
	utils.AddOrderSideFlag(acceptQuoteCmd)
	utils.AddClientOrderId(acceptQuoteCmd)

	acceptQuoteCmd.Flags().String(utils.QuoteIdFlag, "", "The quote id returned by the create quote request")

	acceptQuoteCmd.MarkFlagRequired(utils.SideFlag)
	acceptQuoteCmd.MarkFlagRequired(utils.ProductIdFlag)
	acceptQuoteCmd.MarkFlagRequired(utils.QuoteIdFlag)
}
