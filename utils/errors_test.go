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
	"errors"
	"fmt"
	"strings"
	"testing"

	primeerrors "github.com/coinbase/prime-sdk-go/model/errors"
)

func TestFormatCliError_NonApiErrorUnchanged(t *testing.T) {
	err := errors.New("cannot unmarshal credentials")
	if got := FormatCliError(err); got != err.Error() {
		t.Fatalf("got %q, want %q", got, err.Error())
	}
}

func TestFormatCliError_Nil(t *testing.T) {
	if got := FormatCliError(nil); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func TestFormatCliError_WrappedApiErrorIncludesTraceAndDescription(t *testing.T) {
	apiErr := &primeerrors.ApiError{
		Response: primeerrors.Response{
			Code:    primeerrors.ErrorCodeValidationError,
			Message: "invalid entity",
			Subcode: primeerrors.SubcodeEntityIdInvalid,
			TraceID: "trace-abc-123",
		},
		StatusCode: 400,
		URL:        "/entities/x",
	}
	wrapped := fmt.Errorf("cannot get FCM balance: %w", apiErr)

	got := FormatCliError(wrapped)
	if !strings.Contains(got, "cannot get FCM balance:") {
		t.Fatalf("missing wrap prefix: %q", got)
	}
	if !strings.Contains(got, "trace_id=trace-abc-123") {
		t.Fatalf("missing trace_id: %q", got)
	}
	if !strings.Contains(got, "description=The entity_id is not a valid UUID.") {
		t.Fatalf("missing subcode description: %q", got)
	}
	if !strings.Contains(got, "subcode=ENTITY_ID_INVALID") {
		t.Fatalf("missing subcode: %q", got)
	}
}
