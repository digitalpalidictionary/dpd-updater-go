# DPD Updater Go

Standalone DPD Updater written in Go using the Fyne GUI toolkit.

## Project Structure
- `cmd/dpd-updater`: Entry point for the application.
- `internal/config`: Configuration management.
- `internal/github`: GitHub API integration.
- `internal/installer`: Installation and backup logic.
- `internal/system`: OS detection and process management.
- `internal/ui`: Fyne GUI components.

## Development
To run the application:
```bash
go run ./cmd/dpd-updater
```
