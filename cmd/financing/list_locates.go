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

var listLocatesCmd = &cobra.Command{
	Use:   "list-locates",
	Short: "List locates for a portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		locateDate, err := cmd.Flags().GetString("date")
		if err != nil {
			return err
		}

		locateIds, err := cmd.Flags().GetStringSlice("locate-ids")
		if err != nil {
			return err
		}

		request := &prime.ListExistingLocatesRequest{
			PortfolioId: portfolioId,
			LocateDate:  locateDate,
			LocateIds:   locateIds,
		}

		response, err := listLocates(svc, request)
		if err != nil {
			return err
		}

		if err := utils.PrintJsonDocs(cmd, response.Locates); err != nil {
			return err
		}

		return nil
	},
}

func listLocates(
	svc prime.FinancingService,
	req *prime.ListExistingLocatesRequest,
) (*prime.ListExistingLocatesResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListExistingLocates(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list locates: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listLocatesCmd)

	listLocatesCmd.Flags().StringSlice("locate-ids", []string{}, "The IDs of specific locates to filter for")
	listLocatesCmd.Flags().String("date", "", "The date of the locates in YYYY-MM-DD format")

	utils.AddPortfolioIdFlag(listLocatesCmd)
}
