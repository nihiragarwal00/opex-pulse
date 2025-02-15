package src

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHeaders(t *testing.T) {
	expectedHeaders := []string{"MetricName", "PromQLQuery", "Approved", "StatOperation"}

	// Test with matching headers
	actualHeaders := []string{"MetricName", "PromQLQuery", "Approved", "StatOperation"}
	assert.True(t, validateHeaders(actualHeaders, expectedHeaders))

	// Test with non-matching headers
	actualHeaders = []string{"MetricName", "PromQLQuery", "Approved", "InvalidHeader"}
	assert.False(t, validateHeaders(actualHeaders, expectedHeaders))
}

func TestParseRow(t *testing.T) {
	// Test with valid row
	row := []string{"Min. Latency (ms) (P95)", "histogram_quantile(0.95,sum(rate(oe_grpc_server_handling_seconds_bucket{k8s_pod=~\"offers-engine-live.*\", kubernetes_namespace=\"offers-engine\"}[5m])) by (le)) * 1000", "true", "MIN"}
	config, err := parseRow(row)
	assert.NoError(t, err)
	assert.Equal(t, "Min. Latency (ms) (P95)", config.MetricName)

	// Test with invalid row (incorrect length)
	row = []string{"Min. Latency (ms) (P95)", "histogram_quantile(0.95,sum(rate(oe_grpc_server_handling_seconds_bucket{k8s_pod=~\"offers-engine-live.*\", kubernetes_namespace=\"offers-engine\"}[5m])) by (le)) * 1000", "true"}
	_, err = parseRow(row)
	assert.Error(t, err)
}

func TestParseTimeArgs(t *testing.T) {
	// Test with valid arguments
	os.Args = []string{"cmd", "-start=1633046400", "-end=1633132800", "-cookie_file=cookie.txt", "-timeout=30"}
	cliArgs, err := ParseTimeArgs()
	assert.NoError(t, err)
	assert.Equal(t, int64(1633046400000), cliArgs.Start)
	assert.Equal(t, int64(1633132800000), cliArgs.End)
	assert.Equal(t, "cookie.txt", cliArgs.Cookie)
	assert.Equal(t, int64(30), cliArgs.Timeout)
}

func TestReadFileAndExtract(t *testing.T) {
	// Create a temporary file
	file, err := os.CreateTemp("", "testfile.txt")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Write content to the file
	content := "This is a test file.\nWith multiple lines.\n"
	_, err = file.WriteString(content)
	assert.NoError(t, err)
	file.Close()

	// Test with valid file
	result, err := ReadFileAndExtract(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, "This is a test file.\nWith multiple lines.", result)

	// Test with non-existent file
	_, err = ReadFileAndExtract("non_existent_file.txt")
	assert.Error(t, err)
}
