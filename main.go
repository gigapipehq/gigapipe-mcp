package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultGigapipeHost = "localhost:3100"
)

type Config struct {
	Host     string
	Username string
	Password string
}

func getConfig() Config {
	host := os.Getenv("GIGAPIPE_HOST")
	if host == "" {
		host = defaultGigapipeHost
	}

	return Config{
		Host:     host,
		Username: os.Getenv("GIGAPIPE_USERNAME"),
		Password: os.Getenv("GIGAPIPE_PASSWORD"),
	}
}

func newHTTPClient(config Config) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func makeRequest(ctx context.Context, client *http.Client, config Config, method, path string, params url.Values) ([]byte, error) {
	baseURL := fmt.Sprintf("http://%s", config.Host)
	if config.Username != "" && config.Password != "" {
		baseURL = fmt.Sprintf("https://%s", config.Host)
	}

	reqURL := fmt.Sprintf("%s%s", baseURL, path)
	if params != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return json.Marshal(result)
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Gigapipe MCP 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add Prometheus query tool
	promTool := mcp.NewTool("prometheus_query",
		mcp.WithDescription("Query metrics from Prometheus API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("PromQL query to execute"),
		),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("step",
			mcp.Description("Query resolution step width (e.g. 15s, 1m, 1h)"),
		),
	)

	// Add Prometheus labels tool
	promLabelsTool := mcp.NewTool("prometheus_labels",
		mcp.WithDescription("List all available label names in Prometheus"),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
	)

	// Add Prometheus label values tool
	promLabelValuesTool := mcp.NewTool("prometheus_label_values",
		mcp.WithDescription("List all values for a specific label in Prometheus"),
		mcp.WithString("label",
			mcp.Required(),
			mcp.Description("Label name to get values for"),
		),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
	)

	// Add Loki query tool
	lokiTool := mcp.NewTool("loki_query",
		mcp.WithDescription("Query logs from Loki API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("LogQL query to execute"),
		),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("limit",
			mcp.Description("Maximum number of entries to return"),
		),
	)

	// Add Loki labels tool
	lokiLabelsTool := mcp.NewTool("loki_labels",
		mcp.WithDescription("List all available label names in Loki"),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
	)

	// Add Loki label values tool
	lokiLabelValuesTool := mcp.NewTool("loki_label_values",
		mcp.WithDescription("List all values for a specific label in Loki"),
		mcp.WithString("label",
			mcp.Required(),
			mcp.Description("Label name to get values for"),
		),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
	)

	// Add Tempo query tool
	tempoTool := mcp.NewTool("tempo_query",
		mcp.WithDescription("Query traces from Tempo API"),
		mcp.WithString("trace_id",
			mcp.Required(),
			mcp.Description("Trace ID to query"),
		),
		mcp.WithString("start",
			mcp.Description("Start timestamp (RFC3339 or Unix timestamp)"),
		),
		mcp.WithString("end",
			mcp.Description("End timestamp (RFC3339 or Unix timestamp)"),
		),
	)

	// Add Tempo tags tool
	tempoTagsTool := mcp.NewTool("tempo_tags",
		mcp.WithDescription("List all available trace tags in Tempo"),
	)

	// Add Tempo tag values tool
	tempoTagValuesTool := mcp.NewTool("tempo_tag_values",
		mcp.WithDescription("List all values for a specific trace tag in Tempo"),
		mcp.WithString("tag",
			mcp.Required(),
			mcp.Description("Tag name to get values for"),
		),
	)

	// Add tool handlers
	s.AddTool(promTool, prometheusQueryHandler)
	s.AddTool(promLabelsTool, prometheusLabelsHandler)
	s.AddTool(promLabelValuesTool, prometheusLabelValuesHandler)
	s.AddTool(lokiTool, lokiQueryHandler)
	s.AddTool(lokiLabelsTool, lokiLabelsHandler)
	s.AddTool(lokiLabelValuesTool, lokiLabelValuesHandler)
	s.AddTool(tempoTool, tempoQueryHandler)
	s.AddTool(tempoTagsTool, tempoTagsHandler)
	s.AddTool(tempoTagValuesTool, tempoTagValuesHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func prometheusQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := request.GetString("start", "")
	end := request.GetString("end", "")
	step := request.GetString("step", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	params.Set("query", query)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	if step != "" {
		params.Set("step", step)
	}

	result, err := makeRequest(ctx, client, config, "GET", "/api/v1/query_range", params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func prometheusLabelsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := request.GetString("start", "")
	end := request.GetString("end", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}

	result, err := makeRequest(ctx, client, config, "GET", "/api/v1/labels", params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func prometheusLabelValuesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	label, err := request.RequireString("label")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := request.GetString("start", "")
	end := request.GetString("end", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}

	result, err := makeRequest(ctx, client, config, "GET", fmt.Sprintf("/api/v1/label/%s/values", label), params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func lokiQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := request.GetString("start", "")
	end := request.GetString("end", "")
	limit := request.GetString("limit", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	params.Set("query", query)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	if limit != "" {
		params.Set("limit", limit)
	}

	result, err := makeRequest(ctx, client, config, "GET", "/loki/api/v1/query_range", params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func lokiLabelsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := request.GetString("start", "")
	end := request.GetString("end", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}

	result, err := makeRequest(ctx, client, config, "GET", "/loki/api/v1/label", params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func lokiLabelValuesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	label, err := request.RequireString("label")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := request.GetString("start", "")
	end := request.GetString("end", "")

	config := getConfig()
	client := newHTTPClient(config)

	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}

	result, err := makeRequest(ctx, client, config, "GET", fmt.Sprintf("/loki/api/v1/label/%s/values", label), params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func tempoQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traceID, err := request.RequireString("trace_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	config := getConfig()
	client := newHTTPClient(config)

	result, err := makeRequest(ctx, client, config, "GET", fmt.Sprintf("/api/traces/%s/json", traceID), nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func tempoTagsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config := getConfig()
	client := newHTTPClient(config)

	result, err := makeRequest(ctx, client, config, "GET", "/api/search/tags", nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func tempoTagValuesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag, err := request.RequireString("tag")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	config := getConfig()
	client := newHTTPClient(config)

	result, err := makeRequest(ctx, client, config, "GET", fmt.Sprintf("/api/search/tag/%s/values", tag), nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
} 