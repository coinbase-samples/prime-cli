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

package products

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/products"

	"github.com/spf13/cobra"
)

var getProductCandlesCmd = &cobra.Command{
	Use:   "get-candles",
	Short: "Gets candlestick data for a product",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := products.NewProductsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		startTime, endTime, err := utils.GetStartEndFlagsAsTime(cmd)
		if err != nil {
			return fmt.Errorf("cannot parse start/end times: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &products.GetProductCandlesRequest{
			PortfolioId: portfolioId,
			ProductId:   utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			StartTime:   startTime,
			EndTime:     endTime,
			Granularity: model.CandleGranularity(utils.GetFlagStringValue(cmd, utils.GranularityFlag)),
		}

		response, err := svc.GetProductCandles(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get product candles: %w", err)
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
	Cmd.AddCommand(getProductCandlesCmd)

	getProductCandlesCmd.Flags().String(utils.GranularityFlag, "", "Candle granularity (e.g. ONE_MINUTE, ONE_HOUR, ONE_DAY) (Required)")

	utils.AddPortfolioIdFlag(getProductCandlesCmd)
	utils.AddProductIdFlag(getProductCandlesCmd)
	utils.AddStartEndFlags(getProductCandlesCmd)

	getProductCandlesCmd.MarkFlagRequired(utils.ProductIdFlag)
	getProductCandlesCmd.MarkFlagRequired(utils.GranularityFlag)
	getProductCandlesCmd.MarkFlagRequired(utils.StartFlag)
	getProductCandlesCmd.MarkFlagRequired(utils.EndFlag)
}
