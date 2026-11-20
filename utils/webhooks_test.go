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
	"encoding/json"
	"testing"

	"github.com/coinbase/prime-sdk-go/webhooks"
)

func TestWebhookAuthTypeRequestJSON(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"", ""},
		{" signature_only ", "WEBHOOK_AUTH_TYPE_SIGNATURE_ONLY"},
		{"HMAC", "WEBHOOK_AUTH_TYPE_SIGNATURE_ONLY"},
		{"WEBHOOK_AUTH_TYPE_HMAC", "WEBHOOK_AUTH_TYPE_SIGNATURE_ONLY"},
		{"oauth", "WEBHOOK_AUTH_TYPE_OAUTH"},
		{"WEBHOOK_AUTH_TYPE_MTLS", "WEBHOOK_AUTH_TYPE_MTLS"},
		{"OAUTH_MTLS", "WEBHOOK_AUTH_TYPE_OAUTH_MTLS"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			authType := ParseWebhookAuthType(tc.input)
			for _, req := range []any{
				webhooks.CreateWebhookSubscriptionRequest{AuthType: authType},
				webhooks.UpdateWebhookSubscriptionRequest{AuthType: authType},
			} {
				body, err := json.Marshal(req)
				if err != nil {
					t.Fatal(err)
				}
				var fields map[string]json.RawMessage
				if err := json.Unmarshal(body, &fields); err != nil {
					t.Fatal(err)
				}
				want := ""
				if tc.want != "" {
					want = `"` + tc.want + `"`
				}
				if got := string(fields["auth_type"]); got != want {
					t.Errorf("%T auth_type = %s, want %s", req, got, want)
				}
			}
		})
	}
}
