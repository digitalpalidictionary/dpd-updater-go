# GoldenDict Configuration Research

This document outlines the configuration storage and structure for GoldenDict and GoldenDict-NG across different operating systems.

## 1. Configuration File Locations

GoldenDict stores its settings and dictionary paths in the following default locations:

- **Windows:** `%APPDATA%\GoldenDict` (typically `C:\Users\<User>\AppData\Roaming\GoldenDict`)
- **macOS:** `~/Library/Application Support/GoldenDict`
- **Linux:** `~/.config/goldendict/`

**Portable Mode:** If a folder named `portable` exists in the same directory as the GoldenDict executable, all configuration and dictionary indexes are stored there instead.

## 2. File Format
The primary configuration file is named **`config`** (no file extension). Although it lacks an extension, it is a standard **XML** file.

## 3. Data Structure
The configuration is wrapped in a `<config>` root element. Dictionary sources are managed within the `<paths>` section.

### Example XML Structure:
```xml
<config>
  <paths>
    <path recursive="true" enabled="true">C:/Path/To/Dictionaries</path>
    <path recursive="false" enabled="true">/home/user/dicts</path>
  </paths>
  <!-- Other settings -->
</config>
```

## 4. Dictionary Paths and Recursion

To programmatically determine which dictionary folders are active and whether they scan subdirectories:

1.  **Identify Folders:** Locate the `<path>` elements inside the `<paths>` block.
2.  **Enabled Status:** Check the `enabled="true"` attribute. Only enabled paths are currently being used by the application.
3.  **Recursion:** Check the `recursive` attribute:
    -   `recursive="true"`: GoldenDict performs a recursive scan of the folder.
    -   `recursive="false"`: GoldenDict only scans the top-level files in that folder.

These settings correspond to the UI found in **F3 (Dictionaries) -> Sources -> Files**.
