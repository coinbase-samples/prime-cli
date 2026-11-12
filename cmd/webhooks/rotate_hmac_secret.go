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

var rotateHmacSecretCmd = &cobra.Command{
	Use:   "rotate-hmac-secret",
	Short: "Rotates a subscription's HMAC secret (the new secret is returned only once)",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.RotateWebhookHmacSecret(ctx, &prime.RotateWebhookHmacSecretRequest{
			PortfolioId:               portfolioId,
			SubscriptionId:            utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag),
			RevokePreviousImmediately: utils.GetFlagBoolValue(cmd, utils.RevokePreviousImmediatelyFlag),
		})
		if err != nil {
			return fmt.Errorf("cannot rotate webhook HMAC secret: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(rotateHmacSecretCmd)

	utils.AddPortfolioIdFlag(rotateHmacSecretCmd)
	addSubscriptionIdFlag(rotateHmacSecretCmd)
	rotateHmacSecretCmd.Flags().Bool(
		utils.RevokePreviousImmediatelyFlag,
		false,
		"Invalidate the previous secret immediately instead of after the grace window",
	)
}
