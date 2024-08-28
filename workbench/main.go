package main

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"errors"
	"fmt"
	"regexp"
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

	// ghostStatementToPrintHeader := "Cac{Once E(policy) F(comes into force)} A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance)."
	ghostStatementToPrintHeader := "A,p(relevant) A(regulators) D(must) I(enforce) Bdir(compliance)."
	output, _ := endpoints.ConvertIGScriptToTabularOutput(ghostStatementToPrintHeader, "1", tabular.OUTPUT_TYPE_CSV, "", false, tabular.IncludeHeader(), tabular.DEFAULT_IG_SCRIPT_OUTPUT)

	//run through all the output and print: Output  StatementMap, HeaderSymbols, HeaderNames adn Error

	// for _, output := range output {
	// fmt.Println("Output", output.Output) // Es el output con apostrofe y lineas separadoras de columnas, esta es la que use y no deberia usar
	// fmt.Println("Error", output.Error)                  // Es el error que se genera al convertir el script a tabular
	// fmt.Println("Header Symbols", output.HeaderSymbols) // Es el header con los simbolos de las columnas
	// fmt.Println("Header Names", output.HeaderNames)
	// fmt.Println("Type of Header Names", reflect.TypeOf(output.HeaderNames))

	// loop to print the statement maps
	// fmt.Println("Len of StatementMap", len(output.StatementMap))
	// for _, statementMap := range output.StatementMap {
	// 	fmt.Println("StatementMap", statementMap)
	// 	// loop to run through the map through the headerSymbols and print the pair Key:value. if there is no value, print an empty string
	// 	for _, headerSymbol := range output.HeaderSymbols {
	// 		if val, ok := statementMap[headerSymbol]; ok {
	// 			fmt.Println(headerSymbol, val)
	// 		} else {
	// 			fmt.Println(headerSymbol, "")
	// 		}
	// 	}
	// }

	// Open file anf check for errors
	// f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX.xlsx")
	f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/empty.xlsx")
	// f1, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX_cdedStatementCOLUMN.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	//Deferred close file
	defer func() {
		if err := f1.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get active Sheet to overcome bug of renaming sheeet and different languages
	activeSheet := f1.GetSheetName(f1.GetActiveSheetIndex())

	// Open new StreamWriter
	sw, err := f1.NewStreamWriter(activeSheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	rowCoordinate := 1

	// Loop to write each iteration of the statement map
	for _, statementMap := range output[0].StatementMap {
		rowToWriteInterface := make([]interface{}, len(output[0].HeaderSymbols))

		// Loop to run through the map through the headerSymbols and print the pair Key:value. if there is no value, print an empty string
		for i, headerSymbol := range output[0].HeaderSymbols {
			if val, ok := statementMap[headerSymbol]; ok {
				rowToWriteInterface[i] = val
			} else {
				rowToWriteInterface[i] = ""
			}
		}

		// Printo row from coordinateCell
		coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := sw.SetRow(coordinateCell, rowToWriteInterface); err != nil {
			fmt.Println(err)
			return
		}

		rowCoordinate++
	}

	// Flush the stream writer
	if err := sw.Flush(); err != nil {
		fmt.Println(err)
		return
	}

	// Save the file
	if err := f1.SaveAs("/Users/ignaciopastorebenaim/go/src/IG-Parser/coded/empty_coded.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

}
