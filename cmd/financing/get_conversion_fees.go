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
	prime "github.com/coinbase/prime-sdk-go/financing"
	"github.com/spf13/cobra"
)

var getConversionFeesCmd = &cobra.Command{
	Use:   "get-conversion-fees",
	Short: "Gets organization stablecoin conversion fee tiers",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := prime.NewFinancingService(client)

		request := &prime.GetConversionFeesRequest{}

		response, err := getConversionFees(svc, request)
		if err != nil {
			return err
		}

		if err := utils.PrintJsonDocs(cmd, response.Fees); err != nil {
			return err
		}

		return nil
	},
}

func getConversionFees(
	svc prime.FinancingService,
	req *prime.GetConversionFeesRequest,
) (*prime.GetConversionFeesResponse, error) {

	ctx, cancel := utils.GetContextWithTimeout()
	defer cancel()

	response, err := svc.GetConversionFees(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get conversion fees: %w", err)
	}

	return response, nil
}

func init() {
	Cmd.AddCommand(getConversionFeesCmd)
}
