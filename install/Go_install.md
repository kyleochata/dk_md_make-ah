### **Go** Installation Guide

---

Follow these steps to install Go (Golang) on your computer.

#### 1. Download and Install Go

##### For Windows, macOS, and Linux:

1. **Go to the official Go website:**

   Visit [https://go.dev/dl/](https://go.dev/dl/) to download the latest version of Go.

2. **Download the Go Installer:**

   - For **Windows**: Download the .msi installer.
   - For **macOS**: Download the .pkg installer.
   - For **Linux**: Remove any previous Go installation by deleting the `/usr/local/go` folder (if it exists), then :

   ```bash
   $ rm -rf /usr/local/go
   ```

3. **Run the Installer:**

   - **Windows**: Double-click the .msi file and follow the prompts.
     - By default, installer will instal Go to `Program Files` or `Program Files (x86)`.
   - **macOS**: Double-click the .pkg file and follow the prompts.
     - The package installs the Go distribution to `/usr/local/go`. The package should put the `/usr/local/go/bin` directory in your PATH environment variable. You may need to restart any open Terminal sessions for the cahnge to take effect.
   - **Linux**: extract the archive you just downloaded into `/usr/local`, creating a fresh Go tree in `/usr/local/go`:

   ```bash
   tar -C /usr/local -xzf go1.X.X.linux-amd64.tar.gz
   ```

4. **Add Go to your PATH (Linux/macOS users only)**

- Add the following to your .bashrc, .zshrc, or equivalent shell config file:
  ```bash
  export PATH=$PATH:/usr/local/go/bin
  ```
- Reload shell config:
  ```bash
  source ~/.bashrc # or source ~/.zshrc
  ```

5. **Verify Installation**

- Check your Go version:
  ```bash
  go version
  ```

_For any issues, please visit the official [Golang Installation Guide](https://go.dev/doc/install)._

---
