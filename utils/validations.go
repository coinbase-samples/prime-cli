/**
 * Copyright 2023-present Coinbase Global, Inc.
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
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func ValidateUUID(uuid string) error {
	if uuid == "" {
		return errors.New("the UUID must not be empty")
	}
	r := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[4-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !r.MatchString(uuid) {
		return errors.New("the UUID is not valid")
	}
	return nil
}

func ValidateUUIDFlag(cmd *cobra.Command, flagName string) error {
	uuid, _ := cmd.Flags().GetString(flagName)
	if err := ValidateUUID(uuid); err != nil {
		return fmt.Errorf("%s must be a valid UUID: %w", flagName, err)
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ValidateSide(cmd *cobra.Command) error {
	side, _ := cmd.Flags().GetString(SideFlag)
	if side != OrderSideBuy && side != OrderSideSell {
		return errors.New("side must be either 'BUY' or 'SELL'")
	}
	return nil
}

func ValidateOrderTypeAndLimitPrice(cmd *cobra.Command) error {
	orderType, _ := cmd.Flags().GetString(TypeFlag)
	limitPrice, _ := cmd.Flags().GetString(LimitPriceFlag)
	switch strings.ToUpper(orderType) {
	case OrderTypeMarket:
		// No further validation needed for MARKET
	case OrderTypeLimit, OrderTypeTwap, OrderTypeVwap:
		if limitPrice == "" {
			return errors.New("limit-price is required for LIMIT, TWAP, and VWAP order types")
		}
	default:
		return errors.New("type must be one of MARKET, LIMIT, TWAP, or VWAP")
	}
	return nil
}

func ValidateTimeInForce(cmd *cobra.Command) error {
	timeInForce, _ := cmd.Flags().GetString(TimeInForceFlag)
	if timeInForce != "" {
		validOptions := []string{
			TifFillOrKill,
			TifGoodUntilDateTime,
			TifGoodUntilCancelled,
			TifImmediateOrCancel}
		if !contains(validOptions, timeInForce) {
			return fmt.Errorf("invalid time_in_force: %s. Must be one of: %v", timeInForce, validOptions)
		}
	}
	return nil
}

func ValidateQuantities(cmd *cobra.Command) error {
	baseQuantity, _ := cmd.Flags().GetString(BaseQuantityFlag)
	quoteValue, _ := cmd.Flags().GetString(QuoteValueFlag)
	if baseQuantity != "" && quoteValue != "" {
		return errors.New("either base-quantity or quote-value must be provided, not both")
	}
	if baseQuantity == "" && quoteValue == "" {
		return errors.New("one of base-quantity or quote-value must be provided")
	}
	return nil
}
