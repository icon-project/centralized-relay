package cmd

import (
	"fmt"
	"strings"
)

func displayTable(headers []string, rows [][]interface{}) {
	strRows := [][]string{}
	for _, rowValues := range rows {
		strRow := []string{}
		for _, val := range rowValues {
			strRow = append(strRow, fmt.Sprintf("%v   ", val))
		}
		strRows = append(strRows, strRow)
	}

	colWidths := make([]int, len(headers))
	padding := 1

	for i, header := range headers {
		colWidths[i] = len(header) + padding
		for _, row := range strRows {
			cell := row[i]
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell) + padding
			}
		}
	}

	for i, header := range headers {
		fmt.Printf("%-*s", colWidths[i], header)
	}
	fmt.Println()

	fmt.Println(strings.Repeat("-", sum(colWidths)-padding))

	for _, row := range strRows {
		for i, cell := range row {
			fmt.Printf("%-*v", colWidths[i], cell)
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
