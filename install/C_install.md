### C Installation Guide

---

#### Installing GCC on Your Computer

GCC (GNU Compiler Collection) is a popular open-source compiler that supports various programming languages, including C. Here’s how to install it on different operating systems.

<details>
  <summary>Windows</summary>
  
  1. **Download MinGW**:
     - Go to the [MinGW website](http://www.mingw.org/) and download the MinGW installation manager.

2. **Install MinGW**:

   - Run the installer and follow the setup instructions. Make sure to select the "C Compiler" component during installation.

3. **Add to PATH**:

   - After installation, you need to add the MinGW `bin` directory to your system's PATH.
   - Right-click on **This PC** > **Properties** > **Advanced system settings** > **Environment Variables**.
   - Find the `Path` variable in the System variables section, click **Edit**, and add the path to your MinGW `bin` directory (e.g., `C:\MinGW\bin`).

4. **Verify Installation**:
   - Open Command Prompt and run:
     ```bash
     gcc --version
     ```

</details>

<details>
  <summary>macOS</summary>

1. **Install Homebrew** (if you haven’t already):

   - Open Terminal and run:
     ```bash
     /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
     ```

2. **Install GCC**:

   - In the Terminal, run:
     ```bash
     brew install gcc
     ```

3. **Verify Installation**:
   - After installation, check the version:
     ```bash
     gcc --version
     ```

</details>

<details>
  <summary>Linux</summary>

1. **Using APT (Debian/Ubuntu)**:

   - Open Terminal and run:
     ```bash
     sudo apt update
     sudo apt install build-essential
     ```

2. **Using DNF (Fedora)**:

   - In the Terminal, run:
     ```bash
     sudo dnf groupinstall 'Development Tools'
     ```

3. **Verify Installation**:
   - Check the GCC version:
     ```bash
     gcc --version
     ```

</details>
<br/>

_For any issues please visit [Microsoft's Installation Guide for C/C++](https://learn.microsoft.com/en-us/cpp/build/vscpp-step-0-installation?view=msvc-170)_

---
