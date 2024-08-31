package converter

import (
	"IG-Parser/core/config"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/web/converter/shared"
	"IG-Parser/web/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
Generic handler for pre-processing of web application input.
- Invoked by web handlers provided in Handler.go.
- Invokes output-specific handlers (AppSpecificHandler.go).
*/

/*
Second-level general handler that retrieves and preprocesses information from input.
Delegates to third-order handler for output-specific generation.
Should be invoked by #ConverterHandlerTabular() and #ConverterHandlerVisual().
*/
func converterHandler(w http.ResponseWriter, r *http.Request, templateName string) {

	//// STEP 1: Read all parameters from returned form

	// Reading form to prepopulate response
	message := ""
	transactionID := ""
	formValueRawStmt := r.FormValue(shared.PARAM_RAW_STATEMENT)
	formValueCodedStmt := r.FormValue(shared.PARAM_CODED_STATEMENT)
	formValueStmtId := r.FormValue(shared.PARAM_STATEMENT_ID)
	formValueDynamicOutput := r.FormValue(shared.PARAM_DYNAMIC_SCHEMA)
	formValueIncludeAnnotations := r.FormValue(shared.PARAM_LOGICO_OUTPUT)
	formValueIncludeDoV := r.FormValue(shared.PARAM_DOV)
	formValueIgExtendedOutput := r.FormValue(shared.PARAM_EXTENDED_OUTPUT)
	formValueIncludeHeaders := r.FormValue(shared.PARAM_PRINT_HEADERS)
	formValuePrintOriginalStatement := r.FormValue(shared.PARAM_PRINT_ORIGINAL_STATEMENT)
	formValuePrintIgScript := r.FormValue(shared.PARAM_PRINT_IG_SCRIPT)
	formValueOutputType := r.FormValue(shared.PARAM_OUTPUT_TYPE)
	formValuePropertyTree := r.FormValue(shared.PARAM_PROPERTY_TREE)
	formValueBinaryTree := r.FormValue(shared.PARAM_BINARY_TREE)
	formValueMoveActivationConditionsToTop := r.FormValue(shared.PARAM_ACTIVATION_CONDITION_ON_TOP)
	formValueCanvasHeightValue := r.FormValue(shared.PARAM_HEIGHT)
	formValueCanvasWidthValue := r.FormValue(shared.PARAM_WIDTH)

	// EVALUATE INDIVIDUAL CHECKBOX INPUTS

	// Dynamic output
	dynamicOutput := false
	Println("Form field (tabular) - Dynamic: ", formValueDynamicOutput)
	if formValueDynamicOutput == shared.CHECKBOX_ON {
		formValueDynamicOutput = shared.CHECKBOX_CHECKED
		dynamicOutput = true
	} else {
		formValueDynamicOutput = shared.CHECKBOX_UNCHECKED
		dynamicOutput = false
	}

	// Annotations in output
	includeAnnotations := false
	Println("Form field (both)    - Annotations: ", formValueIncludeAnnotations)
	if formValueIncludeAnnotations == shared.CHECKBOX_ON {
		formValueIncludeAnnotations = shared.CHECKBOX_CHECKED
		includeAnnotations = true
	} else {
		formValueIncludeAnnotations = shared.CHECKBOX_UNCHECKED
		includeAnnotations = false
	}

	// DoV in output
	includeDoV := false
	Println("Form field (visual)  - DoV: ", formValueIncludeDoV)
	if formValueIncludeDoV == shared.CHECKBOX_ON {
		formValueIncludeDoV = shared.CHECKBOX_CHECKED
		includeDoV = true
	} else {
		formValueIncludeDoV = shared.CHECKBOX_UNCHECKED
		includeDoV = false
	}

	// Component-level nesting in output
	produceIGExtendedOutput := false
	Println("Form field (tabular) - IG Extended output: ", formValueIgExtendedOutput)
	if formValueIgExtendedOutput == shared.CHECKBOX_ON {
		formValueIgExtendedOutput = shared.CHECKBOX_CHECKED
		produceIGExtendedOutput = true
	} else {
		formValueIgExtendedOutput = shared.CHECKBOX_UNCHECKED
		produceIGExtendedOutput = false
	}

	// Print headers in output
	printHeaders := false
	// If not received by POST, set headers as default setting
	if formValueIncludeHeaders == "" && r.Method != http.MethodPost {
		formValueIncludeHeaders = shared.CHECKBOX_ON
	}
	Println("Form field (tabular) - Include headers in output: ", formValueIncludeHeaders)
	if formValueIncludeHeaders == shared.CHECKBOX_ON {
		formValueIncludeHeaders = shared.CHECKBOX_CHECKED
		printHeaders = true
	} else {
		formValueIncludeHeaders = shared.CHECKBOX_UNCHECKED
		printHeaders = false
	}

	// Selection for inclusion of Original Statement in output
	// If not received by POST, set Original Statement output as default setting
	if formValuePrintOriginalStatement == "" && r.Method != http.MethodPost {
		formValuePrintOriginalStatement = tabular.DEFAULT_ORIGINAL_STATEMENT_OUTPUT
	}
	Println("Form field (tabular) - Include Original Statement in output: ", formValuePrintOriginalStatement)

	// Selection for inclusion of IG Script in output
	// If not received by POST, set IG Script output as default setting
	if formValuePrintIgScript == "" && r.Method != http.MethodPost {
		formValuePrintIgScript = tabular.DEFAULT_IG_SCRIPT_OUTPUT
	}
	Println("Form field (tabular) - Include IG Script-encoded Statement in output: ", formValuePrintIgScript)

	// Private property printing in output
	printFlatProperties := false
	Println("Form field (visual)  - Private property tree printing: ", formValuePropertyTree)
	if formValuePropertyTree == shared.CHECKBOX_ON {
		formValuePropertyTree = shared.CHECKBOX_CHECKED
		printFlatProperties = false
	} else {
		formValuePropertyTree = shared.CHECKBOX_UNCHECKED
		printFlatProperties = true
	}

	// Binary tree printing in output
	printBinaryTree := false
	Println("Form field (visual)  - Binary tree printing: ", formValueBinaryTree)
	if formValueBinaryTree == shared.CHECKBOX_ON {
		formValueBinaryTree = shared.CHECKBOX_CHECKED
		printBinaryTree = true
	} else {
		formValueBinaryTree = shared.CHECKBOX_UNCHECKED
		printBinaryTree = false
	}

	// Activation condition on top in output
	printActivationConditionsOnTop := false
	Println("Form field (visual)  - Activation conditions on top: ", formValueMoveActivationConditionsToTop)
	if formValueMoveActivationConditionsToTop == shared.CHECKBOX_ON {
		formValueMoveActivationConditionsToTop = shared.CHECKBOX_CHECKED
		printActivationConditionsOnTop = true
	} else {
		formValueMoveActivationConditionsToTop = shared.CHECKBOX_UNCHECKED
		printActivationConditionsOnTop = false
	}

	// Checkbox interpretation finished

	// Prepare return structure with prepopulated information (to be refined prior to return)
	retStruct := shared.ReturnStruct{
		Success:                         false,
		Error:                           false,
		Message:                         message,
		OverrideSavedStmts:              false,
		RawStmt:                         formValueRawStmt,
		DefaultRawStmt:                  shared.RAW_STATEMENT,
		CodedStmt:                       formValueCodedStmt,
		DefaultCodedStmt:                shared.ANNOTATED_STATEMENT,
		StmtId:                          formValueStmtId,
		DefaultStmtId:                   shared.STATEMENT_ID,
		DynamicOutput:                   formValueDynamicOutput,
		IGExtendedOutput:                formValueIgExtendedOutput,
		IncludeAnnotations:              formValueIncludeAnnotations,
		IncludeDoV:                      formValueIncludeDoV,
		IncludeHeaders:                  formValueIncludeHeaders,
		PrintOriginalStatement:          formValuePrintOriginalStatement,
		PrintOriginalStatementSelection: tabular.ORIGINAL_STATEMENT_INCLUSION_OPTIONS,
		PrintIgScript:                   formValuePrintIgScript,
		PrintIgScriptSelection:          tabular.IG_SCRIPT_INCLUSION_OPTIONS,
		OutputType:                      formValueOutputType,
		OutputTypes:                     tabular.OUTPUT_TYPES,
		PrintPropertyTree:               formValuePropertyTree,
		PrintBinaryTree:                 formValueBinaryTree,
		ActivationConditionsOnTop:       formValueMoveActivationConditionsToTop,
		Width:                           shared.WIDTH,
		DefaultWidth:                    shared.WIDTH,
		Height:                          shared.HEIGHT,
		DefaultHeight:                   shared.HEIGHT,
		TransactionId:                   transactionID,
		IGScriptLink:                    shared.HEADER_SCRIPT_LINK,
		IGWebsiteLink:                   shared.HEADER_IG_LINK,
		RawStmtHelp:                     shared.HELP_RAW_STMT,
		CodedStmtHelpRef:                shared.HELP_REF,
		CodedStmtHelp:                   template.HTML(strings.Replace(shared.HELP_CODED_STMT, "\n", "<br>", -1)),
		StmtIdHelp:                      shared.HELP_STMT_ID,
		ParametersHelp:                  shared.HELP_PARAMETERS,
		OutputTypeHelp:                  shared.HELP_OUTPUT_TYPE,
		OriginalStatementInclusionHelp:  shared.HELP_ORIGINAL_STATEMENT_OUTPUT,
		IgScriptInclusionHelp:           shared.HELP_IG_SCRIPT_OUTPUT,
		ReportHelp:                      shared.HELP_REPORT,
		Version:                         config.IG_PARSER_VERSION}

	// Parse UI canvas information (visual parser)

	// Assign width for UI rendering
	if formValueCanvasWidthValue != "" {
		widthVal, err := strconv.Atoi(formValueCanvasWidthValue)
		if err != nil || widthVal < shared.MIN_WIDTH {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = shared.ERROR_INPUT_WIDTH
			err2 := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err2 != nil {
				log.Println("Error generating error response for template processing:", err2.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Stop execution and return to UI
			return
		}
		// Assign input value only in case of success
		retStruct.Width = widthVal
	}

	// Assign height for UI rendering
	if formValueCanvasHeightValue != "" {
		heightVal, err := strconv.Atoi(formValueCanvasHeightValue)
		if err != nil || heightVal < shared.MIN_HEIGHT {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = shared.ERROR_INPUT_HEIGHT
			err2 := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err2 != nil {
				log.Println("Error generating error response for template processing:", err2.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Stop execution and return to UI
			return
		}
		// Assign input value only in case of success
		retStruct.Height = heightVal
	}

	// All form parameters are processed ...

	//// STEP 2: Checking for URL parameters (that may override/refine specific returnStruct fields)

	if r.Method != http.MethodPost {

		// Request will generally be a GET request, but contain URL parameters

		// Just show empty form with prepopulated elements (to ensure they are reset to default values if not parameterized)
		retStruct.RawStmt = shared.RAW_STATEMENT
		retStruct.CodedStmt = shared.ANNOTATED_STATEMENT
		retStruct.StmtId = shared.STATEMENT_ID

		// CHECK FOR URL PARAMETERS TO CUSTOMIZE EXECUTION
		// NOTE: When using URL parameters, ensure to substitute all semicolons in URL parameters with %3B

		// Set switch to indicate potential need to align raw and coded statement field entries
		resetValues := false

		// Parameter: Raw Statement
		val, suc := extractUrlParameters(r, shared.PARAM_RAW_STATEMENT)
		if suc {
			retStruct.RawStmt = val
			resetValues = true
			// Override browser-saved content
			retStruct.OverrideSavedStmts = true
		}

		// Parameter: IG Script-coded statement - consider interaction with raw statement
		val, suc = extractUrlParameters(r, shared.PARAM_CODED_STATEMENT)
		if suc {
			retStruct.CodedStmt = val
			// Check for raw statement if it is still default, while the coded stmt changed; if so, then reset
			if retStruct.CodedStmt != shared.ANNOTATED_STATEMENT && retStruct.RawStmt == shared.RAW_STATEMENT {
				retStruct.RawStmt = ""
			}
			// Override browser-saved content
			retStruct.OverrideSavedStmts = true
		} else if resetValues {
			// Reset value, since the default coded statement will likely not correspond.
			retStruct.CodedStmt = ""
		}

		// TABULAR OUTPUT PARAMETERS

		// Parameter: Statement ID
		val, suc = extractUrlParameters(r, shared.PARAM_STATEMENT_ID)
		if suc {
			retStruct.StmtId = val
			// Override browser-saved content
			retStruct.OverrideSavedStmts = true
		}

		// NOTE: The following parameters would be overridden by browser-saved values where existing!

		// Parameter: Dynamic output
		val, suc = extractUrlParameters(r, shared.PARAM_DYNAMIC_SCHEMA)
		check := evaluateBooleanUrlParameters(shared.PARAM_DYNAMIC_SCHEMA, val, suc)
		// Assign values
		if check {
			retStruct.DynamicOutput = shared.CHECKBOX_CHECKED
			dynamicOutput = true
		} else {
			retStruct.DynamicOutput = shared.CHECKBOX_UNCHECKED
			dynamicOutput = false
		}

		// Parameter: Component-level nesting (IG Extended output)
		val, suc = extractUrlParameters(r, shared.PARAM_EXTENDED_OUTPUT)
		check = evaluateBooleanUrlParameters(shared.PARAM_EXTENDED_OUTPUT, val, suc)
		// Assign values
		if check {
			retStruct.IGExtendedOutput = shared.CHECKBOX_CHECKED
			produceIGExtendedOutput = true
		} else {
			retStruct.IGExtendedOutput = shared.CHECKBOX_UNCHECKED
			produceIGExtendedOutput = false
		}

		// Parameter: IG Logico annotations
		val, suc = extractUrlParameters(r, shared.PARAM_LOGICO_OUTPUT)
		check = evaluateBooleanUrlParameters(shared.PARAM_LOGICO_OUTPUT, val, suc)
		// Assign values
		if check {
			retStruct.IncludeAnnotations = shared.CHECKBOX_CHECKED
			includeAnnotations = true
		} else {
			retStruct.IncludeAnnotations = shared.CHECKBOX_UNCHECKED
			includeAnnotations = false
		}

		// Parameter: Header row printing
		val, suc = extractUrlParameters(r, shared.PARAM_PRINT_HEADERS)
		check = evaluateBooleanUrlParameters(shared.PARAM_PRINT_HEADERS, val, suc)
		// Sets default if no information is passed along in form (e.g., deactivation of headers);
		// Note that formValueIncludeHeaders will be prepopulated based on earlier form processing
		if formValueIncludeHeaders != shared.CHECKBOX_ON {

			// Fall back and assess whether the URL contained relevant parameter
			if !suc {
				check = true
			}
			// Assign values
			if check {
				retStruct.IncludeHeaders = shared.CHECKBOX_CHECKED
				printHeaders = true
			} else {
				retStruct.IncludeHeaders = shared.CHECKBOX_UNCHECKED
				printHeaders = false
			}
		}

		// Parameter: Original Statement inclusion
		val, suc = extractUrlParameters(r, shared.PARAM_PRINT_ORIGINAL_STATEMENT)
		if val != "" {
			// Read from parameter
			retStruct.PrintOriginalStatement = val
		} else {
			// Use default parameter
			retStruct.PrintOriginalStatement = tabular.DEFAULT_ORIGINAL_STATEMENT_OUTPUT
		}

		// Parameter: IG Script inclusion
		val, suc = extractUrlParameters(r, shared.PARAM_PRINT_IG_SCRIPT)
		if val != "" {
			// Read from parameter
			retStruct.PrintIgScript = val
		} else {
			// Use default parameter
			retStruct.PrintIgScript = tabular.DEFAULT_IG_SCRIPT_OUTPUT
		}

		// Parameter: Output type
		val, suc = extractUrlParameters(r, shared.PARAM_OUTPUT_TYPE)
		if val != "" {
			// Read from parameter
			retStruct.OutputType = val
		} else {
			// Use default parameter
			retStruct.OutputType = tabular.DEFAULT_OUTPUT_TYPES
		}

		// VISUAL PARAMETERS

		// Parameter: Property tree
		val, suc = extractUrlParameters(r, shared.PARAM_PROPERTY_TREE)
		check = evaluateBooleanUrlParameters(shared.PARAM_PROPERTY_TREE, val, suc)
		// Manually override if not set - effectively defines default setting
		if !suc {
			check = true
		}
		// Assign values
		if check {
			retStruct.PrintPropertyTree = shared.CHECKBOX_CHECKED
			printFlatProperties = false
		} else {
			retStruct.PrintPropertyTree = shared.CHECKBOX_UNCHECKED
			printFlatProperties = true
		}

		// Parameter: DoV
		val, suc = extractUrlParameters(r, shared.PARAM_DOV)
		check = evaluateBooleanUrlParameters(shared.PARAM_DOV, val, suc)
		// Manually override if not set - effectively defines default setting
		if !suc {
			check = false
		}
		// Assign values
		if check {
			retStruct.IncludeDoV = shared.CHECKBOX_CHECKED
			includeDoV = true
		} else {
			retStruct.IncludeDoV = shared.CHECKBOX_UNCHECKED
			includeDoV = false
		}

		// Parameter: Binary tree
		val, suc = extractUrlParameters(r, shared.PARAM_BINARY_TREE)
		check = evaluateBooleanUrlParameters(shared.PARAM_BINARY_TREE, val, suc)
		// Assign values
		if check {
			retStruct.PrintBinaryTree = shared.CHECKBOX_CHECKED
			printBinaryTree = true
		} else {
			retStruct.PrintBinaryTree = shared.CHECKBOX_UNCHECKED
			printBinaryTree = false
		}

		// Parameter: Activation condition on top
		val, suc = extractUrlParameters(r, shared.PARAM_ACTIVATION_CONDITION_ON_TOP)
		check = evaluateBooleanUrlParameters(shared.PARAM_ACTIVATION_CONDITION_ON_TOP, val, suc)
		// Assign values
		if check {
			retStruct.ActivationConditionsOnTop = shared.CHECKBOX_CHECKED
			printActivationConditionsOnTop = true
		} else {
			retStruct.ActivationConditionsOnTop = shared.CHECKBOX_UNCHECKED
			printActivationConditionsOnTop = false
		}

		// Parameter: Canvas width
		val, suc = extractUrlParameters(r, shared.PARAM_WIDTH)
		if suc {
			width, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Error when interpreting URL parameter '"+shared.PARAM_WIDTH+"':", err)
			} else {
				retStruct.Width = width
			}
		}

		// Parameter: Canvas height
		val, suc = extractUrlParameters(r, shared.PARAM_HEIGHT)
		if suc {
			height, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Error when interpreting URL parameter '"+shared.PARAM_HEIGHT+"':", err)
			} else {
				retStruct.Height = height
			}
		}

		// Parameter: Execution - invokes immediate execution - i.e., continuation of function to the end
		val, suc = extractUrlParameters(r, shared.PARAM_EXECUTE_PARSER)
		check = evaluateBooleanUrlParameters(shared.PARAM_EXECUTE_PARSER, val, suc)

		// BY DEFAULT (i.e., if immediate execution is not set), simply refine populated fields in UI form and return
		if !check {

			// All URL parameters processed, but this is only returned if *no immediate execution* is requested ...

			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error processing default template:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			Println("Provided p/repopulated form")
			// Just p/repopulate template, but do not go beyond
			return
		}
	}

	// Actual parsing occurs from here on ... (if either POST is sent, or if immediate execution is demanded via URL parameter)

	// Initialize request-specific logfile first
	if Logging {
		log.Println("Logging enabled")
		tID, filename := helper.GenerateUniqueIdAndFilename()
		// Assign transaction ID
		retStruct.TransactionId = tID
		// Check whether logging path has terminating slash
		if LoggingPath != "" {
			if LoggingPath[len(LoggingPath)-1:] != "/" {
				LoggingPath += "/"
			}
		}
		// Perform the file redirection
		var err error
		terminateOutput, err = helper.SaveOutputToFile(LoggingPath + filename)

		fmt.Println("TRANSACTION ID: " + retStruct.TransactionId)
		if err != nil {
			fmt.Println("Error when initializing logging: " + err.Error())
		}
	}

	// Logging and catching empty input for tabular and visual output
	if templateName != TEMPLATE_NAME_PARSER_PRODUCTION {
		fmt.Println("Input values:\n" +
			"RAW STATEMENT: " + retStruct.RawStmt + "\n" +
			"ANNOTATED STATEMENT: " + retStruct.CodedStmt + "\n")

		// Check for empty input statement first
		if retStruct.CodedStmt == "" {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error generating error response for empty input:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}

			// Final comment in log
			Println("Error: No input to parse.")
			// Ensure logging is terminated
			err2 := terminateOutput(ERROR_SUFFIX)
			if err2 != nil {
				log.Println("Error when finalizing log file: ", err2.Error())
			}
			return
		}
	} else {
		// Delegate to specific output handlers ...
		if templateName == TEMPLATE_NAME_PARSER_TABULAR {
			Println("Tabular output requested")
			handleTabularOutput(w, retStruct.RawStmt, retStruct.CodedStmt, retStruct.StmtId, retStruct, dynamicOutput, produceIGExtendedOutput, includeAnnotations, retStruct.OutputType, printHeaders, formValuePrintOriginalStatement, formValuePrintIgScript)
		} else if templateName == TEMPLATE_NAME_PARSER_VISUAL {
			Println("Visual output requested")
			handleVisualOutput(w, retStruct.CodedStmt, retStruct.StmtId, retStruct, printFlatProperties, printBinaryTree, printActivationConditionsOnTop, dynamicOutput, produceIGExtendedOutput, includeAnnotations, includeDoV)
		} else if templateName == TEMPLATE_NAME_PARSER_PRODUCTION {
			Println("Production output requested")
			handleProductionOutput(w, r, retStruct, produceIGExtendedOutput, includeAnnotations)
		} else {
			log.Fatal("Output variant " + templateName + " not found.")
		}
	}
}

/*
Extract URL parameters for further processing. It takes the parameter key (as specified in UrlParameters.go) as input
and returns the associated value, as well as a flag indicating success of extraction (i.e., whether a parameter had
been present in the first place).
*/
func extractUrlParameters(r *http.Request, parameterKey string) (string, bool) {
	keys, ok := r.URL.Query()[parameterKey]
	if ok && len(keys[0]) > 0 {

		// Assume single item
		key := keys[0]

		// Return entry as string and signal success (even if value is empty)
		return string(key), true
	}
	// Return empty string (and signal absence of URL parameter)
	return "", false
}

/*
Evaluates URL parameter value input for boolean variants:
- "true", "t", and "1" are interpreted as "on"
- "false", "f", and "0" are interpreted as "off"
- If success is false, the function returns false (i.e., no URL parameter of the given name found).
*/
func evaluateBooleanUrlParameters(parameter string, value string, success bool) bool {
	if success {
		switch value {
		case "t":
			return true
		case "true":
			return true
		case "1":
			return true
		case "f":
			return false
		case "false":
			return false
		case "0":
			return false
		default:
			log.Println("Invalid URL parameter value for parameter '" + parameter + "': " + value)
		}
	}
	return false
}
