package production

import (
	"errors"
)

// Indicates founding no matches in the input header for "Coded Statement"
const HEADER_MATCHING_ERROR_NO_MATCH_FOR_CODED_STATEMENT = "NO_MATCH_FOR_CODED_STATEMENT"

// Indicates founding no matches in the input header for "Coded Statement"
const HEADER_MATCHING_ERROR_MULTIPLE_MATCHES_FOR_CODED_STATEMENT = "MULTIPLE_MATCHES_FOR_CODED_STATEMENT"

// Indicates exceeding the maximum upload size
const UPLOAD_ERROR_FILE_TOO_BIG = "FILE_TOO_BIG"

// Indicates the file is not an excel file
const UPLOAD_ERROR_NOT_EXCEL_FILE = "NOT_EXCEL_FILE"

// Indicates error in parsing the multipart form
const UPLOAD_ERROR_PARSING_FORM = "ERROR_PARSING_FORM"

// Indicates error in getting the file from the form
const UPLOAD_ERROR_GETTING_FILE = "ERROR_GETTING_FILE"

// Indicates error in creating the temp file
const UPLOAD_ERROR_CREATING_FILE = "ERROR_CREATING_FILE"

// Indicates error in copying the file
const UPLOAD_ERROR_COPYING_FILE = "ERROR_COPYING_FILE"

// Indicates no error during the production process
const PRODUCTION_NO_ERROR = "PRODUCTION_NO_ERROR"

// Indicates error in creating the folder whether in upload or saving
const UPLOAD_SAVE_ERROR_CREATING_TEMP_FOLDER = "ERROR_CREATING_TEMP_FOLDER"

// Indicates error in saving the file
const SAVE_ERROR_SAVING_FILE = "ERROR_SAVING_FILE"

// Indicates error in removing the file from uploads
const REMOVE_ERROR_ERASING_FILE = "ERROR_ERASING_FILE"

// Indicates error in opening file to process
const PROCESS_ERROR_OPENING_FILE = "ERROR_OPENING_FILE"

// Indicates error in closing the file after processing
const PROCESS_ERROR_CLOSING_FILE = "ERROR_CLOSING_FILE"

// Indicaates error in creating the streamwriter
const PROCESS_ERROR_CREATING_STREAMWRITER = "ERROR_CREATING_STREAMWRITER"

// Indicates error in getting the rows from the excel file
const PROCESS_ERROR_GETTING_ROWS = "ERROR_GETTING_ROWS"

// Indicates error in converting to row and column coordinates in the excelize package
const PROCESS_ERROR_COORDINATE_CONVERSION = "ERROR_COORDINATE_CONVERSION"

// Indicates error in setting the row in the excel file with the streamwriter
const PROCESS_ERROR_SETTING_ROW = "ERROR_SETTING_ROW"

// Indicates error in flushing the streamwriter
const PROCESS_ERROR_FLUSHING_STREAMWRITER = "ERROR_FLUSHING_STREAMWRITER"

/*
Error type signaling errors during the upload process
*/
type ProductionError struct {
	ErrorCode    string
	ErrorMessage string
}

func (e *ProductionError) Error() error {
	return errors.New("Production Error " + e.ErrorCode + ": " + e.ErrorMessage)
}
