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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coinbase-samples/prime-sdk-go/client"
	"github.com/coinbase-samples/prime-sdk-go/credentials"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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

func GetClientFromEnv() (client.RestClient, error) {
	creds := &credentials.Credentials{}
	if err := json.Unmarshal([]byte(os.Getenv("PRIME_CREDENTIALS")), creds); err != nil {
		return nil, fmt.Errorf("cannot unmarshal credentials: %w", err)
	}

	restClient := client.NewRestClient(creds, http.Client{})
	return restClient, nil
}

func GetFlagStringValue(cmd *cobra.Command, flagName string) string {
	value, _ := cmd.Flags().GetString(flagName)
	return value
}

func PrintJsonDocs[T any](cmd *cobra.Command, items []T) error {

	for _, item := range items {
		docStr, err := FormatResponseAsJson(cmd, item)
		if err != nil {
			return err
		}
		fmt.Println(docStr)
	}

	return nil
}

func AddPortfolioIdFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
}

func AddStartEndFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(StartFlag, "", "", "Start time in RFC3339 format")
	cmd.Flags().StringP(EndFlag, "", "", "End time in RFC3339 format")
}

func AddPaginationFlags(cmd *cobra.Command, includeSortLimit bool) {

	if includeSortLimit {
		cmd.Flags().StringP(LimitFlag, "l", LimitDefault, "Pagination limit")
		cmd.Flags().StringP(SortDirectionFlag, "d", SortDirectionDefault, "Sort direction")
	}

	cmd.Flags().BoolP(AllFlag, "", false, "Set to print all results without manually paging through results")
	cmd.Flags().BoolP(InteractiveFlag, "", false, "Iterate through all results by manually paging through results")
}

func GetPaginationParams(cmd *cobra.Command) (*model.PaginationParams, error) {
	limitStr, err := cmd.Flags().GetString("limit")
	if err != nil {
		return nil, fmt.Errorf("cannot parse limit: %w", err)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid limit value: %w", err)
	}

	sortDirection, err := cmd.Flags().GetString("sort-direction")
	if err != nil {
		return nil, fmt.Errorf("cannot parse sort direction: %w", err)
	}

	return &model.PaginationParams{
		Cursor:        "",
		Limit:         int32(limit),
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
	formatFlagValue, err := cmd.Flags().GetBool(FormatFlag)
	if err != nil {
		return false, fmt.Errorf("cannot read format flag: %w", err)
	}
	return formatFlagValue, nil
}

func isAllFlagSet(cmd *cobra.Command) (bool, error) {
	return isBoolFlagSet(cmd, AllFlag)
}

func isInteractiveFlagSet(cmd *cobra.Command) (bool, error) {
	return isBoolFlagSet(cmd, InteractiveFlag)
}

func isBoolFlagSet(cmd *cobra.Command, name string) (bool, error) {
	flag, err := cmd.Flags().GetBool(name)
	if err != nil {
		return false, fmt.Errorf("cannot read %s flag: %w", name, err)
	}
	return flag, nil
}

func printInteractivePrompt() {
	fmt.Print("Press space to continue, q to quit: ")
}

type ListCmdCallback func(paginationParams *model.PaginationParams) (*model.Pagination, error)

func HandleListCmd(cmd *cobra.Command, callback ListCmdCallback) error {

	paginationParams, err := GetPaginationParams(cmd)
	if err != nil {
		return err
	}

	all, err := isAllFlagSet(cmd)
	if err != nil {
		return err
	}

	interactive, err := isInteractiveFlagSet(cmd)
	if err != nil {
		return err
	}

	var nextCursor string
	for {

		paginationParams.Cursor = nextCursor

		pagination, err := callback(paginationParams)
		if err != nil {
			return err
		}

		shouldContinue, shouldBreak, cursor, err := continueBreakInteractive(
			all,
			interactive,
			pagination,
		)
		if err != nil {
			return err
		}

		nextCursor = cursor

		if shouldBreak {
			break
		}

		if shouldContinue {
			continue
		}
	}

	return nil
}

func continueBreakInteractive(
	all,
	interactive bool,
	pagination *model.Pagination,
) (shouldContinue bool, shouldBreak bool, nextCursor string, err error) {

	if !all && !interactive {
		shouldBreak = true
		return
	}

	nextCursor = pagination.NextCursor

	if len(nextCursor) == 0 {
		shouldBreak = true
		return
	}

	if !interactive {
		shouldContinue = true
		return
	}

	var quit bool

	if quit, err = interactiveCheck(); err != nil {
		return
	} else if quit {
		shouldBreak = true
		return
	}

	printInteractiveNewline()
	return
}

func printInteractiveNewline() {
	fmt.Print("\r\n")
}

func isInteractiveQuit(input byte) bool {
	return input == 'q'
}

func interactiveCheck() (bool, error) {

	printInteractivePrompt()

	input, err := readRawStateByte()
	if err != nil {
		return false, err
	}

	return isInteractiveQuit(input), nil
}

func readRawStateByte() (input byte, err error) {
	rawState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), rawState)
	input, err = bufio.NewReader(os.Stdin).ReadByte()
	if err != nil {
		err = fmt.Errorf("unable to read from terminal: %w", err)
	}
	return
}

func GetPortfolioId(cmd *cobra.Command, client client.RestClient) (string, error) {
	portfolioId, err := cmd.Flags().GetString(PortfolioIdFlag)
	if err != nil {
		return "", fmt.Errorf("error retrieving portfolio ID: %w", err)
	}

	if portfolioId == "" {
		creds := client.Credentials()
		if creds == nil {
			return "", errors.New("client credentials are nil")
		}
		portfolioId = creds.PortfolioId
		if portfolioId == "" {
			return "", errors.New("portfolio ID is not provided in both flag and client credentials")
		}
	}

	return portfolioId, nil
}

func GetEntityId(cmd *cobra.Command, client client.RestClient) (string, error) {
	entityId, err := cmd.Flags().GetString(EntityIdFlag)
	if err != nil {
		return "", fmt.Errorf("error retrieving entity ID: %w", err)
	}

	if entityId == "" {
		creds := client.Credentials()
		if creds == nil {
			return "", errors.New("client credentials are nil")
		}
		entityId = creds.EntityId
		if entityId == "" {
			return "", errors.New("entity ID is not provided in both flag and client credentials")
		}
	}

	return entityId, nil
}

func FormatResponseAsJson(cmd *cobra.Command, response interface{}) (string, error) {
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

func GetFlagBoolValue(cmd *cobra.Command, flagName string) bool {
	value, _ := cmd.Flags().GetBool(flagName)
	return value
}

func NewUuid() uuid.UUID {
	return uuid.New()
}

func NewUuidStr() string {
	return NewUuid().String()
}
