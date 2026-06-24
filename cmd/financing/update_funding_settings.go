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

package financing

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	prime "github.com/coinbase/prime-sdk-go/financing"
	"github.com/spf13/cobra"
)

const (
	designatedFundingPortfolioIdFlag = "designated-funding-portfolio-id"
	automaticConversionEnabledFlag   = "automatic-conversion-enabled"
	automaticLoanEnabledFlag         = "automatic-loan-enabled"
	automaticExcessReturnEnabledFlag = "automatic-excess-return-enabled"
	excessFundsTargetAmountFlag      = "excess-funds-target-amount"
)

var updateFundingSettingsCmd = &cobra.Command{
	Use:   "update-funding-settings",
	Short: "Updates FCM funding settings for an entity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return err
		}

		automaticConversionEnabled, err := cmd.Flags().GetBool(automaticConversionEnabledFlag)
		if err != nil {
			return err
		}

		automaticLoanEnabled, err := cmd.Flags().GetBool(automaticLoanEnabledFlag)
		if err != nil {
			return err
		}

		automaticExcessReturnEnabled, err := cmd.Flags().GetBool(automaticExcessReturnEnabledFlag)
		if err != nil {
			return err
		}

		request := &prime.UpdateFundingSettingsRequest{
			EntityId:                     entityId,
			DesignatedFundingPortfolioId: utils.GetFlagStringValue(cmd, designatedFundingPortfolioIdFlag),
			AutomaticConversionEnabled:   automaticConversionEnabled,
			AutomaticLoanEnabled:         automaticLoanEnabled,
			AutomaticExcessReturnEnabled: automaticExcessReturnEnabled,
			ExcessFundsTargetAmount:      utils.GetFlagStringValue(cmd, excessFundsTargetAmountFlag),
		}

		response, err := updateFundingSettings(svc, request)
		if err != nil {
			return err
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func updateFundingSettings(
	svc prime.FinancingService,
	req *prime.UpdateFundingSettingsRequest,
) (*prime.UpdateFundingSettingsResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.UpdateFundingSettings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot update funding settings: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(updateFundingSettingsCmd)

	utils.AddEntityIdFlag(updateFundingSettingsCmd)
	updateFundingSettingsCmd.Flags().String(designatedFundingPortfolioIdFlag, "", "Derivatives funding portfolio ID")
	updateFundingSettingsCmd.Flags().Bool(automaticConversionEnabledFlag, false, "Convert USDC to USD automatically to meet FCM margin calls")
	updateFundingSettingsCmd.Flags().Bool(automaticLoanEnabledFlag, false, "Allow Coinbase affiliates to initiate loans to meet FCM margin calls")
	updateFundingSettingsCmd.Flags().Bool(automaticExcessReturnEnabledFlag, false, "Sweep FCM balance above margin requirements back to the derivatives funding portfolio")
	updateFundingSettingsCmd.Flags().String(excessFundsTargetAmountFlag, "", "Target amount to maintain in the futures account above margin requirements")

	updateFundingSettingsCmd.MarkFlagRequired(designatedFundingPortfolioIdFlag)
}
