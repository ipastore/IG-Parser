package production

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// Func to log the content of an Excel file
func logExcelContent(file *excelize.File) {
	for _, sheetName := range file.GetSheetList() {
		log.Printf("Sheet: %s\n", sheetName)
		rows, err := file.GetRows(sheetName)
		if err != nil {
			log.Printf("Error reading rows from sheet %s: %v\n", sheetName, err)
			continue
		}
		for i, row := range rows {
			log.Printf("Row %d: %v\n", i, row)
		}
	}
}

// Compare two slices of integers
func compareSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CompareExcelFiles compares the contents of two Excel files and returns an error if there are any mismatches.
func compareExcelFiles(expectedF, actualF *excelize.File) error {
	expectedSheets := expectedF.GetSheetList()
	actualSheets := actualF.GetSheetList()

	if len(expectedSheets) != len(actualSheets) {
		return fmt.Errorf("number of sheets mismatch: expected %d, got %d", len(expectedSheets), len(actualSheets))
	}

	for _, sheet := range expectedSheets {
		expectedRows, err := expectedF.GetRows(sheet)
		if err != nil {
			return fmt.Errorf("failed to get rows from expected file: %v", err)
		}
		actualRows, err := actualF.GetRows(sheet)
		if err != nil {
			return fmt.Errorf("failed to get rows from actual file: %v", err)
		}

		if len(expectedRows) != len(actualRows) {
			return fmt.Errorf("number of rows mismatch in sheet %s: expected %d, got %d", sheet, len(expectedRows), len(actualRows))
		}

		for i, expectedRow := range expectedRows {
			actualRow := actualRows[i]
			if len(expectedRow) != len(actualRow) {
				return fmt.Errorf("number of columns mismatch in sheet %s, row %d: expected %d, got %d", sheet, i+1, len(expectedRow), len(actualRow))
			}

			for j, expectedCell := range expectedRow {
				actualCell := actualRow[j]
				if expectedCell != actualCell {
					return fmt.Errorf("cell mismatch in sheet %s, row %d, column %d: expected %q, got %q", sheet, i+1, j+1, expectedCell, actualCell)
				}
			}
		}
	}

	return nil
}
