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
	"errors"
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/model"
	prime "github.com/coinbase/prime-sdk-go/webhooks"
	"github.com/spf13/cobra"
)

var updateSubscriptionCmd = &cobra.Command{
	Use:   "update-subscription",
	Short: "Updates a webhook subscription. Only the flags you set are changed",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, portfolioId, err := newService(cmd)
		if err != nil {
			return err
		}

		paths := updateMaskPaths(cmd.Flags().Changed)
		if len(paths) == 0 {
			return errors.New("no fields to update: set at least one update flag")
		}

		subscriptionId := utils.GetFlagStringValue(cmd, utils.SubscriptionIdFlag)

		// The update is optimistic-concurrency checked, so base it on the
		// current version unless the caller supplied one.
		expectedVersion := utils.GetFlagStringValue(cmd, utils.ExpectedVersionFlag)
		if expectedVersion == "" {
			ctx, cancel := utils.GetContextWithTimeout()
			current, err := svc.GetWebhookSubscription(ctx, &prime.GetWebhookSubscriptionRequest{
				PortfolioId:    portfolioId,
				SubscriptionId: subscriptionId,
			})
			cancel()
			if err != nil {
				return fmt.Errorf("cannot get webhook subscription: %w", err)
			}
			if current.Subscription == nil {
				return errors.New("cannot get webhook subscription: empty response")
			}
			expectedVersion = current.Subscription.Version
		}

		eventTypes, _ := cmd.Flags().GetStringSlice(utils.EventTypesFlag)

		rawHeaders, _ := cmd.Flags().GetStringArray(utils.CustomHeaderFlag)
		headers, err := parseCustomHeaders(rawHeaders)
		if err != nil {
			return err
		}

		request := &prime.UpdateWebhookSubscriptionRequest{
			PortfolioId:     portfolioId,
			SubscriptionId:  subscriptionId,
			EventTypes:      eventTypes,
			TargetUrl:       utils.GetFlagStringValue(cmd, utils.TargetUrlFlag),
			AuthType:        model.WebhookAuthType(withPrefix(authTypePrefix, utils.GetFlagStringValue(cmd, utils.AuthTypeFlag))),
			OauthConfig:     buildOauthConfig(cmd),
			CustomHeaders:   headers,
			Description:     utils.GetFlagStringValue(cmd, utils.DescriptionFlag),
			UpdateMask:      &model.WebhookUpdateMask{Paths: paths},
			ExpectedVersion: expectedVersion,
		}

		if cmd.Flags().Changed(utils.IsEnabledFlag) {
			isEnabled := utils.GetFlagBoolValue(cmd, utils.IsEnabledFlag)
			request.IsEnabled = &isEnabled
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		response, err := svc.UpdateWebhookSubscription(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot update webhook subscription: %w", err)
		}

		return printResponse(cmd, response)
	},
}

func init() {
	Cmd.AddCommand(updateSubscriptionCmd)

	utils.AddPortfolioIdFlag(updateSubscriptionCmd)
	addSubscriptionIdFlag(updateSubscriptionCmd)
	addOauthFlags(updateSubscriptionCmd)
	addCustomHeaderFlag(updateSubscriptionCmd)

	updateSubscriptionCmd.Flags().StringSlice(utils.EventTypesFlag, nil, "Replacement event types, comma separated")
	updateSubscriptionCmd.Flags().String(utils.TargetUrlFlag, "", "New target URL")
	updateSubscriptionCmd.Flags().String(utils.AuthTypeFlag, "", "New auth type")
	updateSubscriptionCmd.Flags().String(utils.DescriptionFlag, "", "New description")
	updateSubscriptionCmd.Flags().Bool(utils.IsEnabledFlag, true, "Enable or disable the subscription, e.g. --is-enabled=false")
	updateSubscriptionCmd.Flags().String(utils.ExpectedVersionFlag, "", "Subscription version the update is based on. Fetched automatically if blank")
}
