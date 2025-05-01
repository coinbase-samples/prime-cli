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
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/products"

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

		svc := products.NewProductsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listProducts(svc, portfolioId, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Products); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			})
	},
}

func listProducts(
	svc products.ProductsService,
	portfolioId string,
	pagination *model.PaginationParams,
) (*products.ListProductsResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &products.ListProductsRequest{
		PortfolioId: portfolioId,
		Pagination:  pagination,
	}

	response, err := svc.ListProducts(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list products: %w", err)
	}

	return response, nil
}

func init() {
	rootCmd.AddCommand(listProductsCmd)
	utils.AddPortfolioIdFlag(listProductsCmd)
	utils.AddPaginationFlags(listProductsCmd, true)
}
