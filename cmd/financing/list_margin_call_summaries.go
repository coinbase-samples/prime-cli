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

var listMarginCallSummariesCmd = &cobra.Command{
	Use:   "list-margin-call-summaries",
	Short: "List margin call summaries for an entity",
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

		startDate, err := cmd.Flags().GetString("start-date")
		if err != nil {
			return err
		}

		endDate, err := cmd.Flags().GetString("end-date")
		if err != nil {
			return err
		}

		request := &prime.ListMarginCallSummariesRequest{
			EntityId:  entityId,
			StartDate: startDate,
			EndDate:   endDate,
		}

		response, err := listMarginCallSummaries(svc, request)
		if err != nil {
			return err
		}

		if err := utils.PrintJsonDocs(cmd, response.MarginSummaries); err != nil {
			return err
		}

		return nil
	},
}

func listMarginCallSummaries(
	svc prime.FinancingService,
	req *prime.ListMarginCallSummariesRequest,
) (*prime.ListMarginCallSummariesResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListMarginCallSummaries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list margin call summaries: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listMarginCallSummariesCmd)

	listMarginCallSummariesCmd.Flags().String("entity-id", "", "Entity ID")
	listMarginCallSummariesCmd.MarkFlagRequired("entity-id")

	listMarginCallSummariesCmd.Flags().String("start-date", "", "Start date in RFC3339 format")
	listMarginCallSummariesCmd.Flags().String("end-date", "", "End date in RFC3339 format")
}
