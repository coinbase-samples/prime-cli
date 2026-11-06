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

var getSubscriptionStatsCmd = &cobra.Command{
	Use:   "get-subscription-stats",
	Short: "Gets delivery statistics for a webhook subscription",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.GetWebhookSubscriptionStats(ctx, &prime.GetWebhookSubscriptionStatsRequest{
			PortfolioId:    portfolioId,
			SubscriptionId: utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag),
			Period: model.WebhookStatsPeriod(
				withPrefix(statsPeriodPrefix, utils.GetFlagStringValue(cmd, utils.PeriodFlag)),
			),
		})
		if err != nil {
			return fmt.Errorf("cannot get webhook subscription stats: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(getSubscriptionStatsCmd)

	utils.AddPortfolioIdFlag(getSubscriptionStatsCmd)
	addSubscriptionIdFlag(getSubscriptionStatsCmd)
	getSubscriptionStatsCmd.Flags().String(utils.PeriodFlag, "", "Stats period: 24H, 7D or 30D")
}
