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

// RuleSummary holds the aggregated cost and count per rule
type RuleSummary struct {
	RuleID     string
	RuleTitle  string
	TotalCost  float64
	TotalSaved float64
	Count      int
}

// EnvSummary holds the aggregated cost and count per environment
type EnvSummary struct {
	Environment string
	TotalCost   float64
	TotalSaved  float64
	Count       int
}

func main() {
	// Ensure correct number of arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <csv-file-path> [discount-percentage]")
	}

	// Parse optional discount percentage
	discountPercentage := 0.0
	if len(os.Args) >= 3 {
		var err error
		discountPercentage, err = strconv.ParseFloat(os.Args[2], 64)
		if err != nil || discountPercentage < 0 || discountPercentage > 100 {
			log.Fatal("Invalid discount percentage. Provide a value between 0 and 100.")
		}
		discountPercentage /= 100 // Convert from percentage (e.g., 10 -> 0.10)
	}

	// Open the CSV file
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
	var envIdx, ruleIdIdx, ruleTitleIdx, costIdx, savingsIdx int = -1, -1, -1, -1, -1
	for i, column := range header {
		switch strings.ToLower(strings.TrimSpace(column)) {
		case "environment":
			envIdx = i
		case "rule id":
			ruleIdIdx = i
		case "rule title":
			ruleTitleIdx = i
		case "cost":
			costIdx = i
		case "savings":
			savingsIdx = i
		}
	}

	// Verify all required columns exist
	if envIdx == -1 || ruleIdIdx == -1 || ruleTitleIdx == -1 || costIdx == -1 || savingsIdx == -1 {
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

		// Skip if not enough columns
		if len(record) <= savingsIdx {
			skippedCount++
			continue
		}

		// Get values using found indices
		environment := record[envIdx]
		ruleID := record[ruleIdIdx]
		ruleTitle := record[ruleTitleIdx]
		costStr := record[costIdx]
		savingsStr := record[savingsIdx]

		// Assign a default value if Environment is empty
		if environment == "" {
			environment = "Unknown"
		}

		// Parse cost (allow empty values)
		cost, err := strconv.ParseFloat(costStr, 64)
		if err != nil || costStr == "" {
			cost = 0 // Default to zero if missing
		}

		// Apply discount
		cost *= (1 - discountPercentage)

		// Parse savings
		savings, err := strconv.ParseFloat(savingsStr, 64)
		if err != nil || savingsStr == "" {
			savings = 0 // Default to zero if missing
		}

		// Always count occurrences, even if cost is zero
		if summary, exists := ruleSummaries[ruleID]; exists {
			summary.TotalCost += cost
			summary.TotalSaved += savings
			summary.Count++ // Always count the rule, even if cost is zero
		} else {
			ruleSummaries[ruleID] = &RuleSummary{
				RuleID:     ruleID,
				RuleTitle:  ruleTitle,
				TotalCost:  cost,
				TotalSaved: savings,
				Count:      1,
			}
		}

		// Update or create environment summary
		if summary, exists := envSummaries[environment]; exists {
			summary.TotalCost += cost
			summary.TotalSaved += savings
			summary.Count++
		} else {
			envSummaries[environment] = &EnvSummary{
				Environment: environment,
				TotalCost:   cost,
				TotalSaved:  savings,
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
	fmt.Printf("Applied Discount: %.2f%%\n\n", discountPercentage*100)

	// Calculate grand total
	var grandTotalCost, grandTotalSavings float64
	for _, summary := range envSummaries {
		grandTotalCost += summary.TotalCost
		grandTotalSavings += summary.TotalSaved
	}

	// Print environment breakdown (without Total Savings)
	fmt.Println("\nEnvironment Cost Breakdown:")
	fmt.Printf("%-20s %-15s %-10s\n", "Environment", "Total Cost", "Count")
	fmt.Println(strings.Repeat("-", 45))

	for _, summary := range envSummaries {
		// Skip environments with 0.00 cost
		if summary.TotalCost > 0 {
			fmt.Printf("%-20s $%-14.2f %-10d\n",
				summary.Environment,
				summary.TotalCost,
				summary.Count)
		}
	}
	fmt.Println(strings.Repeat("-", 45))
	fmt.Printf("%-20s $%-14.2f\n", "GRAND TOTAL:", grandTotalCost)

	// Print rule breakdown (without Total Savings)
	fmt.Println("\nRule Cost Breakdown:")
	fmt.Printf("%-15s %-40s %-15s %-10s\n", "Rule ID", "Rule Title", "Total Cost", "Count")
	fmt.Println(strings.Repeat("-", 75))

	for _, summary := range ruleSummaries {
		// Skip rules with 0.00 cost
		if summary.TotalCost > 0 {
			fmt.Printf("%-15s %-40s $%-14.2f %-10d\n",
				summary.RuleID,
				truncateString(summary.RuleTitle, 37),
				summary.TotalCost,
				summary.Count)
		}
	}
	fmt.Println(strings.Repeat("-", 75))
	fmt.Printf("%-56s $%-14.2f\n", "GRAND TOTAL:", grandTotalCost)

}

// truncateString shortens a string to the specified length, adding "..." if truncated.
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
