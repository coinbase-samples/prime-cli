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

var listSubscriptionsCmd = &cobra.Command{
	Use:   "list-subscriptions",
	Short: "Lists webhook subscriptions for a portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		includeDisabled, _ := cmd.Flags().GetBool(utils.IncludeDisabledFlag)

		request := &prime.ListWebhookSubscriptionsRequest{
			PortfolioId:     portfolioId,
			IncludeDisabled: includeDisabled,
			EventTypeFilter: utils.GetFlagStringValue(cmd, utils.EventTypeFilterFlag),
			Cursor:          utils.GetFlagStringValue(cmd, utils.CursorFlag),
		}

		response, err := listSubscriptions(svc, request)
		if err != nil {
			return err
		}

		if !utils.GetFlagBoolValue(cmd, utils.AllFlag) {
			return printResponse(cmd, response)
		}

		var all []*model.WebhookSubscription
		for response != nil {
			all = append(all, response.Subscriptions...)

			ctx, cancel := utils.GetContextWithTimeout()
			next, err := response.Next(ctx)
			cancel()
			if err != nil {
				return fmt.Errorf("cannot list webhook subscriptions: %w", err)
			}
			response = next
		}

		return utils.PrintJsonDocs(cmd, all)
	},
}

func listSubscriptions(
	svc prime.WebhooksService,
	req *prime.ListWebhookSubscriptionsRequest,
) (*prime.ListWebhookSubscriptionsResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListWebhookSubscriptions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list webhook subscriptions: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listSubscriptionsCmd)

	utils.AddPortfolioIdFlag(listSubscriptionsCmd)
	addCursorFlags(listSubscriptionsCmd)
	listSubscriptionsCmd.Flags().Bool(utils.IncludeDisabledFlag, false, "Include disabled subscriptions")
	listSubscriptionsCmd.Flags().String(utils.EventTypeFilterFlag, "", "Only return subscriptions for this event type")
}
