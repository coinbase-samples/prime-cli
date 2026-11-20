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

package utils

import (
	"strings"

	"github.com/coinbase/prime-sdk-go/model"
)

// ParseWebhookAuthType accepts short or API enum values. An empty value lets
// the API choose the default on create and leaves auth unchanged on update
// unless explicitly masked.
func ParseWebhookAuthType(value string) model.WebhookAuthType {
	const prefix = "WEBHOOK_AUTH_TYPE_"
	value = strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(value)), prefix)
	switch value {
	case "":
		return ""
	case "HMAC":
		return model.WebhookAuthTypeSignatureOnly
	default:
		return model.WebhookAuthType(prefix + value)
	}
}
