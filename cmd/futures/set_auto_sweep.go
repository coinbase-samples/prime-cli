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

package futures

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/futures"

	"github.com/spf13/cobra"
)

var setAutoSweepCmd = &cobra.Command{
	Use:   "set-auto-sweep",
	Short: "Enables or disables auto sweep for an entity's futures account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		svc := futures.NewFuturesService(client)

		entityId, err := utils.GetEntityId(cmd, client)
		if err != nil {
			return fmt.Errorf("cannot get entity ID: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &futures.SetAutoSweepRequest{
			EntityId:  entityId,
			AutoSweep: utils.GetFlagBoolValue(cmd, utils.AutoSweepEnabledFlag),
		}

		response, err := svc.SetAutoSweep(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot set auto sweep: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func init() {
	Cmd.AddCommand(setAutoSweepCmd)

	utils.AddEntityIdFlag(setAutoSweepCmd)
	setAutoSweepCmd.Flags().Bool(utils.AutoSweepEnabledFlag, false, "Whether auto sweep is enabled (Required)")
	setAutoSweepCmd.MarkFlagRequired(utils.AutoSweepEnabledFlag)
}
