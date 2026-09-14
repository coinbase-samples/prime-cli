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

const liquidationIdFlag = "liquidation-id"

var getCrossMarginLiquidationCmd = &cobra.Command{
	Use:   "get-cross-margin-liquidation",
	Short: "Gets detailed cross-margin liquidation data for an entity",
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

		request := &prime.GetCrossMarginLiquidationRequest{
			EntityId:      entityId,
			LiquidationId: utils.GetFlagStringValue(cmd, liquidationIdFlag),
		}

		response, err := getCrossMarginLiquidation(svc, request)
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

func getCrossMarginLiquidation(
	svc prime.FinancingService,
	req *prime.GetCrossMarginLiquidationRequest,
) (*prime.GetCrossMarginLiquidationResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.GetCrossMarginLiquidation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get cross margin liquidation: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(getCrossMarginLiquidationCmd)

	utils.AddEntityIdFlag(getCrossMarginLiquidationCmd)
	getCrossMarginLiquidationCmd.Flags().String(liquidationIdFlag, "", "Optional liquidation ID")
}
