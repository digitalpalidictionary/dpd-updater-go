# DPD Updater Go

Standalone DPD Updater written in Go using the Fyne GUI toolkit. This application allows users to easily update their Digital Pāḷi Dictionary (DPD) installation for GoldenDict.

## Features
- **Cross-Platform:** Supports Windows, macOS, and Linux.
- **Fast:** Compiled Go binary with a modern Fyne GUI.
- **Automatic Path Detection:** Attempts to find GoldenDict configuration paths.
- **GitHub Integration:** Fetches the latest releases directly from GitHub.
- **Transactional Updates:** Includes backup and temporary extraction for safety.
- **Duplicate Detection:** Automatically identifies and helps clean up multiple DPD copies in your GoldenDict folder to ensure the best performance.

## Installation
You can download the latest pre-compiled binaries from the GitHub Releases page.

## Development
To build and run from source:

1. Install Go 1.22+
2. Install Fyne dependencies (on Linux: `sudo apt-get install libgl1-mesa-dev xorg-dev`)
3. Run the application:
   ```bash
   go run ./cmd/dpd-updater
   ```

To build an optimized production binary:
```bash
go build -ldflags="-s -w" -o dpd-updater ./cmd/dpd-updater
```

### Building for Multiple Platforms

Fyne uses platform-specific graphics libraries (OpenGL), so **cross-compilation is not supported**. 

To build for Windows, use one of these methods:

**Option 1: Use GitHub Actions (Recommended)**
The project includes a workflow that builds for all platforms automatically. Go to Actions → "Build & Release" and run the workflow.

**Option 2: Build on Windows directly**
On a Windows machine with Go installed:
```bash
go build -ldflags="-s -w" -o dpd-updater.exe ./cmd/dpd-updater
```

## License
MIT License
