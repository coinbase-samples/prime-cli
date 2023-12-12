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
	"context"
	"encoding/json"
	"fmt"
	"github.com/coinbase-samples/prime-sdk-go"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getDefaultTimeoutDuration() time.Duration {
	envTimeout := os.Getenv("primeCliTimeout")
	if envTimeout != "" {
		if value, err := strconv.Atoi(envTimeout); err == nil && value > 0 {
			return time.Duration(value) * time.Second
		}
	}
	return 7 * time.Second
}

func GetContextWithTimeout() (context.Context, context.CancelFunc) {
	timeoutDuration := getDefaultTimeoutDuration()
	return context.WithTimeout(context.Background(), timeoutDuration)
}

func GetClientFromEnv() (*prime.Client, error) {
	credentials := &prime.Credentials{}
	if err := json.Unmarshal([]byte(os.Getenv("PRIME_CREDENTIALS")), credentials); err != nil {
		return nil, fmt.Errorf("cannot unmarshal credentials: %w", err)
	}

	client := prime.NewClient(credentials, http.Client{})
	return client, nil
}

func GetFlagStringValue(cmd *cobra.Command, flagName string) string {
	value, _ := cmd.Flags().GetString(flagName)
	return value
}

func GetPaginationParams(cmd *cobra.Command) (*prime.PaginationParams, error) {
	cursor, err := cmd.Flags().GetString("cursor")
	if err != nil {
		return nil, fmt.Errorf("cannot parse cursor: %w", err)
	}
	limit, err := cmd.Flags().GetString("limit")
	if err != nil {
		return nil, fmt.Errorf("cannot parse limit: %w", err)
	}

	sortDirection, err := cmd.Flags().GetString("sort-direction")
	if err != nil {
		return nil, fmt.Errorf("cannot parse sort direction: %w", err)
	}

	return &prime.PaginationParams{
		Cursor:        cursor,
		Limit:         limit,
		SortDirection: sortDirection,
	}, nil
}

func ParseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	var start, end time.Time
	var err error
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return start, end, fmt.Errorf("invalid start time: %w", err)
		}
	}
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return start, end, fmt.Errorf("invalid end time: %w", err)
		}
	}
	return start, end, nil
}

func MarshalJSON(data interface{}, format bool) ([]byte, error) {
	if format {
		return json.MarshalIndent(data, "", JsonIndent)
	}
	return json.Marshal(data)
}

func CheckFormatFlag(cmd *cobra.Command) (bool, error) {
	formatFlagValue, err := cmd.Flags().GetString(FormatFlag)
	if err != nil {
		return false, fmt.Errorf("cannot read format flag: %w", err)
	}
	return formatFlagValue == "true", nil
}

func GetPortfolioId(cmd *cobra.Command, client *prime.Client) string {
	portfolioId := GetFlagStringValue(cmd, PortfolioIdFlag)
	if portfolioId == "" {
		portfolioId = client.Credentials.PortfolioId
	}
	return portfolioId
}

func FormatResponseAsJSON(cmd *cobra.Command, response interface{}) (string, error) {
	shouldFormat, err := CheckFormatFlag(cmd)
	if err != nil {
		return "", err
	}

	jsonResponse, err := MarshalJSON(response, shouldFormat)
	if err != nil {
		return "", fmt.Errorf("cannot marshal response to JSON: %w", err)
	}

	return string(jsonResponse), nil
}
