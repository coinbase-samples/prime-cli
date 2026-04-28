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

package transactions

import (
	"encoding/json"
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/transactions"

	"github.com/spf13/cobra"
)

const (
	originatorFlag                    = "originator"
	beneficiaryFlag                   = "beneficiary"
	isSelfFlag                        = "is-self"
	optOutOfOwnershipVerificationFlag = "opt-out-of-ownership-verification"
)

var submitDepositTravelRuleDataCmd = &cobra.Command{
	Use:   "submit-deposit-travel-rule-data",
	Short: "Submits travel rule data for a deposit transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		transactionsService := transactions.NewTransactionsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		request := &transactions.SubmitDepositTravelRuleDataRequest{
			PortfolioId:                   portfolioId,
			TransactionId:                 utils.GetFlagStringValue(cmd, utils.TransactionIdFlag),
			IsSelf:                        utils.GetFlagBoolValue(cmd, isSelfFlag),
			OptOutOfOwnershipVerification: utils.GetFlagBoolValue(cmd, optOutOfOwnershipVerificationFlag),
		}

		if originatorJson := utils.GetFlagStringValue(cmd, originatorFlag); originatorJson != "" {
			var originator model.TravelRuleParty
			if err := json.Unmarshal([]byte(originatorJson), &originator); err != nil {
				return fmt.Errorf("invalid originator format: %w", err)
			}
			request.Originator = &originator
		}

		if beneficiaryJson := utils.GetFlagStringValue(cmd, beneficiaryFlag); beneficiaryJson != "" {
			var beneficiary model.TravelRuleParty
			if err := json.Unmarshal([]byte(beneficiaryJson), &beneficiary); err != nil {
				return fmt.Errorf("invalid beneficiary format: %w", err)
			}
			request.Beneficiary = &beneficiary
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := transactionsService.SubmitDepositTravelRuleData(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot submit deposit travel rule data: %w", err)
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
	Cmd.AddCommand(submitDepositTravelRuleDataCmd)

	utils.AddPortfolioIdFlag(submitDepositTravelRuleDataCmd)
	submitDepositTravelRuleDataCmd.Flags().String(utils.TransactionIdFlag, "", "Transaction ID (Required)")
	submitDepositTravelRuleDataCmd.Flags().String(originatorFlag, "", "JSON string of the originator travel rule party")
	submitDepositTravelRuleDataCmd.Flags().String(beneficiaryFlag, "", "JSON string of the beneficiary travel rule party")
	submitDepositTravelRuleDataCmd.Flags().Bool(isSelfFlag, false, "Whether the deposit is to self")
	submitDepositTravelRuleDataCmd.Flags().Bool(optOutOfOwnershipVerificationFlag, false, "Whether to opt out of ownership verification")

	submitDepositTravelRuleDataCmd.MarkFlagRequired(utils.TransactionIdFlag)
}
