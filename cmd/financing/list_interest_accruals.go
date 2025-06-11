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

package financing

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	prime "github.com/coinbase-samples/prime-sdk-go/financing"
	"github.com/spf13/cobra"
)

var listInterestAccrualsCmd = &cobra.Command{
	Use:   "list-interest-accruals",
	Short: "List interest accruals for an entity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		entityId, err := cmd.Flags().GetString("entity-id")
		if err != nil {
			return err
		}

		portfolioId, err := cmd.Flags().GetString("portfolio-id")
		if err != nil {
			return err
		}

		startDate, err := cmd.Flags().GetString("start-date")
		if err != nil {
			return err
		}

		endDate, err := cmd.Flags().GetString("end-date")
		if err != nil {
			return err
		}

		request := &prime.ListInterestAccrualsRequest{
			EntityId:    entityId,
			PortfolioId: portfolioId,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		response, err := listInterestAccruals(svc, request)
		if err != nil {
			return err
		}

		if err := utils.PrintJsonDocs(cmd, response.Accruals); err != nil {
			return err
		}

		return nil
	},
}

func listInterestAccruals(
	svc prime.FinancingService,
	req *prime.ListInterestAccrualsRequest,
) (*prime.ListInterestAccrualsResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListInterestAccruals(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list interest accruals: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listInterestAccrualsCmd)

	listInterestAccrualsCmd.Flags().String("entity-id", "", "Entity ID")
	listInterestAccrualsCmd.MarkFlagRequired("entity-id")

	listInterestAccrualsCmd.Flags().String("portfolio-id", "", "Portfolio ID")

	listInterestAccrualsCmd.Flags().String("start-date", "", "Start date in RFC3339 format")
	listInterestAccrualsCmd.Flags().String("end-date", "", "End date in RFC3339 format")
}
