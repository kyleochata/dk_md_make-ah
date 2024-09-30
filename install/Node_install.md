### **Node.js** Installation Guide

---

Follow these steps to install Node.js on your computer.

#### 1. Download and Install Node.js

##### For Windows, macOS, and Linux:

1. **Go to the official Node.js website:**

   Visit [https://nodejs.org/](https://nodejs.org/) to download the latest version of Node.js.

2. **Download the Node.js Installer:**

   - For **Windows**: Download the .msi installer.
   - For **macOS**: Download the .pkg installer.
   - For **Linux**: Download the .tar.xz archive or use your package manager.

3. **Run the Installer:**

   - **Windows**: Double-click the .msi file and follow the prompts.
   - **macOS**: Double-click the .pkg file and follow the prompts.
   - **Linux**: Extract the .tar.xz archive and install manually:

   ```bash
   tar -C /usr/local --strip-components 1 -xJf node-vX.X.X-linux-x64.tar.xz
   ```

4. **Add Node.js to your PATH (Linux/macOS users only):**

   If you installed Node.js manually, add the following to your .bashrc, .zshrc, or equivalent shell config file:

   ```bash
   export PATH=$PATH:/usr/local/bin
   ```

   Reload the shell config:

   ```bash
   source ~/.bashrc  # or source ~/.zshrc
   ```

#### 2. Verify Installation

Check your Node.js version:

```bash
node -v
```

Check your npm version (npm comes with Node.js):

```bash
npm -v
```

_For any issues please visit the official [Node.js Installation Guide](https://nodejs.org/en/learn/getting-started/how-to-install-nodejs)_

---
