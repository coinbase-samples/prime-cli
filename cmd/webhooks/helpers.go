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
	"strings"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/model"
	prime "github.com/coinbase/prime-sdk-go/webhooks"
	"github.com/spf13/cobra"
)

// newService builds the webhooks service and resolves the portfolio ID from
// the --portfolio-id flag or the credentials.
func newService(cmd *cobra.Command) (prime.WebhooksService, string, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return nil, "", fmt.Errorf("failed to initialize client: %w", err)
	}

	portfolioId, err := utils.GetPortfolioId(cmd, client)
	if err != nil {
		return nil, "", err
	}

	return prime.NewWebhooksService(client), portfolioId, nil
}

func printResponse(cmd *cobra.Command, response interface{}) error {
	jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
	if err != nil {
		return err
	}

	fmt.Println(jsonResponse)

	return nil
}

// withPrefix upper-cases an enum flag value and adds the API prefix when the
// caller used the short form, e.g. "7d" -> "WEBHOOK_STATS_PERIOD_7D".
func withPrefix(prefix, value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" || strings.HasPrefix(value, prefix) {
		return value
	}
	return prefix + value
}

const (
	authTypePrefix       = "WEBHOOK_AUTH_TYPE_"
	deliveryStatusPrefix = "WEBHOOK_DELIVERY_STATUS_"
	statsPeriodPrefix    = "WEBHOOK_STATS_PERIOD_"
)

func addSubscriptionIdFlag(cmd *cobra.Command) {
	cmd.Flags().String(utils.SubscriptionIdFlag, "", "Webhook subscription ID (Required)")
	cmd.MarkFlagRequired(utils.SubscriptionIdFlag)
}

func addCursorFlags(cmd *cobra.Command) {
	cmd.Flags().String(utils.CursorFlag, "", "Cursor from a previous response's next_cursor")
	cmd.Flags().Bool(utils.AllFlag, false, "Set to print all results without manually paging through results")
}

func addOauthFlags(cmd *cobra.Command) {
	cmd.Flags().String(utils.OauthTokenUrlFlag, "", "OAuth token URL")
	cmd.Flags().String(utils.OauthClientIdFlag, "", "OAuth client ID")
	cmd.Flags().String(utils.OauthClientSecretFlag, "", "OAuth client secret")
	cmd.Flags().StringSlice(utils.OauthScopesFlag, nil, "OAuth scopes, comma separated")
}

func addCustomHeaderFlag(cmd *cobra.Command) {
	cmd.Flags().StringArray(utils.CustomHeaderFlag, nil, "Custom delivery header as key=value. Repeat for multiple headers")
}

func oauthFlagsChanged(cmd *cobra.Command) bool {
	for _, name := range []string{
		utils.OauthTokenUrlFlag,
		utils.OauthClientIdFlag,
		utils.OauthClientSecretFlag,
		utils.OauthScopesFlag,
	} {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

// buildOauthConfig returns nil when no OAuth flag was set.
func buildOauthConfig(cmd *cobra.Command) *model.WebhookOAuthConfigInput {
	if !oauthFlagsChanged(cmd) {
		return nil
	}

	scopes, _ := cmd.Flags().GetStringSlice(utils.OauthScopesFlag)

	return &model.WebhookOAuthConfigInput{
		TokenUrl:     utils.GetFlagStringValue(cmd, utils.OauthTokenUrlFlag),
		ClientId:     utils.GetFlagStringValue(cmd, utils.OauthClientIdFlag),
		ClientSecret: utils.GetFlagStringValue(cmd, utils.OauthClientSecretFlag),
		Scopes:       scopes,
	}
}

// parseCustomHeaders converts key=value pairs into a header map. It returns
// nil when no headers were provided.
func parseCustomHeaders(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}

	headers := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid --%s value %q: expected key=value", utils.CustomHeaderFlag, pair)
		}
		headers[key] = value
	}

	return headers, nil
}

// updateMaskPaths returns the snake_case field paths for the update flags
// that were explicitly set, in a stable order.
func updateMaskPaths(changed func(name string) bool) []string {
	fields := []struct {
		flag []string
		path string
	}{
		{[]string{utils.EventTypesFlag}, "event_types"},
		{[]string{utils.TargetUrlFlag}, "target_url"},
		{[]string{utils.AuthTypeFlag}, "auth_type"},
		{[]string{
			utils.OauthTokenUrlFlag,
			utils.OauthClientIdFlag,
			utils.OauthClientSecretFlag,
			utils.OauthScopesFlag,
		}, "oauth_config"},
		{[]string{utils.CustomHeaderFlag}, "custom_headers"},
		{[]string{utils.DescriptionFlag}, "description"},
		{[]string{utils.IsEnabledFlag}, "is_enabled"},
	}

	var paths []string
	for _, f := range fields {
		for _, name := range f.flag {
			if changed(name) {
				paths = append(paths, f.path)
				break
			}
		}
	}

	return paths
}
