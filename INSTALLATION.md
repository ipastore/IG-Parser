# Installation Guide

## IG-Parser-Production
Follow these steps to install the IG-Parser-Production on your local machine:

### Step 1: Verify Your Oprtating System Version

To ensure you download the correct installers, first verify your operating system version.

**Instructions:**

**For Windows:**
1. Click on the **Start** button in the bottom bar of your screen.
2. Type **“System Information”** and select the application that appears.
3. In the **System Information** window, look for the following two lines:
   - **System Type:** Indicates whether your system is 32-bit or 64-bit.
   - **OS Name:** Shows the version of Windows you are using.

**For MacOS:**
1.	Click on the Apple menu in the top-left corner of your screen.
2.	Select About This Mac.
3.	In the window that appears, look for:
	•	macOS Version: Shows the version of macOS you are using (e.g., macOS Ventura 13.0).
	•	Processor Information: Indicates if your Mac is using an Intel or Apple Silicon (M1, M2) processor.


### Step 2: Install Git

Git is a version control tool that allows you to download (clone) the web application’s source code from GitHub.

**Instructions:**

1. Open your internet browser and go to [https://git-scm.com](https://git-scm.com).
2. Download the appropriate version of Git for your operating system.
3. Follow the installation instructions, accepting the license terms and leaving the default options selected. Click **Next** until you see the **Install** button.
4. Click **Install** and wait for the installation to complete.

### Step 3: Install Go (Golang)

Go is the programming language used to develop the application. You need to install Go to run the web application.

**Instructions:**

1. Open your internet browser and go to [https://go.dev/dl/](https://go.dev/dl/).
2. Select the version of Go that is compatible with your operating system.
3. Download the installer and double-click the file to start the installation.
4. Follow the installer instructions, accepting the license terms and leaving the default options selected.
5. Click **Install** and wait for the installation to complete.

### Step 4: Clone the GitHub Repository

Now, clone the GitHub repository where the web application is located.

**Instructions:**

1. Under Windows open **Git Bash** (an application installed with Git). You can find it in the Windows Start menu by searching for **"Git Bash"**. Under Linux (MacOS) open **Terminal**
2. Type the following command and press Enter:

   ```sh
   git clone https://github.com/ipastore/IG-Parser-Production.git
   ```

3. Wait for the cloning process to complete. This will download the application’s source code to your computer.

### Step 5: Run the Go Build Command

Compile the source code to create the executable file for the web application.

**Instructions:**

1.	Open Git Bash or Terminal again or continue in the same window.
2.	Navigate to the folder where the repository was cloned by typing:


  ```sh
    cd IG-Parser-Production
   ```

3. Compile the application by running the following command:

Under Windows in Git Bash:
```sh
    go build -o ig-parser.exe ./web 
```

Under Linux (MacOS) in Terminal:
```sh
    go build -o ig-parser ./web 
```

This will create an executable file named ig-parser in the repository folder (IG-Parser-Production).

### Conclusion

You have successfully installed the IG-Parser-Production on your local machine. Now you can proceed to the [Usage Guide](USAGE.md) to learn how to run and use the application.

## Further help

If you are a Windows user and need a more detailed step-by-step installation guide with illustrations and examples, please refer to this [Installation and Usage Guide](docs/InstallationAndUsageGuide.pdf) in PDF. If you are a Linux user this guide should also be useful to understand the process of Installation.

