package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/razorpay/opex-pulse/src"
)

type MetricConfig struct {
	MetricName    string
	PromQLQuery   string
	Approved      bool
	StatOperation src.StatOperation
}

type MetricResult struct {
	MetricConfig
	Value float64
}

type CLIArgs struct {
	Start  int64
	End    int64
	Cookie string
}

func ParseMetricConfigCSV(filename string) ([]MetricConfig, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	// Validate header columns
	expectedHeaders := []string{"MetricName", "PromQLQuery", "Approved", "StatOperation"}
	if !validateHeaders(header, expectedHeaders) {
		return nil, fmt.Errorf("invalid headers. Expected: %v, Got: %v", expectedHeaders, header)
	}

	var configs []MetricConfig

	// Read the rest of the rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}

		// Parse the row into MetricConfig
		config, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %v", err)
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// validateHeaders checks if the CSV headers match the expected headers
func validateHeaders(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, header := range actual {
		if strings.TrimSpace(strings.ToLower(header)) != strings.ToLower(expected[i]) {
			return false
		}
	}
	return true
}

// parseRow converts a CSV row into a MetricConfig struct
func parseRow(row []string) (MetricConfig, error) {
	if len(row) != 4 {
		return MetricConfig{}, fmt.Errorf("invalid row length: expected 4, got %d", len(row))
	}

	// Parse approved field
	//approved := false
	//if strings.ToLower(strings.TrimSpace(row[2])) == "true" {
	//	approved = true
	//}

	return MetricConfig{
		MetricName:    strings.TrimSpace(row[0]),
		PromQLQuery:   strings.TrimSpace(row[1]),
		Approved:      true,
		StatOperation: src.StatOperation(strings.TrimSpace(row[3])),
	}, nil
}

func ParseTimeArgs() (*CLIArgs, error) {
	// Define command line flags
	var startEpoch int64
	var endEpoch int64
	var cookie string
	flag.Int64Var(&startEpoch, "start", 0, "Start time in epoch seconds")
	flag.Int64Var(&endEpoch, "end", 0, "End time in epoch seconds (defaults to current time)")
	flag.StringVar(&cookie, "cookie", "", "Cookie is required")

	// Parse the flags
	flag.Parse()

	// Validate start time is provided
	if startEpoch == 0 {
		return nil, fmt.Errorf("start time is required")
	}

	// If end time is not provided, use current time
	if endEpoch == 0 {
		endEpoch = time.Now().Unix()
	}

	// Validate time range
	if startEpoch >= endEpoch {
		return nil, fmt.Errorf("start time must be before end time")
	}

	if cookie == "" {
		return nil, fmt.Errorf("your cookie is very essential")
	}

	return &CLIArgs{
		Start:  startEpoch * 1000,
		End:    endEpoch * 1000,
		Cookie: cookie,
	}, nil
}

type GrafanaQuery struct {
	Queries []Query `json:"queries"`
	From    string  `json:"from"`
	To      string  `json:"to"`
}

// Query represents a single Grafana query
type Query struct {
	Datasource     Datasource `json:"datasource"`
	Exemplar       bool       `json:"exemplar"`
	Expr           string     `json:"expr"`
	Format         string     `json:"format"`
	Hide           bool       `json:"hide"`
	Interval       string     `json:"interval"`
	IntervalFactor int        `json:"intervalFactor"`
	LegendFormat   string     `json:"legendFormat"`
	RefID          string     `json:"refId"`
	EditorMode     string     `json:"editorMode"`
	Range          bool       `json:"range"`
	RequestID      string     `json:"requestId"`
	UtcOffsetSec   int        `json:"utcOffsetSec"`
	DatasourceID   int        `json:"datasourceId"`
	IntervalMs     int        `json:"intervalMs"`
	MaxDataPoints  int        `json:"maxDataPoints"`
}

// Datasource represents the Grafana datasource
type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// GrafanaClient represents a client for querying Grafana
type GrafanaClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewGrafanaClient creates a new Grafana client
func NewGrafanaClient(baseURL string, timeout time.Duration, cookie string) *GrafanaClient {
	return &GrafanaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: map[string]string{
			"accept":           "application/json, text/plain, */*",
			"content-type":     "application/json",
			"x-grafana-org-id": "1",
			"x-plugin-id":      "prometheus",
			"x-datasource-uid": "000000002",
			"x-dashboard-uid":  "E2KboMgBY",
			"cookie":           cookie,
		},
	}
}

// QueryMetrics queries Grafana with the given expression and time range
func (c *GrafanaClient) QueryMetrics(expr string, start, end int64) (*src.Root, error) {
	query := GrafanaQuery{
		Queries: []Query{
			{
				Datasource: Datasource{
					Type: "prometheus",
					UID:  "000000002",
				},
				Exemplar:       true,
				Expr:           expr,
				Format:         "time_series",
				Hide:           false,
				IntervalFactor: 1,
				RefID:          "C",
				EditorMode:     "code",
				Range:          true,
				RequestID:      "305C",
				UtcOffsetSec:   19800,
				DatasourceID:   2,
				IntervalMs:     60000,
				MaxDataPoints:  1320,
			},
		},
		From: fmt.Sprintf("%d", start),
		To:   fmt.Sprintf("%d", end),
	}

	// Marshal the query to JSON
	payload, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", c.baseURL+"/api/ds/query?ds_type=prometheus&requestId=Q256", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error from server: %s", string(body))
	}
	var responseRoot src.Root
	if err = json.Unmarshal(body, &responseRoot); err != nil {
		return nil, err
	}
	return &responseRoot, nil
}
func ConvertToFloat64(input []interface{}) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		if num, ok := v.(float64); ok {
			result[i] = num
		} else {
			return nil, fmt.Errorf("invalid type at index %d: expected float64, got %T", i, v)
		}
	}
	return result, nil
}
func getStatistic(operation src.StatOperation, rootResponse *src.Root) float64 {
	values, _ := ConvertToFloat64(rootResponse.Results["C"].Frames[0].Data.Values[1])
	return RoundToTwoDecimals(src.StatsFuncs[operation](values))
}

func RoundToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}

func GetAllMetrics(metricConfigs []MetricConfig, cliArgs CLIArgs) []MetricResult {
	client := NewGrafanaClient("https://vajra.razorpay.com", 10*time.Second, cliArgs.Cookie)
	metricsResult := []MetricResult{}
	for _, metricConfig := range metricConfigs {
		if metricConfig.StatOperation != src.NotSet {
			resp, err := client.QueryMetrics(metricConfig.PromQLQuery, cliArgs.Start, cliArgs.End)
			if err != nil {
				fmt.Printf("Error querying metrics: %v\n", err)
				continue
			}
			metricResult := MetricResult{
				MetricConfig: metricConfig,
				Value:        getStatistic(metricConfig.StatOperation, resp),
			}
			metricsResult = append(metricsResult, metricResult)
		}
	}
	return metricsResult
}
func WriteMetricResultsToCSV(results []MetricResult, filename string) error {
	// Open file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"MetricName", "PromQLQuery", "Approved", "StatOperation", "Value"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write data rows
	for _, result := range results {
		row := []string{
			result.MetricName,
			result.PromQLQuery,
			strconv.FormatBool(result.Approved),
			string(result.StatOperation),
			fmt.Sprintf("%.2f", result.Value), // Round to 2 decimal places
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}

	return nil
}

func main() {
	cliArgs, err := ParseTimeArgs()
	if err != nil {
		fmt.Printf("Error parsing time arguments: %v\n", err)
		return
	}
	metricConfigs, err := ParseMetricConfigCSV("sample.csv")
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	metricResult := GetAllMetrics(metricConfigs, *cliArgs)
	err = WriteMetricResultsToCSV(metricResult, "output.csv")
	if err != nil {
		fmt.Println("Error in converting to csv")
	}

}
