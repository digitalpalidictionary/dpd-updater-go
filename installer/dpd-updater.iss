; DPD Updater Installer Script
; Inno Setup Script for DPD Updater

#define MyAppName "DPD Updater"
#define MyAppVersion "{VERSION_PLACEHOLDER}"
#define MyAppPublisher "Digital Pali Dictionary"
#define MyAppURL "https://github.com/digitalpalidictionary/dpd-updater-go"
#define MyAppExeName "dpd-updater.exe"
#define MyAppId "{{9F8E7D6C-5B4A-3210-9876-543210FEDCBA}}"

[Setup]
; Basic application info
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes

; Installation scope - let user choose
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

; Output configuration
OutputDir=.\Output
OutputBaseFilename=dpd-updater-windows-{#MyAppVersion}
Compression=lzma2
SolidCompression=yes

; Architecture support
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

; Version info for the installer itself
VersionInfoVersion={#MyAppVersion}
VersionInfoProductVersion={#MyAppVersion}
VersionInfoProductName={#MyAppName}
VersionInfoCompany={#MyAppPublisher}
VersionInfoDescription="DPD Updater for GoldenDict"
VersionInfoTextVersion={#MyAppVersion}

; Icon and branding
SetupIconFile=..\assets\icon.ico
UninstallDisplayIcon={app}\{#MyAppExeName}
UninstallDisplayName={#MyAppName}
WizardStyle=modern

; Other settings
CreateAppDir=yes
CreateUninstallRegKey=yes
Uninstallable=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Main executable
Source: "..\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
; Start Menu shortcuts
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{app}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"

; Desktop shortcut (optional based on task)
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{app}"; Tasks: desktopicon

[Run]
; Optional: Show message after install
Filename: "{app}\{#MyAppExeName}"; Description: "Launch DPD Updater"; Flags: nowait postinstall skipifsilent

[Registry]
; Store installation path for auto-detect in app
Root: HKCU; Subkey: "Software\Digital Pali Dictionary\DPD Updater"; ValueType: string; ValueName: "InstallDir"; ValueData: "{app}"; Flags: uninsdeletekey

[UninstallDelete]
; Clean up registry on uninstall
Type: dirifempty; Name: "{app}"

[Messages]
; Custom welcome message
WelcomeLabel1=Welcome to the DPD Updater Setup Wizard
WelcomeLabel2=This will install DPD Updater on your computer.%n%nDPD Updater helps you keep your Digital Pali Dictionary (DPD) up to date for GoldenDict.%n%nIt is recommended that you close GoldenDict before continuing.
