package main

import (
	"IG-Parser/core/config"
	"IG-Parser/web/converter"
	"IG-Parser/web/helper"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
This file is the main entry point for the IG Parser as a web application.
It relies on the IG Parser core package functionality.
*/

// Environment variables (port, logging activation)
const ENV_VAR_PORT = "IG_PARSER_PORT"
const ENV_VAR_LOGGING = "IG_PARSER_LOGGING"
const ENV_VAR_LOGGING_PATH = "IG_PARSER_LOGGING_PATH"

// Default values
const DEFAULT_LOGGING_PATH = "./logs"
const DEFAULT_PORT = "8080"

// Control whether stdout console output should be suppressed (only works if logging is deactivated)
const SUPPRESS_CONSOLE_OUTPUT = false

// API Paths (embed trailing slashes to retain all URL control here)
const TABULAR_PATH = "" // empty per default
const VISUAL_PATH = "visual/"
const HELP_PATH = "help/"
const PRODUCTION_PATH = "production/"

// Embed external files in compiled binary filesystem

//go:embed css/default.css css/favicon.ico
var cssFiles embed.FS

//go:embed libraries/d3.v7.min.js libraries/ace/ace.js
var libraryFiles embed.FS

//go:embed converter/templates/ig-parser-user-guide.html
var helpFiles embed.FS

/*
Main entry point for web version of IG Parser.
*/
func main() {

	// Initializes templating and determines correct relative path for templates and CSS
	converter.Init()

	// Register static resources

	// D3 & ACE libraries
	http.Handle("/libraries/", http.FileServer(http.FS(libraryFiles)))
	// CSS folder mapping (for CSS and favicon)
	http.Handle("/css/", http.FileServer(http.FS(cssFiles)))

	// Register handlers

	// Conventional tabular output handler (path per default empty)
	http.HandleFunc("/"+TABULAR_PATH, converter.ConverterHandlerTabular)
	// Visual tree output handler
	http.HandleFunc("/"+VISUAL_PATH, converter.ConverterHandlerVisual)
	// Help handler
	http.HandleFunc("/"+HELP_PATH, converter.HelpHandler)
	// Production handler
	http.HandleFunc("/"+PRODUCTION_PATH, converter.ConverterHandlerProduction)

	// Check for custom port
	port := os.Getenv(ENV_VAR_PORT)
	if port == "" {
		port = DEFAULT_PORT
	}

	// Check for logging specification (default activated)
	logEnv := os.Getenv(ENV_VAR_LOGGING)
	if logEnv == "" || strings.ToLower(logEnv) == "true" {
		converter.Logging = true
	} else {
		converter.Logging = false
	}

	// Check for logging path (if logging is enabled)
	if converter.Logging == true {
		logPath := os.Getenv(ENV_VAR_LOGGING_PATH)
		if logPath != "" {
			// Path from environment
			converter.LoggingPath = logPath
		} else {
			// Choose default path if not specified (but logging activated)
			converter.LoggingPath = DEFAULT_LOGGING_PATH
		}
	}

	// Suppress stdout (to be used with care) - only works if logging is deactivated
	if SUPPRESS_CONSOLE_OUTPUT && converter.Logging == false {
		os.Stdout = nil
	}

	// Compose port suffix
	portSuffix := ":" + port

	// Print configuration in console
	log.Println("Launching IG Parser ...")
	log.Println(" - Version: " + config.IG_PARSER_VERSION)
	log.Println(" - Website: https://newinstitutionalgrammar.org/ig-parser")
	log.Println(" - Logging enabled: " + fmt.Sprint(converter.Logging))
	log.Println(" - Logging path: " + fmt.Sprint(converter.LoggingPath))
	log.Printf("Navigate to the URL http://localhost%s/"+TABULAR_PATH+" in your browser to open the tabular output version of IG Parser.\n", portSuffix)
	log.Printf("Navigate to the URL http://localhost%s/"+VISUAL_PATH+" in your browser to open the visual output version of IG Parser.\n", portSuffix)
	log.Printf("Navigate to the URL http://localhost%s/"+PRODUCTION_PATH+" in your browser to open the production output version of IG Parser.\n", portSuffix)

	// Attempt launch of URL in browser
	err0 := helper.OpenBrowser("http://localhost" + portSuffix + "/" + PRODUCTION_PATH)
	if err0 != nil {
		log.Println("Browser launch failed (Error: " + err0.Error() + "). Please launch browser manually using URLs above.")
	}
	// Launch web server
	err := http.ListenAndServe(portSuffix, nil)
	if err != nil {
		log.Fatal("Web service stopped. Error:", err)
	}

}
