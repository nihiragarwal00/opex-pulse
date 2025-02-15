### README

# Metric Configuration and Query Tool

This project provides tools for parsing metric configurations from a CSV file, querying metrics from Grafana, and writing the results to a CSV file.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Functions](#functions)
    - [ParseMetricConfigCSV](#parsemetricconfigcsv)
    - [validateHeaders](#validateheaders)
    - [parseRow](#parserow)
    - [ParseTimeArgs](#parsetimeargs)
    - [ReadFileAndExtract](#readfileandextract)
    - [GetAllMetrics](#getallmetrics)
    - [WriteMetricResultsToCSV](#writemetricresultstocsv)
- [Testing](#testing)

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/nihiragarwal00/opex-pulse.git
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

1. Prepare a CSV file with metric configurations. Check Sample file for reference.
2. Run the main program with the required arguments:
   ```sh
   go run main.go -start=<start_epoch> -end=<end_epoch> -cookie_file=<cookie_file> -timeout=<timeout>
   ```
If end epoch is not provided, it defaults to the current time.   

If you want to use the sample file, you can run the following command:
   ```sh
   go run main.go -start=1735385000 -end=1739557800 -cookie_file="cookie.txt" -timeout=50
   ```

## Testing

To run the tests, use the following command:
```sh
go test ./...
```