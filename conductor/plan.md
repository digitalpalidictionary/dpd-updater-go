# Implementation Plan: DPD Updater Go Conversion

## Phase 1: Environment & Scaffolding
- [x] Task: Initialize standalone project structure
    - [x] Initialize `git` and `go mod` in `resources/dpd-updater-go`
    - [x] Create directory structure: `cmd/`, `internal/config`, `internal/github`, `internal/installer`, `internal/system`, `internal/ui`
- [x] Task: Setup Fyne and cross-platform build tooling
    - [x] Install Fyne dependencies
    - [x] Create a basic "Hello World" window to verify build
- [x] Task: Conductor - User Manual Verification 'Phase 1' (Protocol in workflow.md)

## Phase 2: Logic Layer Implementation
- [x] Task: Implement Config and System managers
    - [x] JSON configuration persistence (using `internal/config`)
    - [x] OS detection and path resolution logic (using `internal/system`)
    - [x] GoldenDict process control (list/kill/restart)
    - [x] Directory scanning and version detection logic
- [x] Task: Implement GitHub and Installer logic
    - [x] GitHub API client for release fetching (using `internal/github`)
    - [x] HTTP Downloader with progress reporting (io.Reader wrapping)
    - [x] Zip extractor with backup mechanism (using `internal/installer`)
- [x] Task: Conductor - User Manual Verification 'Phase 2' (Protocol in workflow.md)

## Phase 3: GUI Development (Fyne)
- [x] Task: Implement UI Framework and State
    - [x] Define UI state and messaging channels (Go channels for updates)
    - [x] Basic application theme and layout
- [x] Task: Implement UI Views
    - [x] First-run Setup Wizard view (Path selection/validation)
    - [x] Main Dashboard view (Status/Action/Log)
    - [x] Settings dialog
- [x] Task: Integrate UI with Logic
    - [x] Bind logic events to UI updates (Reactive UI)
    - [x] Handle asynchronous operations with loading states
- [~] Task: Conductor - User Manual Verification 'Phase 3' (Protocol in workflow.md)

## Phase 4: CI/CD & Final Delivery
- [ ] Task: Logo and Asset Integration
    - [ ] Bundle application icon and assets using `fyne bundle`
- [ ] Task: Create GitHub Build Workflow
    - [ ] Port `build-and-release.yml` logic to Go
    - [ ] Ensure Windows, Linux, and macOS targets are supported using `fyne-cross`
- [ ] Task: Documentation and Quality Control
    - [ ] Write `README.md` for `dpd-updater-go`
    - [ ] Final verification of feature parity
    - [ ] Run `go fmt` and `go vet`
- [ ] Task: Conductor - User Manual Verification 'Phase 4' (Protocol in workflow.md)
