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
	"fmt"
	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go"

	"github.com/spf13/cobra"
)

var getAllocationCmd = &cobra.Command{
	Use:   "get-allocation",
	Short: "Get an allocation using an allocation ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId := utils.GetPortfolioId(cmd, client)

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.GetPortfolioAllocationRequest{
			PortfolioId:  portfolioId,
			AllocationId: utils.GetFlagStringValue(cmd, utils.AllocationIdFlag),
		}

		response, err := client.GetPortfolioAllocation(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get allocation: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJSON(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAllocationCmd)

	getAllocationCmd.Flags().StringP(utils.AllocationIdFlag, "i", "", "ID for allocation lookup (Required)")
	getAllocationCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getAllocationCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getAllocationCmd.MarkFlagRequired(utils.AllocationIdFlag)
}
