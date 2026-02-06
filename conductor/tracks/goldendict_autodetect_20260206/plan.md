# Implementation Plan: GoldenDict Config Auto-Detection

## Phase 1: Research & Core Logic Implementation

- [x] Task: Research Go XML parsing best practices for potentially large or malformed files.
- [x] Task: Research OS-specific environment variable access in Go (e.g., `APPDATA`, `HOME`).
- [x] Task: Implement `system.GetGoldenDictConfigPath()` to return the platform-specific path to the config file.
    - [x] Create `internal/system/goldendict_test.go`
    - [x] Write failing tests for path resolution on different OS (using environment variable overrides).
    - [x] Implement path resolution in `internal/system/goldendict.go`.
- [x] Task: Implement `system.ParseGoldenDictPaths(configPath string)` to extract enabled paths and recursion status.
    - [x] Write tests with sample XML content (single, multiple, recursive, disabled).
    - [x] Implement XML parsing in `internal/system/goldendict.go`.
- [x] Task: Implement `system.AnalyzeGoldenDictPaths(paths []GDPath)` for the suggestion logic.
    - [x] Write tests for:
        - Single recursive path -> suggest.
        - Single non-recursive path -> no suggestion.
        - Multiple paths with common parent -> suggest parent.
        - Multiple paths with no common parent -> no suggestion.
    - [x] Implement logic in `internal/system/goldendict.go`.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Research & Core Logic Implementation' (Protocol in workflow.md)

## Phase 2: UI Integration & User Guidance

- [~] Task: Update `ui.SetupWizard` to invoke the auto-detection logic during initialization.
- [~] Task: Modify the folder selection screen to pre-fill the detected path and show an "Auto-detected" status message.
- [~] Task: Implement the "Best Practice Guidance" tip for fragmented folder structures.
- [~] Task: Ensure the manual file picker remains functional and overrides any auto-detected path.
- [~] Task: Implement error handling (toast/warning) for cases where auto-detection fails.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: UI Integration & User Guidance' (Protocol in workflow.md)

## Phase 3: Finalization & Documentation

- [ ] Task: Run project-wide linting and formatting (`uv run ruff check --fix`, `uv run ruff format`).
- [ ] Task: Update `docs/goldendict_config.md` with technical details of the implementation.
- [ ] Task: Update `README.md` if any new configuration flags or behaviors were added.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Finalization & Documentation' (Protocol in workflow.md)
