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

var sendSimulationCmd = &cobra.Command{
	Use:   "send-simulation",
	Short: "Sends a simulated event to a webhook subscription",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.SendWebhookSimulation(ctx, &prime.SendWebhookSimulationRequest{
			PortfolioId:    portfolioId,
			SubscriptionId: utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag),
			EventType:      utils.GetFlagStringValue(cmd, utils.EventTypeFlag),
		})
		if err != nil {
			return fmt.Errorf("cannot send webhook simulation: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(sendSimulationCmd)

	utils.AddPortfolioIdFlag(sendSimulationCmd)
	addSubscriptionIdFlag(sendSimulationCmd)
	sendSimulationCmd.Flags().String(utils.EventTypeFlag, "", "Event type to simulate, from get-event-types (Required)")
	sendSimulationCmd.MarkFlagRequired(utils.EventTypeFlag)
}
