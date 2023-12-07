/* Copyright 2023-present Coinbase Global, Inc.
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
	"github.com/spf13/cobra"
)

var createAllocationCmd = &cobra.Command{
	Use:   "create-allocation",
	Short: "Create a portfolio allocation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		orderIds, err := cmd.Flags().GetStringArray(utils.OrderIdsFlag)
		if err != nil {
			return fmt.Errorf("cannot get order ids: %w", err)
		}

		allocationLegsJson, err := cmd.Flags().GetString(utils.AllocationLegsFlag)
		if err != nil {
			return fmt.Errorf("cannot get allocatio legs: %w", err)
		}

		var allocationLegs []*prime.AllocationLeg
		if err := json.Unmarshal([]byte(allocationLegsJson), &allocationLegs); err != nil {
			return fmt.Errorf("invalid allocation legs format: %w", err)
		}

		request := &prime.CreatePortfolioAllocationsRequest{
			AllocationId:                    utils.GetFlagStringValue(cmd, utils.AllocationIdFlag),
			SourcePortfolioId:               utils.GetFlagStringValue(cmd, utils.SourcePortfolioIdFlag),
			ProductId:                       utils.GetFlagStringValue(cmd, utils.ProductIdFlag),
			OrderIds:                        orderIds,
			AllocationLegs:                  allocationLegs,
			SizeType:                        utils.GetFlagStringValue(cmd, utils.SizeTypeFlag),
			RemainderDestinationPortfolioId: utils.GetFlagStringValue(cmd, utils.RemainderDestPortfolioIdFlag),
		}

		response, err := client.CreatePortfolioAllocations(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create portfolio allocations: %w", err)
		}

		jsonResponse, err := utils.MarshalJSON(response, cmd.Flags().Lookup(utils.FormatFlag).Changed)
		if err != nil {
			return fmt.Errorf("cannot marshal response to JSON: %w", err)
		}
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createAllocationCmd)

	createAllocationCmd.Flags().StringP(utils.AllocationIdFlag, "i", "", "ID of the allocation (Required)")
	createAllocationCmd.Flags().StringP(utils.SourcePortfolioIdFlag, "s", "", "ID of the source portfolio (Required)")
	createAllocationCmd.Flags().StringP(utils.ProductIdFlag, "p", "", "ID of the product (Required)")
	createAllocationCmd.Flags().StringP(utils.SizeTypeFlag, "t", "", "Size type of the allocation (Required)")
	createAllocationCmd.Flags().StringP(utils.RemainderDestPortfolioIdFlag, "r", "", "ID of the remainder destination portfolio (Required)")
	createAllocationCmd.Flags().StringP(utils.AllocationLegsFlag, "l", "", "JSON string of allocation legs (Required)")
	createAllocationCmd.Flags().StringArrayP(utils.OrderIdsFlag, "o", []string{}, "List of order IDs")
	createAllocationCmd.Flags().BoolP(utils.FormatFlag, "", false, "Format the JSON output")

	createAllocationCmd.MarkFlagRequired(utils.AllocationIdFlag)
	createAllocationCmd.MarkFlagRequired(utils.SourcePortfolioIdFlag)
	createAllocationCmd.MarkFlagRequired(utils.ProductIdFlag)
	createAllocationCmd.MarkFlagRequired(utils.SizeTypeFlag)
	createAllocationCmd.MarkFlagRequired(utils.RemainderDestPortfolioIdFlag)
	createAllocationCmd.MarkFlagRequired(utils.AllocationLegsFlag)
}
