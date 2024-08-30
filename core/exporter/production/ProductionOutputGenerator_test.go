package production

import (
	"IG-Parser/core/exporter/tabular"
	"os"
	"path/filepath"
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
			header:      []string{"Statement", "Coded", "Statement", "Code", "Coded Statement"},
			expectedIdx: []int{4},
			expectedErr: PRODUCTION_NO_ERROR,
		},
		{ // No element match the Regex
			header: []string{"Statement", "S. C.", "Code", "C. St.",
				"Co Sta", "Cod S", "Coded", "C Statement"},
			expectedIdx: nil,
			expectedErr: HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT,
		},
		{ // All this elemnts match the Regex
			header: []string{"Coded Statement", "Statement Coded",
				"CODED STATEMENT", "STATEMENT CODED", "Cod. St.", "Cod Statmnt",
				"CodedStatemnt", "Cod Stmnt", "Coded St", "Coded_-/Staement123"},
			expectedIdx: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
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

func TestProductionExcelWithoutParsingError(t *testing.T) {

	// IG Extended output: false
	tabular.SetProduceIGExtendedOutput(false)
	// IG Logical output: false
	tabular.SetIncludeAnnotations(false)
	// Choose the input file
	inputFilename := "TestProductionWithoutParsingError.xlsx"
	// Choose the expected output file
	expectedFilename := "TestProductionWithoutParsingError_CODED.xlsx"

	// Get path of the input file
	inputPath := filepath.Join(
		"testing", "input", inputFilename)
	inputPath, err := filepath.Abs(inputPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path of input file: %v", err)
	}
	// Ensure the output directory exists
	outputDir := filepath.Join(LIBRARY_DIRECTORY_NAME)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			t.Fatalf("Failed to create output directory: %v", err)
		}
	}

	//outputPath: IG-Library path + filename
	outputPath, err1 := ProcessExcelFile(inputPath, inputFilename)
	if err1.ErrorCode != PRODUCTION_NO_ERROR {
		t.Fatalf("ProcessExcelFile returned an error: %v", err1)
	}
	t.Logf("Output file path: %s", outputPath)

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
		t.Error(err)
	}

	// Clean up: remove the output file from IG-Library folder
	if err := os.Remove(outputPath); err != nil {
		t.Fatalf("Failed to remove output file: %v", err)
	}

	if err := os.Remove(outputDir); err != nil {
		t.Fatalf("Failed to remove output file: %v", err)
	}
}
