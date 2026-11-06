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
	prime "github.com/coinbase/prime-sdk-go/webhooks"
	"github.com/spf13/cobra"
)

var replayEventCmd = &cobra.Command{
	Use:   "replay-event",
	Short: "Replays a webhook event",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.ReplayWebhookEvent(ctx, &prime.ReplayWebhookEventRequest{
			PortfolioId:    portfolioId,
			EventId:        utils.GetFlagStringValue(cmd, utils.EventIdFlag),
			SubscriptionId: utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag),
		})
		if err != nil {
			return fmt.Errorf("cannot replay webhook event: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(replayEventCmd)

	utils.AddPortfolioIdFlag(replayEventCmd)
	replayEventCmd.Flags().String(utils.EventIdFlag, "", "Webhook event ID (Required)")
	replayEventCmd.Flags().String(utils.SubscriptionIdFlag, "", "Limit the replay to a single subscription")
	replayEventCmd.MarkFlagRequired(utils.EventIdFlag)
}
