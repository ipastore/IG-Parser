# IG-Parser-Production

[![License](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

<div align="center">
  <img src="" alt="RESILIENT RULES Logo" height="100">
  <img src="assets/00logo-ERC.png" alt="ERC and EU Funding Logo" height="130">
</div>


New Production feature for IG-Parser

### Contact: 
- <u>Developer:</u> Ignacio Pastore Benaim (ignaciopastorebenaim@gmail.com)
- <u>Principal Investigator:</u> Irene Pérez Ibarra (perezibarra@unizar.es)

See [Contributors](contributors.md) 

### Introduction

This repository is a fork of the [IG-Parser](https://github.com/chrfrantz/IG-Parser) created by [chrfrantz](https://github.com/chrfrantz). The purpose of this fork is to enhance the original IG-Parser with additional features and improvements, specifically designed for the **RESILIENT RULES** project (ERC-2021-CoG).

**Note:** This version of the IG-Parser-Production is currently designed for **local deployment only**. Future updates will include development for server deployment.

Additionally, both the **tabular** and **visual outputs** are fully functional and it is **recommended** to use them before populating the Excel file as input. This helps in understanding the structure and required formatting of the input data as explained in the original repository.

### New Features in This Fork

- **Excel Production Module**: Developed specifically for the RESILIENT RULES project, allowing for the production of parsed outputs from input Excel files with encoded statements.
- **User-Friendly Interface**: Accessible via a web browser at `http://localhost:8080/production/` for ease of use.
- **New Functionalities Planned**:
  - Automatic Statement ID generation based on selected Excel columns.
  - Option to display full names or symbols in headers.
  - Expanded coding options for multiple sheets within an Excel file.
  - **Server Deployment**: Future versions will include development for server deployment.


### Installation and local deployment

The purpose of building a local executable is to run IG Parser on a local machine (primarily for personal use on your own machine).

For detailed installation instructions fow Windows please refer to the [Installation Guide](INSTALLATION.md). This guide is tailored for users who are new to programming or need step-by-step instructions.

If you are familiar with programming, you can follow the quick setup instructions below:

#### Quick Setup for Experienced Users

* Prerequisites:
  * Install [Go (Programming Language)](https://go.dev/dl/)
  * Clone this repository into a dedicated folder on your local machine
  * Navigate to the repository folder
  * Compile IG Parser in the corresponding console
    * Under Windows, execute `go build -o ig-parser.exe ./web`
      * This creates the executable `ig-parser.exe` in the repository folder
    * Under Linux, execute `go build -o ig-parser ./web`
      * This creates the executable `ig-parser` in the repository folder
  * Run the created executable
    * Under Windows, run `ig-parser` (or `ig-parser.exe`) either via command line or by doubleclicking
    * Under Linux (or Windows PowerShell), run `./ig-parser`
  * Once started, it should automatically open your browser and navigate to http://localhost:8080/visual. Alternatively, use your browser to manually navigate to one of the URLs listed in the console output. By default, this is the URL http://localhost:8080 (and http://localhost:8080/visual respectively)
  * Press `Ctrl` + `C` in the console window to terminate the execution (or simply close the console window)

### Usage

For detailed instructions on how to run and use the application, please refer to the [Usage Guide](USAGE.md).


### License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE.txt) file for details.

### Acknowledgements

- [chrfrantz](https://github.com/chrfrantz) for the original IG-Parser project.




