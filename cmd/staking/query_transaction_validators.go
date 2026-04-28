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

package staking

import (
	"fmt"
	"strconv"

	"github.com/coinbase-samples/prime-cli/utils"
	primeStaking "github.com/coinbase-samples/prime-sdk-go/staking"
	"github.com/spf13/cobra"
)

const transactionIdsFlag = "transaction-ids"

var queryTransactionValidatorsCmd = &cobra.Command{
	Use:   "query-validators",
	Short: "Queries transaction validators for a portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := primeStaking.NewStakingService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		transactionIds, err := cmd.Flags().GetStringSlice(transactionIdsFlag)
		if err != nil {
			return err
		}

		request := &primeStaking.QueryTransactionValidatorsRequest{
			PortfolioId:    portfolioId,
			TransactionIds: transactionIds,
			Cursor:         utils.GetFlagStringValue(cmd, "cursor"),
			SortDirection:  utils.GetFlagStringValue(cmd, utils.SortDirectionFlag),
		}

		if limitStr := utils.GetFlagStringValue(cmd, utils.LimitFlag); limitStr != "" {
			parsed, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				return fmt.Errorf("cannot parse limit: %w", err)
			}
			request.Limit = int32(parsed)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.QueryTransactionValidators(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot query transaction validators: %w", err)
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
	Cmd.AddCommand(queryTransactionValidatorsCmd)
	utils.AddPortfolioIdFlag(queryTransactionValidatorsCmd)

	queryTransactionValidatorsCmd.Flags().StringSlice(transactionIdsFlag, []string{}, "List of transaction IDs to query")
	queryTransactionValidatorsCmd.Flags().String(utils.LimitFlag, "", "Maximum number of results to return")
	queryTransactionValidatorsCmd.Flags().String(utils.SortDirectionFlag, "", "Sort direction (ASC or DESC)")
	queryTransactionValidatorsCmd.Flags().String("cursor", "", "Pagination cursor")
}
