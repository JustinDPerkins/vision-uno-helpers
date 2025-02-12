# Cloud Cost Savings Analyzer

A Go-based tool that analyzes cloud cost savings data from CSV files and provides breakdowns by environment and rule.

## Prerequisites

- Go 1.21.3 or later

## Installation

Clone the repository and navigate to the directory:

```bash
git clone <repository-url>
cd cloud/cost-savings-evp
```

## Usage

Run the program with a CSV file as input:
```bash
go run main.go <path-to-csv-file>
```

### CSV File Format

The input CSV file should contain the following columns:
- Environment
- Service
- Rule ID
- Rule Title
- Cost

Example CSV format:
```csv
Environment,Service,Rule ID,Rule Title,Cost
Production,EC2,R1,Underutilized Instances,100.50
Staging,RDS,R2,Idle Databases,75.25
```

## Output

The tool provides three main sections of output:

1. Processing Summary
   - Total rows processed
   - Rows skipped
   - Number of unique rules
   - Number of unique environments

2. Environment Cost Breakdown
   - Lists costs by environment
   - Shows percentage of total costs
   - Provides subtotals for each environment

3. Rule Analysis
   - Breaks down costs by rule
   - Shows frequency of rule occurrence
   - Highlights top cost-saving opportunities

## Example Output

```
=== Processing Summary ===
Total Rows: 150
Skipped Rows: 0
Unique Rules: 12
Unique Environments: 3

=== Environment Breakdown ===
Production: $1,234.56 (45%)
Staging:    $823.45 (30%)
Dev:        $685.90 (25%)
Total:      $2,743.91

=== Top Rules by Cost ===
1. Underutilized Instances: $856.78
2. Idle Databases: $534.67
3. Oversized Resources: $423.45
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.