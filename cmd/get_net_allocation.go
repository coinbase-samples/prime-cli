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

package cmd

import (
	"fmt"
	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/allocations"

	"github.com/spf13/cobra"
)

var getNetAllocationCmd = &cobra.Command{
	Use:   "get-net-allocation",
	Short: "Get a net allocation using a net allocation ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		allocationsService := allocations.NewAllocationsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &allocations.GetPortfolioNetAllocationRequest{
			PortfolioId: portfolioId,
			NettingId:   utils.GetFlagStringValue(cmd, utils.NettingIdFlag),
		}

		response, err := allocationsService.GetPortfolioNetAllocation(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get net allocation: %w", err)
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
	rootCmd.AddCommand(getNetAllocationCmd)

	getNetAllocationCmd.Flags().StringP(utils.AllocationIdFlag, "i", "", "ID for allocation lookup (Required)")
	getNetAllocationCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getNetAllocationCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getNetAllocationCmd.MarkFlagRequired(utils.AllocationIdFlag)
}
