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
	"encoding/json"
	"fmt"
	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var createConversionCmd = &cobra.Command{
	Use:   "create-conversion",
	Short: "Convert USD to USDC and vice-versa.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		idempotencyKey := uuid.New().String()

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CreateConversionRequest{
			PortfolioId:         client.Credentials.PortfolioId,
			SourceWalletId:      utils.GetFlagStringValue(cmd, utils.SourceWalletIdFlag),
			SourceSymbol:        utils.GetFlagStringValue(cmd, utils.SourceSymbolFlag),
			DestinationWalletId: utils.GetFlagStringValue(cmd, utils.DestinationWalletIdFlag),
			DestinationSymbol:   utils.GetFlagStringValue(cmd, utils.DestinationSymbolFlag),
			IdempotencyKey:      idempotencyKey,
			Amount:              utils.GetFlagStringValue(cmd, utils.AmountFlag),
		}

		response, err := client.CreateConversion(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating conversion: %v", err)
		}

		jsonResponse, err := json.MarshalIndent(response, "", utils.JsonIndent)
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %v", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createConversionCmd)

	createConversionCmd.Flags().StringP(utils.SourceWalletIdFlag, "i", "", "ID of the source wallet (Required)")
	createConversionCmd.Flags().StringP(utils.SourceSymbolFlag, "s", "", "Symbol of the source wallet (Required)")
	createConversionCmd.Flags().StringP(utils.DestinationWalletIdFlag, "d", "", "ID of the destination wallet (Required)")
	createConversionCmd.Flags().StringP(utils.DestinationSymbolFlag, "f", "", "Symbol of the destination wallet (Required)")
	createConversionCmd.Flags().StringP(utils.AmountFlag, "a", "", "Conversion size (Required)")

	createConversionCmd.MarkFlagRequired(utils.SourceWalletIdFlag)
	createConversionCmd.MarkFlagRequired(utils.SourceSymbolFlag)
	createConversionCmd.MarkFlagRequired(utils.DestinationWalletIdFlag)
	createConversionCmd.MarkFlagRequired(utils.DestinationSymbolFlag)
	createConversionCmd.MarkFlagRequired(utils.AmountFlag)

	createConversionCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateUUIDFlag(cmd, utils.SourceWalletIdFlag); err != nil {
			return err
		}
		if err := utils.ValidateUUIDFlag(cmd, utils.DestinationWalletIdFlag); err != nil {
			return err
		}
		return nil
	}
}