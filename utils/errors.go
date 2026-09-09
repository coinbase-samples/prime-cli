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

	primeerrors "github.com/coinbase/prime-sdk-go/model/errors"
)

// FormatCLIError returns a user-facing error string. Prime API errors include
// trace_id and the spec subcode/code description via APIError.Format().
func FormatCLIError(err error) string {
	if err == nil {
		return ""
	}
	apiErr, ok := primeerrors.From(err)
	if !ok {
		return err.Error()
	}
	formatted := apiErr.Format()
	full := err.Error()
	if old := apiErr.Error(); strings.Contains(full, old) {
		return strings.Replace(full, old, formatted, 1)
	}
	return full + " (" + formatted + ")"
}
