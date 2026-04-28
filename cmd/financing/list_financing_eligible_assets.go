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

package financing

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	prime "github.com/coinbase-samples/prime-sdk-go/financing"
	"github.com/spf13/cobra"
)

var listFinancingEligibleAssetsCmd = &cobra.Command{
	Use:   "list-financing-eligible-assets",
	Short: "Lists assets eligible for financing",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		request := &prime.ListFinancingEligibleAssetsRequest{}

		response, err := listFinancingEligibleAssets(svc, request)
		if err != nil {
			return err
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func listFinancingEligibleAssets(
	svc prime.FinancingService,
	req *prime.ListFinancingEligibleAssetsRequest,
) (*prime.ListFinancingEligibleAssetsResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.ListFinancingEligibleAssets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot list financing eligible assets: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(listFinancingEligibleAssetsCmd)
}
