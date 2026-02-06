# Specification: GoldenDict Config Auto-Detection

## Overview
Enhance the DPD Updater setup wizard to automatically detect existing GoldenDict dictionary paths. By reading the GoldenDict `config` file, the updater can suggest the most likely installation folder, reducing friction for the user during the initial setup or update process.

## Functional Requirements

### 1. GoldenDict Config Discovery
- Locate the GoldenDict configuration file based on the OS:
    - **Linux:** `~/.config/goldendict/config`
    - **Windows:** `%APPDATA%\GoldenDict\config`
    - **macOS:** `~/Library/Application Support/GoldenDict/config`
- Support "Portable Mode" by checking for a `config` file in the executable's directory (if applicable/detectable).

### 2. Path Extraction & Analysis
- Parse the XML `config` file to extract all `<path>` elements within `<paths>`.
- Only consider paths where `enabled="true"`.
- Determine the `recursive` status for each path.

### 3. Suggestion Logic
- **Single Folder (Recursive):** If exactly one enabled path is found and it is marked `recursive="true"`, suggest this folder.
- **Single Folder (Non-Recursive):** If exactly one enabled path is found and it is `recursive="false"`, do not suggest; fallback to manual selection.
- **Multiple Folders (Common Parent):** If multiple enabled paths share a common parent directory, suggest that common parent.
- **Multiple Folders (No Common Parent):** If no common parent is identifiable, do not suggest a specific path. Fallback to manual selection.

### 4. UI Interaction
- **Success:** If a suggestion is found, pre-fill the folder selection field in the Setup Wizard and inform the user that it was auto-detected from GoldenDict.
- **Failure/Missing Config:** If the config is missing or unreadable, show a brief warning/toast ("Could not auto-detect GoldenDict settings") and default to the manual folder picker.
- **Manual Override:** The user must always be able to use the manual file picker to change the auto-detected path.
- **Best Practice Guidance:** In cases where multiple folders are detected with no common parent (or no folders found), display a helper tip: "Recommendation: Organize all your dictionary files into a single master folder (e.g., `Documents/GoldenDict`) and select that folder here."

## Non-Functional Requirements
- **Robustness:** The updater must not crash if the XML is malformed or permission is denied to the config folder.
- **Privacy:** Only read the `<paths>` section; ignore all other user settings (history, etc.).

## Out of Scope
- **File Cleanup:** Removal of duplicate DPD folders or legacy versions is NOT included in this track.

## Acceptance Criteria
- [ ] Updater successfully finds the `config` file on the current OS.
- [ ] Correct paths are extracted from the XML.
- [ ] Suggestion logic behaves as specified for single/multiple folders.
- [ ] UI correctly pre-fills the path and allows manual override.
- [ ] Warning is shown when auto-detection fails.
- [ ] Best practice guidance is displayed when the folder structure is fragmented (multiple folders, no common parent).
