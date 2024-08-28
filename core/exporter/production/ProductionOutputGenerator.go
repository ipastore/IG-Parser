package production

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ConvertExcelToExcelWithTabularOutput(r *http.Request) (string, ProductionError) {

	// Upload file to uploads folder
	filename, uploadPath, err := UploadExcelFile(r)
	if err.ErrorCode != PRODUCTION_NO_ERROR {
		return "", err
	}

	// Process file
	savePath, err := ProcessExcelFile(uploadPath, filename)
	if err.ErrorCode != PRODUCTION_NO_ERROR {
		return "", err
	}

	// // Remove file from uploads folder
	// err = RemoveFileFromUploads(uploadPath)
	// if err.ErrorCode != PRODUCTION_NO_ERROR {
	// 	return "", err
	// }

	return savePath, ProductionError{ErrorCode: PRODUCTION_NO_ERROR}
}

// helpers to write the file etc
// I need to call this functionn within this package in ProductionOutputGenerator.go
func SearchCodedStatementIdx(header []string) ([]int, ProductionError) {
	var indexes []int

	for i, cellString0 := range header {

		cellString := regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(cellString0, "")
		cellString1 := strings.ToLower(cellString)

		regStatement := regexp.MustCompile("(?:sta?t?e?m?e?n?t?)")
		regCoded := regexp.MustCompile("(?:co?d)")
		matchStatement := regStatement.MatchString(cellString1)
		matchCoded := regCoded.MatchString(cellString1)

		if matchStatement && matchCoded {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		errorMsg := "No matches for Coded Statement found in the input header"
		return nil, ProductionError{
			ErrorCode:    HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
			ErrorMessage: errorMsg}
	} else if len(indexes) > 1 {
		errorMsg := "Multiple matches for Coded Statement found in the input header"
		return indexes, ProductionError{
			ErrorCode:    HEADER_MATCHING_ERROR_MULTIPLE_MATCHES_FOR_CODED_STATEMENT,
			ErrorMessage: errorMsg}
	}
	return indexes, ProductionError{ErrorCode: PRODUCTION_NO_ERROR}

}

// Upload File to uploads
func UploadExcelFile(r *http.Request) (string, string, ProductionError) {

	// Parse the form
	if err := r.ParseMultipartForm(FORM_FILE_SIZE); err != nil {
		errorMsg := "Failed to parse the form."
		log.Println(err.Error())
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_PARSING_FORM,
			ErrorMessage: errorMsg}
	}

	// Create the uploads directory if it does not exist
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		errorMsg := "Failed to create the uploads directory."
		log.Println(err.Error())
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_SAVE_ERROR_CREATING_TEMP_FOLDER,
			ErrorMessage: errorMsg}
	}

	// Get the file from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		errorMsg := "Failed to get the file from the form."
		log.Println(err.Error())
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_GETTING_FILE,
			ErrorMessage: errorMsg}
	}

	// Defer Close
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	// Manage possible errors
	// Check if the file is too big
	if fileHeader.Size > MAX_UPLOAD_SIZE {
		errorMsg := "The uploaded file is too big: " + fmt.Sprintf("%d", fileHeader.Size)
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_FILE_TOO_BIG,
			ErrorMessage: errorMsg}
	}

	filename := fileHeader.Filename
	// Check if the file is not an excel file
	if fileHeader.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		errorMsg := "The uploaded file is not an excel file: " + filename
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_NOT_EXCEL_FILE,
			ErrorMessage: errorMsg}
	}

	// Create the file
	uploadPath := "./uploads/" + filename
	newFile, err := os.Create(uploadPath)
	if err != nil {
		errorMsg := "Failed to create the file."
		log.Println(err.Error())
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_CREATING_FILE,
			ErrorMessage: errorMsg}
	}

	// Defer Close
	defer func() {
		if err := newFile.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	// Copy the input file to the temp file in uploads
	if _, err := io.Copy(newFile, file); err != nil {
		errorMsg := "Failed to copy the file."
		log.Println(err.Error())
		return "", "", ProductionError{
			ErrorCode:    UPLOAD_ERROR_COPYING_FILE,
			ErrorMessage: errorMsg}
	}

	// return the file header and the path with no error
	return filename, uploadPath, ProductionError{ErrorCode: PRODUCTION_NO_ERROR}
}

// Save File
func SavExcelizeFile(file *excelize.File, filename string) (string, ProductionError) {

	// Create the coded directory if it does not exist
	if err := os.MkdirAll("./coded", os.ModePerm); err != nil {
		errorMsg := "Failed to create the coded directory."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    UPLOAD_SAVE_ERROR_CREATING_TEMP_FOLDER,
			ErrorMessage: errorMsg}
	}

	// Get the name, and append CODED to it
	saveAsFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + "_CODED.xlsx"
	savePath := "./coded/" + saveAsFilename

	// Save the file
	if err := file.SaveAs(savePath); err != nil {
		errorMsg := "Failed to save the excelize file."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    SAVE_ERROR_SAVING_FILE,
			ErrorMessage: errorMsg}
	}

	return savePath, ProductionError{ErrorCode: PRODUCTION_NO_ERROR}
}

// Remove file from /uploads
func RemoveFileFromUploads(uploadPath string) ProductionError {
	// Erase the file from uploads
	if err := os.Remove(uploadPath); err != nil {
		errorMsg := "Failed to remove the file from uploads."
		log.Println(err.Error())
		return ProductionError{
			ErrorCode:    REMOVE_ERROR_ERASING_FILE,
			ErrorMessage: errorMsg}
	}

	return ProductionError{ErrorCode: PRODUCTION_NO_ERROR}
}

// Process excel file with Excelize and the engine of the parser

func ProcessExcelFile(uploadPath string, filename string) (string, ProductionError) {

	// Open file
	file, err := excelize.OpenFile(uploadPath)
	if err != nil {
		errorMsg := "Failed to open the file to process."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    PROCESS_ERROR_OPENING_FILE,
			ErrorMessage: errorMsg}
	}

	// Defer Close
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	// Get active Sheet
	activeSheet := file.GetSheetName(file.GetActiveSheetIndex())

	// Open new StreamWriter
	sw, err := file.NewStreamWriter(activeSheet)
	if err != nil {
		errorMsg := "Failed create the streamwriter."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    PROCESS_ERROR_CREATING_STREAMWRITER,
			ErrorMessage: errorMsg}
	}

	// Get the rows in a 2D matrix
	matrix, err := file.GetRows(activeSheet)
	if err != nil {
		errorMsg := "Failed to get the rows of the excel file."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    PROCESS_ERROR_GETTING_ROWS,
			ErrorMessage: errorMsg}
	}

	// Initialize stmtID , rowCoordinate and codedStatementColumn
	stmtId := "1"
	rowCoordinate := 1
	var codedStatementColumn int

	// Loop through the rows of the matrix
	for r, rowMatrix := range matrix {

		// Initialize interfaces:
		// Final interface to pass it to the streamwriter
		rowToWriteInterface := make([]interface{}, 0)
		// Interface to catch the parsing error from the parser
		errorMessageInterface := make([]interface{}, 1)
		//Interface to catch the existing row
		rowMatrixInterface := make([]interface{}, len(matrix[0]))

		// Copy to 2D matrix to interface
		for i, v := range rowMatrix {
			rowMatrixInterface[i] = v
		}

		// First iteration: printing the header by catching the existing row and adding
		//the header with a dummy statement
		if r == 0 {

			// Search for column of Coded Statement (err1 is of Type ProductionError)
			arrayCodedStatementColumn, err1 := SearchCodedStatementIdx(rowMatrix)
			if err1.ErrorCode != PRODUCTION_NO_ERROR {
				return "", err1
			}

			// Get the first and only element of the array indicating the column of the coded statement
			codedStatementColumn = arrayCodedStatementColumn[0]

			// Make dummy statement to print header
			dummyStatementToPrintHeader := "Cac{Once E(policy) F(comes into force)} A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance)."

			/*
				Get the output of the parsing: stmtId is not used, OutputTYpe is indifferent, no filename,
				no overwrite (because the file is not needed), printHeaders is indifferent, printIgScriptInput is indifferent.
				The only field used is HeaderNames to print the header
			*/

			output, _ := endpoints.ConvertIGScriptToTabularOutput(dummyStatementToPrintHeader, stmtId, tabular.OUTPUT_TYPE_CSV,
				"", false, tabular.IncludeHeader(), tabular.DEFAULT_IG_SCRIPT_OUTPUT)

			// Get the HeaderNames from the output
			headerArray := output[0].HeaderNames

			// Append the error column to the HeeaderNames
			headerArray = append([]string{"Error"}, headerArray...)

			// Copy the Error + HeaderNames to an interface
			headerInterface := make([]interface{}, len(headerArray))
			for i, v := range headerArray {
				headerInterface[i] = v
			}

			// Append the Error + HeaderNames to the existing header
			rowToWriteInterface = append(rowMatrixInterface, headerInterface...)

			// Get the coordinates of the cell to write (err1 is of Type error)
			coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
			if err != nil {
				errorMsg := "Failed to convert coordinates to cell name."
				log.Println(err.Error())
				return "", ProductionError{
					ErrorCode:    PROCESS_ERROR_COORDINATE_CONVERSION,
					ErrorMessage: errorMsg,
				}
			}

			// Set the row in the coordinate cell with the streamwriter
			if err := sw.SetRow(coordinateCell, rowToWriteInterface); err != nil {
				errorMsg := "Failed to set the row with the streamwriter."
				log.Println(err.Error())
				return "", ProductionError{
					ErrorCode:    PROCESS_ERROR_SETTING_ROW,
					ErrorMessage: errorMsg,
				}
			}
			// Incerment the rowCoordinate in order to pass from A1 to A2
			rowCoordinate += 1

		} else {

			/*
				Rest of iterations: getting the parsed output and appending it to the existing row
				Get the output of the parsing: stmtId is incremented(in the future there can
				be a function to assign columns to print an specific type of stmtID),
				OutputTYpe Ris indifferent, no filename and overwrite (because the file is not needed),
				printHeaders is indifferent, printIgScriptInput is indifferent.
			*/

			// Get the output of the parsing. err2 is of type tree.ParsingError
			output, err2 := endpoints.ConvertIGScriptToTabularOutput(rowMatrix[codedStatementColumn], stmtId, tabular.OUTPUT_TYPE_CSV,
				"", false, tabular.IncludeHeader(), tabular.DEFAULT_IG_SCRIPT_OUTPUT)

			// Catching the parsing error and store it in the interface. If there is a error, print the error and continue to the next row
			// If there is no error: print OK and print the output of the parsing
			if err2.ErrorCode != tree.PARSING_NO_ERROR {

				errorMessageInterface[0] = err2.ErrorMessage
				rowToWriteInterface = append(rowMatrixInterface, errorMessageInterface)

				// Get the coordinates of the cell to write
				coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
				if err != nil {
					errorMsg := "Failed to convert coordinates to cell name."
					log.Println(err.Error())
					return "", ProductionError{
						ErrorCode:    PROCESS_ERROR_COORDINATE_CONVERSION,
						ErrorMessage: errorMsg,
					}
				}
				// Set the row in the coordinate cell with the streamwriter
				if err := sw.SetRow(coordinateCell, rowToWriteInterface); err != nil {
					errorMsg := "Failed to set the row with the streamwriter."
					log.Println(err.Error())
					return "", ProductionError{
						ErrorCode:    PROCESS_ERROR_SETTING_ROW,
						ErrorMessage: errorMsg,
					}
				}
				rowCoordinate += 1
				continue
			}

			// Write OK in the error Column
			errorMessageInterface[0] = "OK"

			// Increment the stmtId
			stmtIdint, _ := strconv.Atoi(stmtId)
			stmtIdint += 1
			stmtId = strconv.Itoa(stmtIdint)

			// Loop to write each elemnt of the statement map
			for _, statementMap := range output[0].StatementMap {
				singleOutputRowInterface := make([]interface{}, len(output[0].HeaderSymbols))

				// Loop to run through the map through the headerSymbols and print the pair Key:value. if there is no value, print an empty string
				for i, headerSymbol := range output[0].HeaderSymbols {
					if val, ok := statementMap[headerSymbol]; ok {
						singleOutputRowInterface[i] = val
					} else {
						singleOutputRowInterface[i] = ""
					}
				}

				// Append the existing row to the error message + the output of the parsing
				rowToWriteInterface = append(rowMatrixInterface, errorMessageInterface)
				rowToWriteInterface = append(rowToWriteInterface, singleOutputRowInterface...)

				// Get the coordinates of the cell to write
				coordinateCell, err := excelize.CoordinatesToCellName(1, rowCoordinate)
				if err != nil {
					errorMsg := "Failed to convert coordinates to cell name."
					log.Println(err.Error())
					return "", ProductionError{
						ErrorCode:    PROCESS_ERROR_COORDINATE_CONVERSION,
						ErrorMessage: errorMsg,
					}
				}
				// Set the row in the coordinate cell with the streamwriter
				if err := sw.SetRow(coordinateCell, rowToWriteInterface); err != nil {
					errorMsg := "Failed to set the row with the streamwriter."
					log.Println(err.Error())
					return "", ProductionError{
						ErrorCode:    PROCESS_ERROR_SETTING_ROW,
						ErrorMessage: errorMsg,
					}
				}
				// Increment the rowCoordinate
				rowCoordinate += 1

			}
		}
	}

	// Flush
	if err := sw.Flush(); err != nil {
		errorMsg := "Failed to flush the streamwriter."
		log.Println(err.Error())
		return "", ProductionError{
			ErrorCode:    PROCESS_ERROR_FLUSHING_STREAMWRITER,
			ErrorMessage: errorMsg,
		}
	}

	savePath, err1 := SavExcelizeFile(file, filename)
	if err1.ErrorCode != PRODUCTION_NO_ERROR {
		return "", err1
	}

	return savePath, ProductionError{ErrorCode: PRODUCTION_NO_ERROR}
}
