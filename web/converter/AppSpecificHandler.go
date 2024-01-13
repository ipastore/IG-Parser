package converter

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"IG-Parser/web/converter/shared"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

/*
Third-level handler generating Production output in response to web request.
Should be invoked by #converterHandler().
*/

type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}
	fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}

// Esta funcion la tengo que modificar
func handleProductionOutput(w http.ResponseWriter, r *http.Request, retStruct shared.ReturnStruct, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, outputType string, printHeaders bool, printIgScriptInput string) {
	//DE ACA

	// Run default configuration
	shared.SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	Println("Setting dynamic output:", dynamicOutput)
	tabular.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	Println("Setting IG Extended output:", produceIGExtendedOutput)
	tabular.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	Println("Setting annotations:", includeAnnotations)
	tabular.SetIncludeAnnotations(includeAnnotations)
	// Define whether header row is included
	Println("Setting header row:", printHeaders)
	tabular.SetIncludeHeaders(printHeaders)
	// Output type
	Println("Output type:", outputType)

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

	// Progress is used to track the progress of a file upload.
	// It implements the io.Writer interface so it can be passed
	// to an io.TeeReader()

	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]

	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.MkdirAll("./coded", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// filetype := http.DetectContentType(buff)
		// if filetype != "image/jpeg" && filetype != "image/png" {
		// 	http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		// 	return
		// }

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f0, err := os.Create("./uploads/" + fileHeader.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f0.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f0, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filename := fileHeader.Filename
		saveAsFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + "_CODED.xlsx"

		// /IG-parser if running from executable
		joinpath := filepath.Join(
			// "IG-parser",
			"uploads", filename)

		joinSaveAsPath := filepath.Join(
			// "IG-parser",
			"coded", saveAsFilename)

		abspath, err101 := filepath.Abs(joinpath)
		if err101 != nil {
			fmt.Println(err101)
		}

		absSaveAsPath, err109 := filepath.Abs(joinSaveAsPath)
		if err109 != nil {
			fmt.Println(err109)
		}

		f, err34 := excelize.OpenFile(abspath)
		if err34 != nil {
			fmt.Println(err34)
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = err34.Error()

			err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_PRODUCTION, retStruct)
			if err != nil {
				log.Println("Error processing default template:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Final comment in log
			Println("Error: " + fmt.Sprint(err34))
			// Ensure logging is terminated
			err3 := terminateOutput(ERROR_SUFFIX)
			if err3 != nil {
				log.Println("Error when finalizing log file: ", err3.Error())
			}
			return
		}

		// Defer Close
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		// EXCELIZE WORK

		// SheetName := f.GetSheetName(1)

		// Get all the rows in the Sheet1.
		rows, err40 := f.GetRows("Sheet1")
		if err40 != nil {
			fmt.Println(err40)
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = err40.Error()

			err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_PRODUCTION, retStruct)
			if err != nil {
				log.Println("Error processing default template:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Final comment in log
			Println("Error: " + fmt.Sprint(err40))
			// Ensure logging is terminated
			err3 := terminateOutput(ERROR_SUFFIX)
			if err3 != nil {
				log.Println("Error when finalizing log file: ", err3.Error())
			}
			return
		}

		// New File to populate
		f1 := excelize.NewFile()

		// Open Stream Writer to populate Sheet1
		streamWriter, err35 := f1.NewStreamWriter("Sheet1")

		if err35 != nil {
			fmt.Print(err35)
			return
		}

		// Header: Populate string for first row
		headerString := "Statement ID,Attributes,Attributes Property,Attributes Property Reference,Deontic,Aim,Direct Object,Direct Object Reference," +
			"Direct Object Property,Direct Object Property Reference,Indirect Object,Indirect Object Reference,Indirect Object Property,Indirect Object Property Reference," +
			"Activation Condition,Activation Condition Reference,Execution Constraint,Execution Constraint Reference,Constituted Entity,Constituted Entity Property," +
			"Constituted Entity Property Reference,Modal,Constitutive Function,Constituting Properties,Constituting Properties Reference,Constituting Properties Properties," +
			"Constituting Properties Properties Reference,Or Else Reference,Logical Linkage (Statements),Logical Linkage (Components)"

		arrayOfHeader := strings.Split(headerString, ",")

		arrayOfHeaderInterface := make([]interface{}, len(arrayOfHeader))
		for i, v := range arrayOfHeader {
			arrayOfHeaderInterface[i] = v
		}

		// Write Header
		if err := streamWriter.SetRow("A1", arrayOfHeaderInterface); err != nil {
			fmt.Println(err)
			return
		}

		// Here I should run through the matrix and keep the coded text to give it to the parser
		// UPGRADE: I should append a call with nil or "" values. If not the length of row it´s 7 instead of 8
		rowindexExcel := 2
		colindexExcel := "A"
		stmtId := "1"

		for i := 1; i < len(rows); i++ {
			codedStmt := rows[i][8]

			//Create Tabular Output
			output, err0 := endpoints.ConvertIGScriptToTabularOutput(codedStmt, stmtId, tabular.OUTPUT_TYPE_CSV, "", false, false, printIgScriptInput)

			// Getting the error and here I should paste it
			if err0.ErrorCode != tree.PARSING_NO_ERROR {
				fmt.Println(err0.ErrorCode)
			}

			// Spliting Ouptut into rows. The first element of outputRows is 0 or nil. The len of outputRows is 3 so later
			// I must correct the i index. rowindexExcel and colindexExcel are the coordinates from where I should populate the sheet

			outputRows := strings.Split(output[0].Output, tabular.StmtIdPrefix) //Be aware here with the apostrophe
			outputMatrix := make([][]string, len(outputRows)-1)
			outputSingleInterface := make([]interface{}, len(outputRows[1]))

			for j := 0; j < len(outputRows)-1; j++ {
				outputMatrix[j] = strings.Split(outputRows[j+1], tabular.CellSeparator)
				for k, v := range outputMatrix[j] {
					outputSingleInterface[k] = v

				}

				err := streamWriter.SetRow(colindexExcel+strconv.Itoa(rowindexExcel), outputSingleInterface)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					rowindexExcel += 1
				}
			}
			stmtIdint, _ := strconv.Atoi(stmtId)
			stmtIdint += 1
			stmtId = strconv.Itoa(stmtIdint)
		}

		// FLush
		if err := streamWriter.Flush(); err != nil {
			fmt.Println(err)
			return
		}

		// Save (I should put a success message). I also should get the number of the interview to pass it to the string path
		// AND ERASE TEMP FILE (/uploads)

		err90 := f1.SaveAs(absSaveAsPath)
		if err90 != nil {
			fmt.Println(err90)

		}
		fmt.Println("Successfully written as: " + absSaveAsPath)
		retStruct.Success = true
		retStruct.Error = false
		retStruct.Message = "Successfully saved as " + absSaveAsPath
		// El loop llega hasta ACA y finalmente se ejecuta el template con el retStruct que deberia ser un mensaje de success
		err = tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_PRODUCTION, retStruct)
		if err != nil {
			log.Println("Error processing default template:", err.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		Println("Success")
		// Ensure logging is terminated
		err3 := terminateOutput(SUCCESS_SUFFIX)
		if err3 != nil {
			log.Println("Error when finalizing log file: ", err3.Error())
		}
	}
	return
}

/*
Third-level handler generating tabular output in response to web request.
Should be invoked by #converterHandler().
*/
func handleTabularOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct shared.ReturnStruct, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, outputType string, printHeaders bool, printIgScriptInput string) {
	// Run default configuration
	shared.SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	Println("Setting dynamic output:", dynamicOutput)
	tabular.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	Println("Setting IG Extended output:", produceIGExtendedOutput)
	tabular.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	Println("Setting annotations:", includeAnnotations)
	tabular.SetIncludeAnnotations(includeAnnotations)
	// Define whether header row is included
	Println("Setting header row:", printHeaders)
	tabular.SetIncludeHeaders(printHeaders)
	// Indicate whether IG Script input is included in output
	Println("Include IG Script input in generated output:", printIgScriptInput)
	// Output type
	Println("Output type:", outputType)
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToTabularOutput(codedStmt, stmtId, outputType, "", true, tabular.IncludeHeader(), printIgScriptInput)
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		// Deal with potential errors and prepopulate return message
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		// Execute template
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_TABULAR, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		Println("Error: " + fmt.Sprint(err2))
		// Ensure logging is terminated
		err := terminateOutput(ERROR_SUFFIX)
		if err != nil {
			log.Println("Error when finalizing log file: ", err.Error())
		}
		return
	}
	// Return success if parsing was successful
	retStruct.Success = true
	retStruct.CodedStmt = codedStmt
	tabularOutput := ""
	for _, v := range output {
		tabularOutput += v.Output
	}
	retStruct.Output = tabularOutput
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_TABULAR, retStruct)
	if err != nil {
		log.Println("Error processing default template:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

	// Final comment in log
	Println("Success")
	// Ensure logging is terminated
	err3 := terminateOutput(SUCCESS_SUFFIX)
	if err3 != nil {
		log.Println("Error when finalizing log file: ", err3.Error())
	}
	return
}

/*
Third-level handler generating visual tree output in response to web request.
Should be invoked by #converterHandler().
*/
func handleVisualOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct shared.ReturnStruct, flatOutput bool, binaryOutput bool, moveActivationConditionsToTop bool, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, includeDoV bool) {
	// Run default configuration
	shared.SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	Println("Setting dynamic output:", dynamicOutput)
	tabular.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	Println("Setting IG Extended output:", produceIGExtendedOutput)
	tabular.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	Println("Setting annotations:", includeAnnotations)
	tabular.SetIncludeAnnotations(includeAnnotations)
	// Define whether Degree of Variability is included
	Println("Setting Degree of Variability (DoV):", includeDoV)
	tabular.SetIncludeDegreeOfVariability(includeDoV)
	// Setting flat printing
	Println("Setting flat printing of properties:", flatOutput)
	tree.SetFlatPrinting(flatOutput)
	Println("Setting binary tree printing:", binaryOutput)
	tree.SetBinaryPrinting(binaryOutput)
	Println("Setting activation condition on top in visual output:", moveActivationConditionsToTop)
	tree.SetMoveActivationConditionsToFront(moveActivationConditionsToTop)
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToVisualTree(codedStmt, stmtId, "")
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_VISUAL, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		Println("Error: " + fmt.Sprint(err2))
		// Ensure logging is terminated
		err := terminateOutput(ERROR_SUFFIX)
		if err != nil {
			log.Println("Error when finalizing log file: ", err.Error())
		}
		return
	}
	// Return success if parsing was successful
	retStruct.Success = true
	retStruct.CodedStmt = codedStmt
	retStruct.Output = output
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_VISUAL, retStruct)
	if err != nil {
		log.Println("Error processing default template:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

	// Final comment in log
	Println("Success")
	// Ensure logging is terminated
	err3 := terminateOutput(SUCCESS_SUFFIX)
	if err3 != nil {
		log.Println("Error when finalizing log file: ", err3.Error())
	}
	return
}
