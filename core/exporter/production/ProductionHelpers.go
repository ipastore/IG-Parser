package production

import (
	"log"

	"github.com/xuri/excelize/v2"
)

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
