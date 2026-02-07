# Specification - Duplicate DPD Detection and Cleanup

## Overview
Users sometimes have multiple copies of the DPD dictionary in their GoldenDict folder (e.g., after manual moves or multiple installation attempts). This causes GoldenDict to show duplicate entries, which impairs functionality. This feature will detect these duplicates, identify the older ones based on the `date` field in `dpd.ifo`, and allow the user to delete them directly from the UI.

## Functional Requirements
1.  **Duplicate Detection:**
    -   Scan the entire GoldenDict dictionary directory recursively for `dpd.ifo` files.
    -   Parse the `date` field from each `dpd.ifo` file (ISO 8601 format).
2.  **UI Feedback:**
    -   Modify the "Version Information" display in the UI.
    -   If multiple copies are found, list each one with its installation date and path.
    -   Use Red or Orange text for the "Installed" status when duplicates exist.
3.  **Cleanup Action:**
    -   Provide a button (e.g., "Delete Old Versions") when duplicates are detected.
    -   The button should remove the directories containing all but the newest version (based on date).
4.  **Execution Timing:**
    -   Perform the scan on application startup.
    -   Perform the scan during the update check process.

## Acceptance Criteria
-   [ ] The application detects all `dpd.ifo` files within the GoldenDict folder.
-   [ ] The UI displays a warning when more than one copy is found.
-   [ ] The "oldest" copies are accurately identified by comparing the `date` field.
-   [ ] The "Delete" action successfully removes the older folders and updates the UI state.
