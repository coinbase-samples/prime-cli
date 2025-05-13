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
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/orders"
	"github.com/spf13/cobra"
)

var createQuoteCmd = &cobra.Command{
	Use:   "create-quote",
	Short: "Create a quote request",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := orders.NewOrdersService(client)

		clientQuoteId := utils.GetFlagStringValue(cmd, utils.ClientQuoteIdFlag)
		if clientQuoteId == "" {
			clientQuoteId = utils.NewUuidStr()
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		request := &orders.CreateQuoteRequest{
			PortfolioId:    portfolioId,
			ProductId:      utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			ClientQuoteId:  clientQuoteId,
			Side:           model.OrderSide(utils.GetFlagStringValue(cmd, utils.SideFlag)),
			BaseQuantity:   utils.GetFlagStringValue(cmd, utils.BaseQuantityFlag),
			QuoteValue:     utils.GetFlagStringValue(cmd, utils.QuoteValueFlag),
			LimitPrice:     utils.GetFlagStringValue(cmd, utils.LimitPriceFlag),
			SettleCurrency: utils.GetFlagStringValue(cmd, utils.SettleCurrencyFlag),
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.CreateQuoteRequest(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create quote request: %w", err)
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
	Cmd.AddCommand(createQuoteCmd)
	utils.AddPortfolioIdFlag(createQuoteCmd)
	utils.AddProductIdFlag(createQuoteCmd)
	utils.AddOrderSideFlag(createQuoteCmd)
	utils.AddLimitPriceFlag(createQuoteCmd)
	utils.AddQuoteValueFlag(createQuoteCmd)
	utils.AddBaseQuantityFlag(createQuoteCmd)

	createQuoteCmd.Flags().String(utils.ClientQuoteIdFlag, "", "A client-generated order ID used for reference purposes")
	createQuoteCmd.Flags().String(utils.SettleCurrencyFlag, "", "The settle currency flag")

	createQuoteCmd.MarkFlagRequired(utils.SideFlag)
	createQuoteCmd.MarkFlagRequired(utils.ProductIdFlag)
}
