package production

import (
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/web/converter/shared"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestSearchCodedStatementIdx(t *testing.T) {
	tests := []struct {
		header      []string
		expectedIdx []int
		expectedErr string
	}{
		{ // Catching Coded Statement.
			header:      []string{"Hello", "Bye", "Hello", "Bye", "Coded Statement"},
			expectedIdx: []int{4},
			expectedErr: PRODUCTION_NO_ERROR,
		},
		{ // No element match the Regex
			header: []string{"Statement", "S. C.", "Code", "C. St.",
				"Co Sta", "Cod S", "Coded", "C Statement", "Encoded"},
			expectedIdx: nil,
			expectedErr: HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
		},
		{ // All this elemnts match the Regex
			header: []string{"Coded Statement", "Statement Coded",
				"CODED STATEMENT", "STATEMENT CODED", "Cod. St.", "Cod Statmnt",
				"CodedStatemnt", "Cod Stmnt", "Coded St", "Coded_-/Staement123",
				"Encoded Statement", "Statement Encoded", "ENCODED STATEMENT", "STATEMENT ENCODED"},
			expectedIdx: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			expectedErr: HEADER_MATCHING_ERROR_MULTIPLE_MATCHES_FOR_CODED_STATEMENT,
		},
	}

	for _, test := range tests {
		actualIdx, actualErr := SearchCodedStatementIdx(test.header)
		if !compareSlices(actualIdx, test.expectedIdx) || actualErr.ErrorCode != test.expectedErr {
			t.Errorf("For header %v, expected indexes %v and error code %s, but got indexes %v and error code %s",
				test.header, test.expectedIdx, test.expectedErr, actualIdx, actualErr.ErrorCode)
		}
	}
}

func TestProcessFileDefaultConfig(t *testing.T) {

	// IG Extended output: false
	tabular.SetProduceIGExtendedOutput(false)
	// IG Logical output: false
	tabular.SetIncludeAnnotations(false)

	tests := []struct {
		inputFilename    string
		expectedFilename string
		expectedErr      string
	}{
		{ // Without parsing error
			inputFilename:    "01_TestProductionWithoutParsingError.xlsx",
			expectedFilename: "01_TestProductionWithoutParsingError_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},

		{ // With parsing error
			inputFilename:    "02_TestProductionParsingError.xlsx",
			expectedFilename: "02_TestProductionParsingError_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},
		{ // Empty coded statement in 1 row having the same length as the header
			inputFilename:    "03_TestProductionEmptyCodedStatementSameRowLength.xlsx",
			expectedFilename: "03_TestProductionEmptyCodedStatementSameRowLength_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},

		{ // Empty coded statement in 1 row having a shorter length than the header
			inputFilename:    "04_TestProductionEmptyCodedStatementLessRowLength.xlsx",
			expectedFilename: "04_TestProductionEmptyCodedStatementLessRowLength_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},

		{ // Empty coded statement in 1 row having a shorter length than the header
			inputFilename:    "05_OneColumnOffset.xlsx",
			expectedFilename: "05_OneColumnOffset_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},
		{ // Empty rows, and parsing errors to test the dinamycally increment of stmtID
			inputFilename:    "06_TestProductionErrorsAndEmptyRows.xlsx",
			expectedFilename: "06_TestProductionErrorsAndEmptyRows_CODED.xlsx",
			expectedErr:      PRODUCTION_NO_ERROR,
		},

		{ // Header without Coded Statement column
			inputFilename:    "101_TestProductionEmptyCellHeaderCodedStatementNoMatch.xlsx",
			expectedFilename: "",
			expectedErr:      HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
		},

		{ // Row larger than header
			inputFilename:    "102_TestProductionRowLargerThanHeader.xlsx",
			expectedFilename: "",
			expectedErr:      PROCESS_ERROR_ROW_LARGER_THAN_HEADER,
		},
		{ // Matrix with offset of 1 row
			inputFilename:    "103_OneRowOffset_CODED.xlsx",
			expectedFilename: "",
			expectedErr:      HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
		},
	}

	// Ensure the output directory exists (IG-Library within the production folder)
	// This folder is removed when finishing the test
	outputDir := filepath.Join(LIBRARY_DIRECTORY_NAME)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			t.Fatalf("Failed to create output directory: %v", err)
		}
	}

	for _, test := range tests {
		// Get the variables of the iteration
		inputFilename := test.inputFilename
		expectedFilename := test.expectedFilename
		expectedErr := test.expectedErr

		// Get path of the input file
		inputPath := filepath.Join(
			"testing", "input", inputFilename)
		inputPath, err := filepath.Abs(inputPath)
		if err != nil {
			t.Fatalf("Failed to get absolute path of input file: %v", err)
		}

		//outputPath: IG-Library path + filename
		outputPath, err1 := ProcessExcelFile(inputPath, inputFilename)
		if err1.ErrorCode != PRODUCTION_NO_ERROR {
			if err1.ErrorCode != expectedErr {
				// Getting an error different than expected without outputfile
				t.Errorf("%v returned an unexpected error: %v."+
					shared.LINEBREAK+"Expected Error: %v", inputFilename, err1.ErrorCode, expectedErr)
				continue
			}
		} else {
			// PRODUCTION_NO_ERROR

			// Open the output file
			actualFile, err := excelize.OpenFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to open output file: %v", err)
			}

			// Defer Close
			defer func() {
				if err := actualFile.Close(); err != nil {
					t.Fatalf("Failed to close the actual file: %v", err)
					return
				}
			}()

			if expectedFilename == "" {
				// Ensure the error directory exists (Errors within the production folder)
				errorDir := filepath.Join("Errors")
				if _, err := os.Stat(errorDir); os.IsNotExist(err) {
					if err := os.MkdirAll(errorDir, os.ModePerm); err != nil {
						t.Fatalf("Failed to create Errors directory: %v", err)
					}
				}

				// Get the name, and append EROR to it
				// Ensure the output directory exists
				errorFileName := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename)) + "_ERROR.xlsx"
				errorPath := filepath.Join(
					errorDir, errorFileName)

				// Open the error file
				errorFile, err := excelize.OpenFile(outputPath)
				if err != nil {
					t.Fatalf("Failed to open output file: %v", err)
				}
				// Defer Close
				defer func() {
					if err := errorFile.Close(); err != nil {
						t.Fatalf("Failed to close the error file: %v", err)
						return
					}
				}()

				// Save the error file
				if err := errorFile.SaveAs(errorPath); err != nil {
					t.Fatalf("Failed to save error file: %v", err)
				}

				// Clean up: remove the output with CODED name from IG-Library
				// in the testing folder
				if err := os.Remove(outputPath); err != nil {
					t.Fatalf("Failed to remove output file: %v", err)
				}

				t.Errorf("%v returned an unexpected file saved in %v. Expected Error: %v", inputFilename, errorPath, expectedErr)

			} else {

				// Get path of the expected file
				expectedPath := filepath.Join(
					"testing", "expected", expectedFilename)
				expectedPath, err2 := filepath.Abs(expectedPath)
				if err2 != nil {
					t.Fatalf("Failed to get absolute path of input file: %v", err)
				}

				// Open the output file
				expectedFile, err := excelize.OpenFile(expectedPath)
				if err != nil {
					t.Fatalf("Failed to open expected expected file: %v", err)
				}

				// Defer Close
				defer func() {
					if err := expectedFile.Close(); err != nil {
						t.Fatalf("Failed to close the expected file: %v", err)
						return
					}
				}()

				// Compare the contents of the expected and actual output files
				if err := compareExcelFiles(expectedFile, actualFile); err != nil {
					// Ensure the error directory exists (Errors within the production folder)
					errorDir := filepath.Join("Errors")
					if _, err := os.Stat(errorDir); os.IsNotExist(err) {
						if err := os.MkdirAll(errorDir, os.ModePerm); err != nil {
							t.Fatalf("Failed to create Errors directory: %v", err)
						}
					}

					// Get the name, and append EROR to it
					// Ensure the output directory exists
					errorFileName := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename)) + "_ERROR.xlsx"
					errorPath := filepath.Join(
						errorDir, errorFileName)

					// Open the error file
					errorFile, err := excelize.OpenFile(outputPath)
					if err != nil {
						t.Fatalf("Failed to open output file: %v", err)
					}
					// Defer Close
					defer func() {
						if err := errorFile.Close(); err != nil {
							t.Fatalf("Failed to close the error file: %v", err)
							return
						}
					}()

					// Save the error file
					if err := errorFile.SaveAs(errorPath); err != nil {
						t.Fatalf("Failed to save error file: %v", err)
					}

					// Clean up: remove the output with CODED name from IG-Library
					// in the testing folder
					if err := os.Remove(outputPath); err != nil {
						t.Fatalf("Failed to remove output file: %v", err)
					}

					t.Errorf("%v output file does not match expected file: %v."+
						shared.LINEBREAK+" Actual file saved in: %v.", inputFilename, expectedFilename, outputPath)

				}
			}

		}
	}
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("Failed to remove output file: %v", err)
	}
}

// Test the upload of the file. Supported extensions are: xlsx and xlsm
// No test to upload a big file larger than maximum supported size (5 MB by default)
// Check ProductionOutputGeneratorConfig.go
func TestUploadExcelFile(t *testing.T) {
	tests := []struct {
		inputFilename string
		expectedErr   string
	}{
		{
			inputFilename: "01_TestProductionWithoutParsingError.xlsx",
			expectedErr:   PRODUCTION_NO_ERROR,
		},

		{
			inputFilename: "200_TestProductionNotExcelFile.docx",
			expectedErr:   UPLOAD_ERROR_NOT_EXCEL_FILE,
		},
		{
			inputFilename: "201_TestProductionNotExcelFile.xls",
			expectedErr:   UPLOAD_ERROR_NOT_EXCEL_FILE,
		},
		{
			inputFilename: "203_TestProductionXLSMFile.xlsm",
			expectedErr:   PRODUCTION_NO_ERROR,
		},
	}

	for _, test := range tests {

		inputFilename := test.inputFilename
		expectedErr := test.expectedErr

		// Get path of the input file
		inputPath := filepath.Join(
			"testing", "input", inputFilename)
		inputPath, err := filepath.Abs(inputPath)
		if err != nil {
			t.Fatalf("Failed to get absolute path of input file: %v", err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/upload", nil)
		if err != nil {
			t.Fatalf("Failed to create HTTP request: %v", err)
		}

		// Create a new multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create a new file part
		fileWriter, err := writer.CreateFormFile("file", inputFilename)
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		// Open the test file
		file, err := os.Open(inputPath)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer file.Close()

		// Copy the file contents to the file part
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			t.Fatalf("Failed to copy file contents: %v", err)
		}

		// Close the multipart writer
		err = writer.Close()
		if err != nil {
			t.Fatalf("Failed to close multipart writer: %v", err)
		}

		// Set the request body and content type
		req.Body = io.NopCloser(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Call the UploadExcelFile function
		_, _, err1 := UploadExcelFile(req)
		if err1.ErrorCode != expectedErr {
			t.Errorf("UploadExcelFile returned an error: %v", err1)
		}

	}
}
