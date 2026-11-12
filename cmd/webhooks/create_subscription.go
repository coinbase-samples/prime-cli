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

var createSubscriptionCmd = &cobra.Command{
	Use:   "create-subscription",
	Short: "Creates a webhook subscription (the HMAC secret is returned only once)",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		eventTypes, _ := cmd.Flags().GetStringSlice(utils.EventTypesFlag)

		rawHeaders, _ := cmd.Flags().GetStringArray(utils.CustomHeaderFlag)
		headers, err := parseCustomHeaders(rawHeaders)
		if err != nil {
			return err
		}

		request := &prime.CreateWebhookSubscriptionRequest{
			PortfolioId:    portfolioId,
			EventTypes:     eventTypes,
			TargetUrl:      utils.GetFlagStringValue(cmd, utils.TargetUrlFlag),
			AuthType:       model.WebhookAuthType(withPrefix(authTypePrefix, utils.GetFlagStringValue(cmd, utils.AuthTypeFlag))),
			OauthConfig:    buildOauthConfig(cmd),
			CustomHeaders:  headers,
			Description:    utils.GetFlagStringValue(cmd, utils.DescriptionFlag),
			IdempotencyKey: utils.GetFlagStringValue(cmd, utils.IdempotencyKeyFlag),
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.CreateWebhookSubscription(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create webhook subscription: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(createSubscriptionCmd)

	utils.AddPortfolioIdFlag(createSubscriptionCmd)
	utils.AddIdempotencyKeyFlag(createSubscriptionCmd)
	addOauthFlags(createSubscriptionCmd)
	addCustomHeaderFlag(createSubscriptionCmd)

	createSubscriptionCmd.Flags().StringSlice(utils.EventTypesFlag, nil, "Event types to subscribe to, comma separated (Required)")
	createSubscriptionCmd.Flags().String(utils.TargetUrlFlag, "", "HTTPS URL that receives deliveries (Required)")
	createSubscriptionCmd.Flags().String(utils.AuthTypeFlag, "", "Auth type: HMAC, OAUTH, MTLS or OAUTH_MTLS (the WEBHOOK_AUTH_TYPE_ prefix is optional)")
	createSubscriptionCmd.Flags().String(utils.DescriptionFlag, "", "Subscription description")

	createSubscriptionCmd.MarkFlagRequired(utils.EventTypesFlag)
	createSubscriptionCmd.MarkFlagRequired(utils.TargetUrlFlag)
}
