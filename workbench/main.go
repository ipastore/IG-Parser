package main

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {

	//Open File

	f, err := excelize.OpenFile("/Users/ignaciopastorebenaim/go/src/IG-Parser/uploads/200_MAX.xlsx")
	if err != nil {
		fmt.Println("Error: open file")
	}

	// Defer Close
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	//Here I could get the name of the active sheet and then getrows
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println("Error: GetRows")
	}

	// New File to populate
	f1 := excelize.NewFile()

	// Open Stream Writer to populate Sheet1
	streamWriter, err := f1.NewStreamWriter("Sheet1")
	if err != nil {
		fmt.Println("Error: NewStreamWriter")
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
		fmt.Printf("Check 72 %v", i)
		codedStmt := rows[i][8]

		//Create Tabular Output
		output, err0 := endpoints.ConvertIGScriptToTabularOutput(codedStmt, stmtId, tabular.OUTPUT_TYPE_CSV, "", false, false, tabular.DEFAULT_IG_SCRIPT_OUTPUT)

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

	err90 := f1.SaveAs("/Users/ignaciopastorebenaim/RESILIENT_IG-Parser/coded/200_MAX_CODED.xlsx")
	if err90 != nil {
		fmt.Println(err90)

	}
	fmt.Println("Successfully written as: " + "/Users/ignaciopastorebenaim/RESILIENT_IG-Parser/coded/200_MAX_CODED.xlsx")

}
