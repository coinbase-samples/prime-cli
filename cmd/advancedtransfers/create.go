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

package advancedtransfers

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/advancedtransfers"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/spf13/cobra"
)

const (
	currencyFlag                = "currency"
	fundMovementIdFlag          = "fund-movement-id"
	sourceTypeFlag              = "source-type"
	sourceValueFlag             = "source-value"
	sourceAddressFlag           = "source-address"
	sourceAccountIdentifierFlag = "source-account-identifier"
	targetTypeFlag              = "target-type"
	targetValueFlag             = "target-value"
	targetAddressFlag           = "target-address"
	targetAccountIdentifierFlag = "target-account-identifier"
	referenceIdFlag             = "reference-id"
	settlementDateFlag          = "settlement-date"
	tradeDateFlag               = "trade-date"
	settlementTimeFlag          = "settlement-time"
)

var createAdvancedTransferCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new advanced transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := advancedtransfers.NewAdvancedTransfersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		fundMovementId := utils.GetFlagStringValue(cmd, fundMovementIdFlag)
		if fundMovementId == "" {
			fundMovementId = utils.GetFlagStringValue(cmd, utils.IdempotencyKeyFlag)
		}
		if fundMovementId == "" {
			fundMovementId = utils.NewUuidStr()
		}

		transfer := &model.AdvancedTransfer{
			Type: model.AdvancedTransferType(utils.GetFlagStringValue(cmd, utils.TransferTypeFlag)),
		}

		movement := &model.FundMovement{
			Id:       fundMovementId,
			Currency: utils.GetFlagStringValue(cmd, currencyFlag),
			Amount:   utils.GetFlagStringValue(cmd, utils.AmountFlag),
		}

		if source := buildTransferLocation(
			utils.GetFlagStringValue(cmd, sourceTypeFlag),
			utils.GetFlagStringValue(cmd, sourceValueFlag),
			utils.GetFlagStringValue(cmd, sourceAddressFlag),
			utils.GetFlagStringValue(cmd, sourceAccountIdentifierFlag),
		); source != nil {
			movement.Source = source
		}

		if target := buildTransferLocation(
			utils.GetFlagStringValue(cmd, targetTypeFlag),
			utils.GetFlagStringValue(cmd, targetValueFlag),
			utils.GetFlagStringValue(cmd, targetAddressFlag),
			utils.GetFlagStringValue(cmd, targetAccountIdentifierFlag),
		); target != nil {
			movement.Target = target
		}

		transfer.FundMovements = []*model.FundMovement{movement}

		referenceId := utils.GetFlagStringValue(cmd, referenceIdFlag)
		settlementDate := utils.GetFlagStringValue(cmd, settlementDateFlag)
		tradeDate := utils.GetFlagStringValue(cmd, tradeDateFlag)
		settlementTime := utils.GetFlagStringValue(cmd, settlementTimeFlag)

		if referenceId != "" || settlementDate != "" || tradeDate != "" || settlementTime != "" {
			transfer.BlindMatchMetadata = &model.BlindMatchMetadata{
				ReferenceId:    referenceId,
				SettlementDate: settlementDate,
				TradeDate:      tradeDate,
				SettlementTime: settlementTime,
			}
		}

		request := &advancedtransfers.CreateAdvancedTransferRequest{
			PortfolioId:      portfolioId,
			AdvancedTransfer: transfer,
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.CreateAdvancedTransfer(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create advanced transfer: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func buildTransferLocation(typeStr, value, address, accountIdentifier string) *model.TransferLocation {
	if typeStr == "" && value == "" && address == "" && accountIdentifier == "" {
		return nil
	}
	return &model.TransferLocation{
		Type:              model.TransferLocationType(typeStr),
		Value:             value,
		Address:           address,
		AccountIdentifier: accountIdentifier,
	}
}

func init() {
	Cmd.AddCommand(createAdvancedTransferCmd)

	utils.AddPortfolioIdFlag(createAdvancedTransferCmd)
	utils.AddIdempotencyKeyFlag(createAdvancedTransferCmd)

	createAdvancedTransferCmd.Flags().String(utils.TransferTypeFlag, "", "Advanced transfer type, e.g. ADVANCED_TRANSFER_TYPE_BLIND_MATCH (Required)")
	createAdvancedTransferCmd.Flags().String(utils.AmountFlag, "", "Amount to transfer (Required)")
	createAdvancedTransferCmd.Flags().String(currencyFlag, "", "Currency symbol for the transfer (Required)")

	createAdvancedTransferCmd.Flags().String(fundMovementIdFlag, "", "Optional client-supplied fund movement ID. Defaults to --idempotency-key or a generated UUID")

	createAdvancedTransferCmd.Flags().String(sourceTypeFlag, "", "Source location type, e.g. WALLET, COUNTERPARTY_ID")
	createAdvancedTransferCmd.Flags().String(sourceValueFlag, "", "Source location value (e.g. wallet ID or counterparty ID)")
	createAdvancedTransferCmd.Flags().String(sourceAddressFlag, "", "Source blockchain address")
	createAdvancedTransferCmd.Flags().String(sourceAccountIdentifierFlag, "", "Source account identifier")

	createAdvancedTransferCmd.Flags().String(targetTypeFlag, "", "Target location type, e.g. WALLET, COUNTERPARTY_ID")
	createAdvancedTransferCmd.Flags().String(targetValueFlag, "", "Target location value (e.g. wallet ID or counterparty ID)")
	createAdvancedTransferCmd.Flags().String(targetAddressFlag, "", "Target blockchain address")
	createAdvancedTransferCmd.Flags().String(targetAccountIdentifierFlag, "", "Target account identifier")

	createAdvancedTransferCmd.Flags().String(referenceIdFlag, "", "Blind match reference ID")
	createAdvancedTransferCmd.Flags().String(settlementDateFlag, "", "Blind match settlement date")
	createAdvancedTransferCmd.Flags().String(tradeDateFlag, "", "Blind match trade date")
	createAdvancedTransferCmd.Flags().String(settlementTimeFlag, "", "Blind match settlement time")

	createAdvancedTransferCmd.MarkFlagRequired(utils.TransferTypeFlag)
	createAdvancedTransferCmd.MarkFlagRequired(utils.AmountFlag)
	createAdvancedTransferCmd.MarkFlagRequired(currencyFlag)
}
