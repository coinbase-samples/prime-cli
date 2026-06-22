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

package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coinbase/prime-sdk-go/client"
	"github.com/coinbase/prime-sdk-go/model"
	mcplib "github.com/mark3labs/mcp-go/mcp"
)

func resolvePortfolioId(c client.RestClient, req mcplib.CallToolRequest) (string, error) {
	id := req.GetString("portfolio_id", "")
	if id != "" {
		return id, nil
	}
	creds := c.Credentials()
	if creds == nil {
		return "", errors.New("client credentials are nil")
	}
	if creds.PortfolioId == "" {
		return "", errors.New("portfolio_id not provided and not set in credentials")
	}
	return creds.PortfolioId, nil
}

func resolveEntityId(c client.RestClient, req mcplib.CallToolRequest) (string, error) {
	id := req.GetString("entity_id", "")
	if id != "" {
		return id, nil
	}
	creds := c.Credentials()
	if creds == nil {
		return "", errors.New("client credentials are nil")
	}
	if creds.EntityId == "" {
		return "", errors.New("entity_id not provided and not set in credentials")
	}
	return creds.EntityId, nil
}

// mcpCtx creates a timeout context rooted at the MCP-provided parent so that
// client cancellation propagates and the primeCliTimeout env var is honoured.
func mcpCtx(parent context.Context) (context.Context, context.CancelFunc) {
	d := 7 * time.Second
	if env := os.Getenv("primeCliTimeout"); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			d = time.Duration(v) * time.Second
		}
	}
	return context.WithTimeout(parent, d)
}

// fetchAllCtx creates a longer-lived context for fetch_all operations that
// may need to traverse many pages.
func fetchAllCtx(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, 120*time.Second)
}

func paginationFor(req mcplib.CallToolRequest) *model.PaginationParams {
	return &model.PaginationParams{
		Cursor: req.GetString("cursor", ""),
		Limit:  int32(req.GetInt("limit", 50)),
	}
}

func marshalResult(v any) (*mcplib.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("internal: cannot marshal response: %w", err)
	}
	return mcplib.NewToolResultText(string(data)), nil
}

func toolErr(format string, a ...any) *mcplib.CallToolResult {
	return mcplib.NewToolResultError(fmt.Sprintf(format, a...))
}

// networkDetailsFor splits a compound network ID (e.g. "base-mainnet") into
// the Id ("base") and Type ("mainnet") fields the API requires.
func networkDetailsFor(networkId string) *model.NetworkDetails {
	if networkId == "" {
		return nil
	}
	if idx := strings.LastIndex(networkId, "-"); idx >= 0 {
		return &model.NetworkDetails{Id: networkId[:idx], Type: networkId[idx+1:]}
	}
	return &model.NetworkDetails{Id: networkId}
}
