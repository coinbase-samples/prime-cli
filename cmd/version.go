/**
 * Copyright 2025-present Coinbase Global, Inc.
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

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"

	"github.com/spf13/cobra"
)

var primectlVersion = `{"version":"0.2.2"}`

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of primectl",
	RunE: func(cmd *cobra.Command, args []string) error {

		var doc map[string]interface{}

		if err := json.Unmarshal([]byte(primectlVersion), &doc); err != nil {
			return fmt.Errorf("cannot marshal version: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, doc)
		if err != nil {
			return fmt.Errorf("cannot create version response: %w", err)
		}

		fmt.Println(jsonResponse)
		return nil
	},
}

func init() {

	rootCmd.AddCommand(versionCmd)

}
