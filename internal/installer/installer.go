package installer

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
)

type ProgressCallback func(message string, percentage int)

type Installer struct {
	Config           *config.Config
	ProgressCallback ProgressCallback
	Client           *http.Client
}

func NewInstaller(cfg *config.Config, cb ProgressCallback) *Installer {
	return &Installer{
		Config:           cfg,
		ProgressCallback: cb,
		Client:           &http.Client{Timeout: 300 * time.Second},
	}
}

func (i *Installer) reportProgress(message string, percentage int) {
	if i.ProgressCallback != nil {
		i.ProgressCallback(message, percentage)
	}
}

func (i *Installer) DownloadRelease(ctx context.Context, url, destDir string) (string, error) {
	i.reportProgress("Starting download...", 0)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := i.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	filename := filepath.Base(url)
	if filename == "" || filename == "." {
		filename = "dpd-update.zip"
	}
	destPath := filepath.Join(destDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	totalSize := resp.ContentLength
	var downloaded int64

	buffer := make([]byte, 8192)
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				out.Write(buffer[:n])
				downloaded += int64(n)
				if totalSize > 0 {
					percentage := int((float64(downloaded) / float64(totalSize)) * 50)
					i.reportProgress(fmt.Sprintf("Downloading... %.1f / %.1f MB", float64(downloaded)/1024/1024, float64(totalSize)/1024/1024), percentage)
				}
			}
			if err == io.EOF {
				goto Done
			}
			if err != nil {
				return "", err
			}
		}
	}

Done:
	i.reportProgress("Download complete", 50)
	return destPath, nil
}

func (i *Installer) BackupExisting(ctx context.Context, gdPath string) (string, error) {
	if !i.Config.BackupBeforeUpdate {
		return "", nil
	}

	i.reportProgress("Creating backup...", 51)
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(gdPath, "backup_"+timestamp)

	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(gdPath)
	if err != nil {
		return "", err
	}

	backedUp := false
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			os.RemoveAll(backupDir)
			return "", ctx.Err()
		default:
			name := entry.Name()
			if stringsHasPrefix(name, "dpd") && name != filepath.Base(backupDir) {
				src := filepath.Join(gdPath, name)
				dst := filepath.Join(backupDir, name)
				if entry.IsDir() {
					if err := copyDir(src, dst); err == nil {
						backedUp = true
					}
				} else {
					if err := copyFile(src, dst); err == nil {
						backedUp = true
					}
				}
			}
		}
	}

	if backedUp {
		i.reportProgress(fmt.Sprintf("Backup created: %s", filepath.Base(backupDir)), 55)
		return backupDir, nil
	}

	os.RemoveAll(backupDir)
	return "", nil
}

func (i *Installer) InstallUpdate(ctx context.Context, zipPath, gdPath string) error {
	i.reportProgress("Extracting files...", 60)
	extractDir := filepath.Join(gdPath, "_dpd_update_temp")
	os.RemoveAll(extractDir)

	if err := unzip(ctx, zipPath, extractDir); err != nil {
		return err
	}

	i.reportProgress("Installing files...", 80)
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			src := filepath.Join(extractDir, entry.Name())
			dst := filepath.Join(gdPath, entry.Name())

			os.RemoveAll(dst)
			if err := os.Rename(src, dst); err != nil {
				if entry.IsDir() {
					copyDir(src, dst)
				} else {
					copyFile(src, dst)
				}
			}
		}
	}

	os.RemoveAll(extractDir)
	os.Remove(zipPath)
	i.reportProgress("Installation complete!", 100)
	return nil
}

func (i *Installer) RestoreFromBackup(backupDir, gdPath string) error {
	if backupDir == "" {
		return fmt.Errorf("no backup directory specified")
	}

	// Check if backup exists
	info, err := os.Stat(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("backup directory does not exist")
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("backup path is not a directory")
	}

	// Read backup contents
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("backup directory is empty")
	}

	// Restore each item from backup
	for _, entry := range entries {
		src := filepath.Join(backupDir, entry.Name())
		dst := filepath.Join(gdPath, entry.Name())

		// Remove any existing broken installation
		os.RemoveAll(dst)

		if entry.IsDir() {
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("failed to restore directory %s: %w", entry.Name(), err)
			}
		} else {
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("failed to restore file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

func DeleteFolders(paths []string) error {
	for _, p := range paths {
		// We delete the parent directory of dpd.ifo because dpd.ifo is inside the dictionary folder
		dir := filepath.Dir(p)
		// Safety check: ensure we are not deleting root or something unexpected
		// Ideally we should check if it's inside the GoldenDict folder but for now basic check
		if dir == "." || dir == "/" || len(dir) < 4 {
			return fmt.Errorf("unsafe deletion path: %s", dir)
		}

		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

// Helper functions (simplified)

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	os.MkdirAll(dst, 0755)
	for _, entry := range entries {
		if entry.IsDir() {
			copyDir(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()))
		} else {
			copyFile(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()))
		}
	}
	return nil
}

func unzip(ctx context.Context, src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fpath := filepath.Join(dest, f.Name)
			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
