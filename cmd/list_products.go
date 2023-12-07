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

var listProductsCmd = &cobra.Command{
	Use:   "list-products",
	Short: "List supported products",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		portfolioId := utils.GetFlagStringValue(cmd, utils.PortfolioIdFlag)
		if portfolioId == "" {
			portfolioId = client.Credentials.PortfolioId
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		request := &prime.ListProductsRequest{
			PortfolioId: portfolioId,
			Pagination:  pagination,
		}

		response, err := client.ListProducts(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list products: %w", err)
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
	rootCmd.AddCommand(listProductsCmd)

	listProductsCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listProductsCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listProductsCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listProductsCmd.Flags().BoolP(utils.FormatFlag, "", false, "Format the JSON output")
	listProductsCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
}
