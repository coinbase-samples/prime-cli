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

package advancedtransfers

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/advancedtransfers"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/spf13/cobra"
)

var listAdvancedTransfersCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists advanced transfers for a portfolio meeting filter criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := advancedtransfers.NewAdvancedTransfersService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		states, err := cmd.Flags().GetStringSlice(utils.StatesFlag)
		if err != nil {
			return fmt.Errorf("cannot get states slice: %w", err)
		}

		var state model.AdvancedTransferState
		if len(states) > 0 {
			state = model.AdvancedTransferState(states[0])
		}

		transferType := model.AdvancedTransferType(
			utils.GetFlagStringValue(cmd, utils.TransferTypeFlag),
		)

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listAdvancedTransfers(svc, portfolioId, state, transferType, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.AdvancedTransfers); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listAdvancedTransfers(
	svc advancedtransfers.AdvancedTransfersService,
	portfolioId string,
	state model.AdvancedTransferState,
	transferType model.AdvancedTransferType,
	pagination *model.PaginationParams,
) (*advancedtransfers.ListAdvancedTransfersResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &advancedtransfers.ListAdvancedTransfersRequest{
		PortfolioId: portfolioId,
		State:       state,
		Type:        transferType,
		Pagination:  pagination,
	}

	response, err := svc.ListAdvancedTransfers(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list advanced transfers: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listAdvancedTransfersCmd)

	utils.AddPortfolioIdFlag(listAdvancedTransfersCmd)
	utils.AddPaginationFlags(listAdvancedTransfersCmd, true)

	listAdvancedTransfersCmd.Flags().StringSlice(utils.StatesFlag, []string{}, "List of advanced transfer states")
	listAdvancedTransfersCmd.Flags().String(utils.TransferTypeFlag, "", "Optional transfer type filter")
}
