# Implementation Plan - Duplicate DPD Detection and Cleanup

## Phase 1: Scanner Logic & Parsing
- [x] Task: Create `internal/system/dpd_info.go` to handle `dpd.ifo` parsing
    - [x] Define `DPDInfo` struct with `Path` and `Date` (time.Time).
    - [x] Implement `ParseIFO(path string) (*DPDInfo, error)` to parse the `.ifo` file format.
- [x] Task: Implement recursive scanner in `internal/system/scanner.go`
    - [x] Add `FindAllDPDInstances(root string) ([]DPDInfo, error)`.
    - [x] Ensure it handles empty or malformed `.ifo` files gracefully.
- [x] Task: Write tests for scanner and parser
    - [x] Mock directory structure with multiple `dpd.ifo` files (varying dates).
    - [x] Verify correct parsing of the `date` field.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Scanner Logic & Parsing' (Protocol in workflow.md)

## Phase 2: UI Integration & "Nag" Display
- [x] Task: Update `internal/ui/state.go` to track DPD instances
    - [x] Add `DPDInstances []system.DPDInfo` to the application state.
- [x] Task: Modify `internal/ui/main_window.go` for Version Info display
    - [x] Update the label/text area to list all found instances.
    - [x] Apply Red/Orange styling if `len(DPDInstances) > 1`.
- [x] Task: Add "Delete Old Versions" button logic
    - [x] Filter `DPDInstances` to find the newest.
    - [x] Implement `DeleteFolders([]string paths)` in `internal/installer/installer.go` or similar.
    - [x] Bind button visibility to the presence of duplicates.
- [x] Task: Conductor - User Manual Verification 'Phase 2: UI Integration & "Nag" Display' (Protocol in workflow.md)

## Phase 3: Lifecycle Integration & Final Cleanup
- [x] Task: Trigger scan on Startup
    - [x] Call `FindAllDPDInstances` during app initialization.
- [x] Task: Trigger scan on Update Check
    - [x] Ensure the scan runs after any update attempt or manual check.
- [x] Task: Documentation update
    - [x] Update `README.md` and any relevant docs in `docs/` regarding the cleanup feature.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Lifecycle Integration & Final Cleanup' (Protocol in workflow.md)
