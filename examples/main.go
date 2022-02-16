package main

import (
	"fmt"
	"flag"
	. "github.com/hhxsv5/go-redis-memory-analysis"
)

func main() {
	limitCount := flag.Int("limitcount", 100, "limit count")

	//Open redis: 127.0.0.1:6379 without password
	analysis, err := NewAnalysisConnection("127.0.0.1", 6379, "")
	if err != nil {
		fmt.Println("something wrong:", err)
		return
	}
	defer analysis.Close()
    
	analysis.Start(*limitCount)

	//Find the csv file in default target folder: ./reports
	//CSV file name format: redis-analysis-{host:port}-{db}.csv
	//The keys order by count desc
	err = analysis.SaveReports("./reports")
	if err == nil {
		fmt.Println("done")
	} else {
		fmt.Println("error:", err)
	}
}
