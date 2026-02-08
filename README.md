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

Download the latest pre-compiled binaries from the [GitHub Releases](https://github.com/digitalpalidictionary/dpd-updater-go/releases) page.

## Development

### Prerequisites

1. Install Go 1.22+
2. Install Fyne dependencies:
   - **Linux**: `sudo apt-get install libgl1-mesa-dev xorg-dev`
   - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
   - **Windows**: No additional dependencies

### Running from Source

```bash
go run .
```

### Building

#### Option 1: Using `fyne package` (Recommended)

The Fyne CLI handles icons, manifests, and app bundles automatically:

```bash
# Install Fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest

# Build for current platform
fyne package -os windows -icon assets/icon.png -appID net.dpdict.dpd-updater -name dpd-updater
fyne package -os darwin  -icon assets/icon.png -appID net.dpdict.dpd-updater -name dpd-updater
fyne package -os linux   -icon assets/icon.png -appID net.dpdict.dpd-updater -name dpd-updater
```

#### Option 2: Using GitHub Actions

The project includes a workflow that builds for all platforms automatically:
1. Go to Actions → "Build & Release"
2. Click "Run workflow"

This is the recommended approach for releases since Fyne requires platform-specific builds.

## Project Structure

```
dpd-updater-go/
├── main.go              # Application entry point
├── internal/            # Internal packages
│   ├── config/          # Configuration handling
│   ├── github/          # GitHub API client
│   ├── installer/       # Installation logic
│   ├── system/          # System detection (GoldenDict, DPD info)
│   └── ui/              # Fyne UI components
├── assets/              # Icons and images
└── .github/workflows/   # CI/CD workflows
```

## License

MIT License
