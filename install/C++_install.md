# C++ Installation Guide

## Installing a C++ Compiler on Your Computer

To write and compile C++ programs, you need a C++ compiler. GCC (GNU Compiler Collection) and Clang are popular compilers that support C++. Here’s how to install them on different operating systems.

<details>
  <summary>Windows</summary>
  
  1. **Download MinGW-w64**:
     - Go to the [MinGW-w64 website](http://mingw-w64.org/doku.php/download) and download the installer.

2. **Install MinGW-w64**:

   - Run the installer and follow the setup instructions. Make sure to select the C++ compiler during installation.

3. **Add to PATH**:

   - After installation, you need to add the MinGW-w64 `bin` directory to your system's PATH.
   - Right-click on **This PC** > **Properties** > **Advanced system settings** > **Environment Variables**.
   - Find the `Path` variable in the System variables section, click **Edit**, and add the path to your MinGW-w64 `bin` directory (e.g., `C:\mingw-w64\bin`).

4. **Verify Installation**:
   - Open Command Prompt and run:
     ```bash
     g++ --version
     ```

</details>

<details>
  <summary>macOS</summary>

1. **Install Homebrew** (if you haven’t already):

   - Open Terminal and run:
     ```bash
     /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
     ```

2. **Install GCC or Clang**:

   - In the Terminal, run:
     ```bash
     brew install gcc
     ```
   - Or to install Clang (which is included with Xcode Command Line Tools):
     ```bash
     xcode-select --install
     ```

3. **Verify Installation**:
   - After installation, check the version:
     ```bash
     g++ --version
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
     sudo dnf install gcc-c++
     ```

3. **Verify Installation**:
   - Check the C++ compiler version:
     ```bash
     g++ --version
     ```

</details>
<br/>

_For any issues please visit [Microsoft's Installation Guide for C/C++](https://learn.microsoft.com/en-us/cpp/build/vscpp-step-0-installation?view=msvc-170)_

---
