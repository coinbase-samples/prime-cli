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
	"log"

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

		pagination, err := utils.GetPaginationParams(cmd)
		if err != nil {
			return fmt.Errorf("cannot get pagination params: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.ListInvoicesRequest{
			EntityId:     client.Credentials.EntityId,
			States:       states,
			BillingYear:  billingYear,
			BillingMonth: billingMonth,
			Pagination:   pagination,
		}

		log.Printf("Sending request: %+v\n", request)

		response, err := client.ListInvoices(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list invoices: %w", err)
		}

		log.Printf("Received response: %+v\n", response)

		jsonResponse, err := utils.FormatResponseAsJSON(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listInvoicesCmd)

	listInvoicesCmd.Flags().StringP(utils.CursorFlag, "c", "", "Pagination cursor")
	listInvoicesCmd.Flags().StringP(utils.LimitFlag, "l", utils.LimitDefault, "Pagination limit")
	listInvoicesCmd.Flags().StringP(utils.SortDirectionFlag, "d", utils.SortDirectionDefault, "Sort direction")
	listInvoicesCmd.Flags().StringSliceP(utils.InvoiceStatesFlag, "s", []string{}, "List of states")
	listInvoicesCmd.Flags().IntSliceP(utils.InvoiceBillingYear, "y", []int{}, "Billing year")
	listInvoicesCmd.Flags().IntSliceP(utils.InvoiceBillingMonth, "m", []int{}, "Billing month")
	listInvoicesCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
}
