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

package webhooks

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/model"
	prime "github.com/coinbase/prime-sdk-go/webhooks"
	"github.com/spf13/cobra"
)

var getDeliveryLogsCmd = &cobra.Command{
	Use:   "get-delivery-logs",
	Short: "Lists delivery log entries for a webhook subscription",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		request := &prime.GetWebhookDeliveryLogsRequest{
			PortfolioId:    portfolioId,
			SubscriptionId: utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag),
			StatusFilter: model.WebhookDeliveryStatus(
				withPrefix(deliveryStatusPrefix, utils.GetFlagStringValue(cmd, utils.StatusFilterFlag)),
			),
			EventTypeFilter: utils.GetFlagStringValue(cmd, utils.EventTypeFilterFlag),
			Cursor:          utils.GetFlagStringValue(cmd, utils.CursorFlag),
		}

		response, err := getDeliveryLogs(svc, request)
		if err != nil {
			return err
		}

		if !utils.GetFlagBoolValue(cmd, utils.AllFlag) {
			return printResponse(cmd, response)
		}

		var all []*model.WebhookDeliveryLogEntry
		for response != nil {
			all = append(all, response.Entries...)

			ctx, cancel := utils.GetContextWithTimeout()
			next, err := response.Next(ctx)
			cancel()
			if err != nil {
				return fmt.Errorf("cannot get webhook delivery logs: %w", err)
			}
			response = next
		}

		return utils.PrintJsonDocs(cmd, all)
	},
}

func getDeliveryLogs(
	svc prime.WebhooksService,
	req *prime.GetWebhookDeliveryLogsRequest,
) (*prime.GetWebhookDeliveryLogsResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.GetWebhookDeliveryLogs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get webhook delivery logs: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(getDeliveryLogsCmd)

	utils.AddPortfolioIdFlag(getDeliveryLogsCmd)
	addSubscriptionIdFlag(getDeliveryLogsCmd)
	addCursorFlags(getDeliveryLogsCmd)
	getDeliveryLogsCmd.Flags().String(utils.StatusFilterFlag, "", "Delivery status: DELIVERED, RETRYING, FAILED or PENDING")
	getDeliveryLogsCmd.Flags().String(utils.EventTypeFilterFlag, "", "Only return entries for this event type")
}
