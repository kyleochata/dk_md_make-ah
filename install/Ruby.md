# Ruby Installation Guide

## Installing Ruby on Your Computer

Ruby is a dynamic, open-source programming language with a focus on simplicity and productivity. Here’s how to install Ruby on different operating systems.

<details>
  <summary>Windows</summary>
  
  1. **Download RubyInstaller**:
     - Go to the [RubyInstaller website](https://rubyinstaller.org/) and download the latest version.

2. **Install Ruby**:

   - Run the installer and follow the setup instructions. Make sure to check the option to add Ruby to your PATH.

3. **Verify Installation**:
   - Open Command Prompt and run:
     ```bash
     ruby -v
     ```

</details>

<details>
  <summary>macOS</summary>

1. **Install Homebrew** (if you haven’t already):

   - Open Terminal and run:
     ```bash
     /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
     ```

2. **Install Ruby**:

   - In the Terminal, run:
     ```bash
     brew install ruby
     ```

3. **Verify Installation**:
   - After installation, check the Ruby version:
     ```bash
     ruby -v
     ```

</details>

<details>
  <summary>Linux</summary>

1. **Using APT (Debian/Ubuntu)**:

   - Open Terminal and run:
     ```bash
     sudo apt update
     sudo apt install ruby-full
     ```

2. **Using DNF (Fedora)**:

   - In the Terminal, run:
     ```bash
     sudo dnf install ruby
     ```

3. **Verify Installation**:
   - Check the Ruby version:
     ```bash
     ruby -v
     ```

</details>
<br>

_For any issues please visit the [Ruby Installation Documentation](https://www.ruby-lang.org/en/documentation/installation/)_

---
