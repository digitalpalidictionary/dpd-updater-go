# Specification: DPD Updater Go Conversion

## Overview
Convert the existing Python/Flet-based DPD Updater (`resources/dpd-updater`) into a high-performance, standalone Go application using the Fyne GUI framework. The new project will reside in its own repository and will be structured as a self-contained repository that can be integrated as a submodule.

## Goals
- Achieve 1:1 feature parity with the current Python implementation.
- Improve startup speed and distribution by using a compiled language (Go).
- Provide a modern, responsive GUI using Fyne.
- Automate cross-platform builds via GitHub Actions.

## Functional Requirements
- **OS Detection & Path Resolution:** 
    - Automatically identify GoldenDict configuration paths on Windows, Linux, and macOS.
    - Support both GoldenDict and GoldenDict-ng.
- **Scanning Logic:**
    - On startup, scan the configured GoldenDict path for existing `dpd*` files.
    - Extract version metadata from existing files to determine the currently installed version.
- **GitHub API Integration:** Fetch latest release data and download assets with progress tracking.
- **Transactional Installation:**
    - Create timestamped backups of existing `dpd*` files.
    - Extract new data to temporary locations before final placement.
    - Clean up temporary files post-installation.
- **Process Management:** Ability to stop and start the GoldenDict application to facilitate data updates.
- **User Interface (Fyne):**
    - **Setup Wizard:** First-time configuration for GoldenDict paths.
    - **Main Dashboard:** 
        - Version status (Installed vs. Latest).
        - Update button (enabled only when update available).
        - "Check for updates" button.
    - **Operation Log:** Scrollable real-time log of background tasks.
    - **Progress Indicators:** Visual feedback during downloads and extraction.
    - **Settings:** Manual path configuration and preference toggles.

## Technical Considerations
- **Concurrency:** Use Go routines for background downloads and file operations to keep the GUI responsive.
- **Error Handling:** Robust error reporting in the UI for network failures or permission issues.
- **Submodule Management:** The project resides as a Git submodule.

## Non-Functional Requirements
- **Standalone:** The `resources/dpd-updater2` directory must be a fully functional Git repository root.
- **Binary Size:** Optimize for a reasonable binary size for distribution.
- **Code Quality:** Adhere to Go idioms and project standards.

## Out of Scope
- Adding new dictionary features not present in the original Python updater.
- Direct database manipulation (the updater handles file distribution only).

## Acceptance Criteria
- [ ] Successful cross-platform build (Windows, Linux, macOS).
- [ ] User can update DPD to the latest version via the Go GUI.
- [ ] Backup created successfully before update.
- [ ] GoldenDict restarts automatically after update (if possible on the platform).
