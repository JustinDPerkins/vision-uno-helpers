package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Rule represents a single rule entry from the CSV
type Rule struct {
	Environment string
	Service     string
	RuleID      string
	RuleTitle   string
	Cost        float64
}

// RuleSummary holds the aggregated cost for each rule
type RuleSummary struct {
	RuleID    string
	RuleTitle string
	TotalCost float64
	Count     int
}

// EnvSummary holds the aggregated cost for each environment
type EnvSummary struct {
	Environment string
	TotalCost   float64
	Count       int
}

func main() {
	// Check for command line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <csv-file-path>")
	}
	
	// Open the CSV file using command line argument
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Unable to read input file:", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	// Read header to find column indices
	header, err := csvReader.Read()
	if err != nil {
		log.Fatal("Error reading header:", err)
	}

	// Find indices for our columns of interest
	var envIdx, serviceIdx, ruleIdIdx, ruleTitleIdx, costIdx int = -1, -1, -1, -1, -1
	for i, column := range header {
		switch strings.ToLower(strings.TrimSpace(column)) {
		case "environment":
			envIdx = i
		case "service":
			serviceIdx = i
		case "rule id":
			ruleIdIdx = i
		case "rule title":
			ruleTitleIdx = i
		case "cost":
			costIdx = i
		}
	}

	// Verify we found all required columns
	if envIdx == -1 || serviceIdx == -1 || ruleIdIdx == -1 || ruleTitleIdx == -1 || costIdx == -1 {
		log.Fatal("Missing required columns")
	}

	// Maps to store summaries
	ruleSummaries := make(map[string]*RuleSummary)
	envSummaries := make(map[string]*EnvSummary)
	rowCount := 0
	skippedCount := 0

	// Read and process each row
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			continue
		}
		rowCount++

		// Skip if we don't have enough columns
		if len(record) <= costIdx {
			skippedCount++
			continue
		}

		// Get values using found indices
		environment := record[envIdx]
		ruleID := record[ruleIdIdx]
		ruleTitle := record[ruleTitleIdx]
		costStr := record[costIdx]

		// Skip if required fields are empty
		if environment == "" || ruleID == "" || costStr == "" {
			skippedCount++
			continue
		}

		// Parse cost
		cost, err := strconv.ParseFloat(costStr, 64)
		if err != nil {
			log.Printf("Error parsing cost: %v (value: %s)", err, costStr)
			skippedCount++
			continue
		}

		// Update or create rule summary
		if summary, exists := ruleSummaries[ruleID]; exists {
			summary.TotalCost += cost
			summary.Count++
		} else {
			ruleSummaries[ruleID] = &RuleSummary{
				RuleID:    ruleID,
				RuleTitle: ruleTitle,
				TotalCost: cost,
				Count:     1,
			}
		}

		// Update or create environment summary
		if summary, exists := envSummaries[environment]; exists {
			summary.TotalCost += cost
			summary.Count++
		} else {
			envSummaries[environment] = &EnvSummary{
				Environment: environment,
				TotalCost:   cost,
				Count:       1,
			}
		}
	}

	// Print processing summary
	fmt.Printf("\nProcessing Summary:\n")
	fmt.Printf("Total rows processed: %d\n", rowCount)
	fmt.Printf("Rows skipped: %d\n", skippedCount)
	fmt.Printf("Unique rules found: %d\n", len(ruleSummaries))
	fmt.Printf("Unique environments found: %d\n\n", len(envSummaries))

	// Calculate grand total
	var grandTotal float64
	for _, summary := range envSummaries {
		grandTotal += summary.TotalCost
	}

	// Print environment breakdown
	fmt.Println("\nEnvironment Cost Breakdown:")
	fmt.Printf("%-20s %-15s %-10s\n", "Environment", "Total Cost", "Count")
	fmt.Println(strings.Repeat("-", 45))

	for _, summary := range envSummaries {
		fmt.Printf("%-20s $%-14.2f %-10d\n",
			summary.Environment,
			summary.TotalCost,
			summary.Count)
	}
	fmt.Println(strings.Repeat("-", 45))
	fmt.Printf("%-20s $%-14.2f\n", "GRAND TOTAL:", grandTotal)

	// Print rule breakdown
	fmt.Println("\nRule Cost Breakdown:")
	fmt.Printf("%-15s %-40s %-15s %-10s\n", "Rule ID", "Rule Title", "Total Cost", "Count")
	fmt.Println(strings.Repeat("-", 80))

	for _, summary := range ruleSummaries {
		fmt.Printf("%-15s %-40s $%-14.2f %-10d\n",
			summary.RuleID,
			truncateString(summary.RuleTitle, 37),
			summary.TotalCost,
			summary.Count)
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-56s $%-14.2f\n", "GRAND TOTAL:", grandTotal)
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
} 