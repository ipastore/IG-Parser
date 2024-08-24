package main

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"IG-Parser/web/converter/shared"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func SearchCodedStatementIdx(header []string) ([]int, HeaderMatchingError) {
	var indexes []int

	for i, cellString0 := range header {

		cellString := regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(cellString0, "")
		cellString1 := strings.ToLower(cellString)

		regStatement := regexp.MustCompile("(?:sta?t?e?m?e?n?t?)")
		regCoded := regexp.MustCompile("(?:co?d)")
		matchStatement := regStatement.MatchString(cellString1)
		matchCoded := regCoded.MatchString(cellString1)

		if matchStatement && matchCoded {
			// log.Println("Match: ", cellString0)
			indexes = append(indexes, i)
		}
		//  else {
		// 	log.Println("No Match: ", cellString0)
		// }
	}

	if len(indexes) == 0 {
		return nil, HeaderMatchingError{ErrorCode: HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
			ErrorMessage: "No matches for Coded Statement found in the input header"}
	} else if len(indexes) > 1 {
		return indexes, HeaderMatchingError{ErrorCode: HEADER_MATCHING_ERROR_MULTIPLE_MATCHES_FOR_CODED_STATEMENT,
			ErrorMessage: "Multiple matches for Coded Statement found in the input header"}
	} else {
		return indexes, HeaderMatchingError{ErrorCode: HEADER_MATCHING_NO_ERROR_MATCH_FOR_CODED_STATEMENT,
			ErrorMessage: "No Maatching Error"}
	}
}

// Indicates founding no matches in the input header for "Coded Statement"
const HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT = "NO_MATCH_FOR_CODED_STATEMENT"

// Indicates founding no matches in the input header for "Coded Statement"
const HEADER_MATCHING_ERROR_MULTIPLE_MATCHES_FOR_CODED_STATEMENT = "MULTIPLE_MATCHES_FOR_CODED_STATEMENT"

// Indicates founding no matches in the input header for "Coded Statement"
const HEADER_MATCHING_NO_ERROR_MATCH_FOR_CODED_STATEMENT = "NO_ERROR_FOR_CODED_STATEMENT"

/*
Error type signaling errors during searching for Coded Statement Column
*/
type HeaderMatchingError struct {
	ErrorCode    string
	ErrorMessage string
}

func (e *HeaderMatchingError) Error() error {
	return errors.New("Header Matching Error " + e.ErrorCode + ": " + e.ErrorMessage)
}

func main() {

	// Copy and append result It is copying without enriched data. I think it could read any matrix. Check of imbalanced matrices
	// f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX.xlsx")
	f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX_error.xlsx")
	// f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX_cdedStatementCOLUMN.xlsx")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f1.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get active Sheet to overcome bug of renaming sheeet and different languages
	activeSheet := f1.GetSheetName(f1.GetActiveSheetIndex())

	// Open new StreamWriter
	sw1, err := f1.NewStreamWriter(activeSheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get 2D array of the book with the stream writer
	matrix, err := f1.GetRows(activeSheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Inicializo stmtID y codedStatementColumn
	stmtId := "1"
	rowCoordinate := 1
	var codedStatementColumn int

	// Set Defaultconfig to overwrite later
	shared.SetDefaultConfig()
	// ACá falta inicializar el index para coordinatesToCellName: no puede ser r

	// General func to write on excel

	//Iterate over rows
	for r, rowMatrix := range matrix {

		//Initialize column index for Coded Statement
		// Row to append and write
		rowToWriteInterface := make([]interface{}, 0)
		errorMessageInterface := make([]interface{}, 1)
		// row to overwrite. Possible bug: rows cant be larger than the header of the matrix. POSSIBLE BUG in len(matrix[0])
		rowMatrixInterface := make([]interface{}, len(matrix[0]))
		//Copy to interface
		for i, v := range rowMatrix {
			rowMatrixInterface[i] = v
		}

		// Append Header
		if r == 0 {
			tabular.SetIncludeHeaders(true)

			//Search for column of Coded Statement
			arrayCodedStatementColumn, err0 := SearchCodedStatementIdx(rowMatrix)
			if err0.ErrorCode != HEADER_MATCHING_NO_ERROR_MATCH_FOR_CODED_STATEMENT {
				fmt.Println(err0.Error())
			}
			codedStatementColumn = arrayCodedStatementColumn[0]

			// Make ghost statement to print header
			ghostStatementToPrintHeader := "Cac{Once E(policy) F(comes into force)} A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance)."
			output, _ := endpoints.ConvertIGScriptToTabularOutput(ghostStatementToPrintHeader, stmtId, tabular.OUTPUT_TYPE_CSV, "", false, tabular.IncludeHeader(), tabular.DEFAULT_IG_SCRIPT_OUTPUT)

			headerArray := output[0].HeaderNames
			headerArray = append([]string{"Error"}, headerArray...)

			headerInterface := make([]interface{}, len(headerArray))
			for i, v := range headerArray {
				headerInterface[i] = v
			}

			rowToWriteInterface = append(rowMatrixInterface, headerInterface...)

			// Printo row from coordinateCell
			coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
			if err != nil {
				fmt.Println(err)
				break
			}
			if err := sw1.SetRow(coordinateCell, rowToWriteInterface); err != nil {
				fmt.Println(err)
				break
			}
			rowCoordinate += 1

		} else {
			tabular.SetIncludeHeaders(false)

			// Append coded staments

			output, err0 := endpoints.ConvertIGScriptToTabularOutput(rowMatrix[codedStatementColumn], stmtId, tabular.OUTPUT_TYPE_CSV, "", false, tabular.IncludeHeader(), tabular.DEFAULT_IG_SCRIPT_OUTPUT)

			//Getting the error and here I should paste it
			if err0.ErrorCode != tree.PARSING_NO_ERROR {
				errorMessageInterface[0] = err0.ErrorMessage

				rowToWriteInterface = append(rowMatrixInterface, errorMessageInterface)

				// Printo row from coordinateCell
				coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
				if err != nil {
					fmt.Println(err)
					break
				}
				if err := sw1.SetRow(coordinateCell, rowToWriteInterface); err != nil {
					fmt.Println(err)
					break
				}
				rowCoordinate += 1
				continue
			}

			errorMessageInterface[0] = "OK"

			stmtIdint, _ := strconv.Atoi(stmtId)
			stmtIdint += 1
			stmtId = strconv.Itoa(stmtIdint)

			//Primero append OK (cuando catchee el error)

			multipleOutputRows := strings.Split(output[0].Output, tabular.StmtIdPrefix) //Be aware here with the apostrophe
			//Leave out first value of Split func beaceuse its empty
			multipleOutputRows = multipleOutputRows[1:]

			//Catch the statements without adding a row
			for _, singleOutputRow := range multipleOutputRows {

				//leave out last value /n
				singleOutputRow := strings.Split(singleOutputRow, tabular.CellSeparator)
				if len(singleOutputRow) > 0 {
					singleOutputRow = singleOutputRow[:len(singleOutputRow)-1]
				}

				singleOutputRowInterface := make([]interface{}, len(singleOutputRow))
				for i, v := range singleOutputRow {
					singleOutputRowInterface[i] = v
				}

				//BUG en error MessageInterface
				rowToWriteInterface = append(rowMatrixInterface, errorMessageInterface)
				rowToWriteInterface = append(rowToWriteInterface, singleOutputRowInterface...)

				// Printo row from coordinateCell
				coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
				if err != nil {
					fmt.Println(err)
					break
				}
				if err := sw1.SetRow(coordinateCell, rowToWriteInterface); err != nil {
					fmt.Println(err)
					break
				}
				rowCoordinate += 1

			}

		}
	}

	if err := sw1.Flush(); err != nil {
		fmt.Println(err)
		return
	}

	if err := f1.SaveAs("/Users/ignaciopastorebenaim/go/src/IG-Parser/coded/200_MAX_coded.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
}
