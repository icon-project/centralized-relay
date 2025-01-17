package cmd

import (
	"fmt"
	"strings"
)

func displayTable(headers []string, items [][]interface{}) {
	// Calculate column widths
	colWidths := make([]int, len(headers))
	padding := 6
	for i, header := range headers {
		colWidths[i] = len(header) + padding // Add padding
		for _, row := range items {
			cell := fmt.Sprintf("%v", row[i]) // Convert cell to string
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell) + padding // Update width with padding
			}
		}
	}

	// Print the header
	for i, header := range headers {
		fmt.Printf("%-*s", colWidths[i], header)
	}
	fmt.Println()

	// Print a separator line
	fmt.Println(strings.Repeat("-", sum(colWidths)-padding))

	// Print the rows
	for _, row := range items {
		for i, cell := range row {
			fmt.Printf("%-*v", colWidths[i], cell) // Format the cell value
		}
		fmt.Println()
	}
}

func sum(widths []int) int {
	total := 0
	for _, width := range widths {
		total += width
	}
	return total
}
