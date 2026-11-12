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
	"reflect"
	"testing"

	"github.com/coinbase-samples/prime-cli/utils"
)

func TestUpdateMaskPaths(t *testing.T) {
	tests := []struct {
		name    string
		changed []string
		want    []string
	}{
		{"none", nil, nil},
		{"target url", []string{utils.TargetUrlFlag}, []string{"target_url"}},
		{
			"stable order regardless of input",
			[]string{utils.IsEnabledFlag, utils.DescriptionFlag, utils.EventTypesFlag},
			[]string{"event_types", "description", "is_enabled"},
		},
		{
			"any oauth flag maps to a single oauth_config path",
			[]string{utils.OauthClientIdFlag, utils.OauthScopesFlag},
			[]string{"oauth_config"},
		},
		{"custom headers", []string{utils.CustomHeaderFlag}, []string{"custom_headers"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := map[string]bool{}
			for _, name := range tt.changed {
				set[name] = true
			}

			got := updateMaskPaths(func(name string) bool { return set[name] })
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("updateMaskPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCustomHeaders(t *testing.T) {
	headers, err := parseCustomHeaders([]string{"X-Env=prod", "X-Token=a=b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if headers["X-Env"] != "prod" || headers["X-Token"] != "a=b" {
		t.Fatalf("unexpected headers: %v", headers)
	}

	if headers, err := parseCustomHeaders(nil); err != nil || headers != nil {
		t.Fatalf("expected nil headers and nil error, got %v, %v", headers, err)
	}

	for _, bad := range []string{"novalue", "=value", "  =value"} {
		if _, err := parseCustomHeaders([]string{bad}); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
}

func TestWithPrefix(t *testing.T) {
	tests := map[string]string{
		"7d":                       "WEBHOOK_STATS_PERIOD_7D",
		"WEBHOOK_STATS_PERIOD_24H": "WEBHOOK_STATS_PERIOD_24H",
		" 30d ":                    "WEBHOOK_STATS_PERIOD_30D",
		"":                         "",
	}

	for in, want := range tests {
		if got := withPrefix(statsPeriodPrefix, in); got != want {
			t.Errorf("withPrefix(%q) = %q, want %q", in, got, want)
		}
	}
}
