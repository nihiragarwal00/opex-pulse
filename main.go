package main

import (
	"fmt"

	"github.com/razorpay/opex-pulse/src"
)

func main() {
	cliArgs, err := src.ParseTimeArgs()
	if err != nil {
		fmt.Printf("Error parsing time arguments: %v\n", err)
		return
	}

	cookie, err := src.ReadFileAndExtract(cliArgs.Cookie)
	if err != nil {
		fmt.Printf("Error reading cookie file: %v\n", err)
		return
	}
	cliArgs.Cookie = cookie

	metricConfigs, err := src.ParseMetricConfigCSV("sample.csv")
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	metricResult := src.GetAllMetrics(metricConfigs, *cliArgs)
	err = src.WriteMetricResultsToCSV(metricResult, "output.csv")
	if err != nil {
		fmt.Println("Error in converting to csv")
	}
}
