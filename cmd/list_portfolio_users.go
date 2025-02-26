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
	"github.com/coinbase-samples/prime-sdk-go/users"

	"github.com/spf13/cobra"
)

var listPortfolioUsersCmd = &cobra.Command{
	Use:   "list-portfolio-users",
	Short: "List users associated with portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		usersService := users.NewUsersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &users.ListPortfolioUsersRequest{
			PortfolioId: portfolioId,
			Pagination:  pagination,
		}

		response, err := usersService.ListPortfolioUsers(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list users: %w", err)
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
	rootCmd.AddCommand(listPortfolioUsersCmd)

	listPortfolioUsersCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listPortfolioUsersCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listPortfolioUsersCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listPortfolioUsersCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	listPortfolioUsersCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

}
