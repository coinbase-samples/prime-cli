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

package invoices

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/invoice"
	"github.com/coinbase-samples/prime-sdk-go/model"

	"github.com/spf13/cobra"
)

var listInvoicesCmd = &cobra.Command{
	Use:   "list-invoices",
	Short: "List invoices matching filter criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := invoice.NewInvoiceService(client)

		states, err := cmd.Flags().GetStringSlice(utils.InvoiceStatesFlag)
		if err != nil {
			return fmt.Errorf("cannot get states slice: %w", err)
		}

		billingYearSlice, err := cmd.Flags().GetIntSlice(utils.InvoiceBillingYear)
		if err != nil {
			return fmt.Errorf("cannot get year slice: %w", err)
		}

		billingMonthSlice, err := cmd.Flags().GetIntSlice(utils.InvoiceBillingMonth)
		if err != nil {
			return fmt.Errorf("cannot get month slice: %w", err)
		}

		var billingYear, billingMonth int32
		if len(billingYearSlice) > 0 {
			billingYear = int32(billingYearSlice[0])
		}

		if len(billingMonthSlice) > 0 {
			billingMonth = int32(billingMonthSlice[0])
		}

		return utils.HandleListCmd(
			cmd,
			func(paginationParams *model.PaginationParams) (*model.Pagination, error) {
				response, err := listInvoices(svc, client.Credentials().EntityId, states, billingYear, billingMonth, paginationParams)
				if err != nil {
					return nil, err
				}

				if err := utils.PrintJsonDocs(cmd, response.Invoices); err != nil {
					return nil, err
				}

				return response.Pagination, nil
			},
		)
	},
}

func listInvoices(
	svc invoice.InvoiceService,
	entityId string,
	states []string,
	billingYear,
	billingMonth int32,
	pagination *model.PaginationParams,
) (*invoice.ListInvoicesResponse, error) {
	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	request := &invoice.ListInvoicesRequest{
		EntityId:     entityId,
		States:       states,
		BillingYear:  billingYear,
		BillingMonth: billingMonth,
		Pagination:   pagination,
	}

	response, err := svc.ListInvoices(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cannot list invoices: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listInvoicesCmd)

	listInvoicesCmd.Flags().StringSliceP(utils.InvoiceStatesFlag, "", []string{}, "List of states")
	listInvoicesCmd.Flags().IntSliceP(utils.InvoiceBillingYear, "", []int{}, "Billing year")
	listInvoicesCmd.Flags().IntSliceP(utils.InvoiceBillingMonth, "", []int{}, "Billing month")

	utils.AddPaginationFlags(listInvoicesCmd, true)

}
