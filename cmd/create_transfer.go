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

var createTransferCmd = &cobra.Command{
	Use:   "create-transfer",
	Short: "Create an internal transfer.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		idempotencyKey := uuid.New().String()

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CreateWalletTransferRequest{
			PortfolioId:         client.Credentials.PortfolioId,
			SourceWalletId:      utils.GetFlagStringValue(cmd, utils.SourceWalletIdFlag),
			Symbol:              utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			DestinationWalletId: utils.GetFlagStringValue(cmd, utils.DestinationWalletIdFlag),
			IdempotencyKey:      idempotencyKey,
			Amount:              utils.GetFlagStringValue(cmd, utils.AmountFlag),
		}

		response, err := client.CreateWalletTransfer(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating transfer: %v", err)
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
	rootCmd.AddCommand(createTransferCmd)

	createTransferCmd.Flags().StringP(utils.SourceWalletIdFlag, "1", "", "ID of the source wallet (Required)")
	createTransferCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Symbol of the asset to be transferred (Required)")
	createTransferCmd.Flags().StringP(utils.DestinationWalletIdFlag, "2", "", "ID of the destination wallet (Required)")
	createTransferCmd.Flags().StringP(utils.AmountFlag, "a", "", "Conversion size (Required)")

	createTransferCmd.MarkFlagRequired(utils.SourceWalletIdFlag)
	createTransferCmd.MarkFlagRequired(utils.SymbolFlag)
	createTransferCmd.MarkFlagRequired(utils.DestinationWalletIdFlag)
	createTransferCmd.MarkFlagRequired(utils.AmountFlag)

	createTransferCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateUUIDFlag(cmd, utils.SourceWalletIdFlag); err != nil {
			return err
		}
		if err := utils.ValidateUUIDFlag(cmd, utils.DestinationWalletIdFlag); err != nil {
			return err
		}
		return nil
	}
}