package backup

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// Engine implements ports.BackupEngine using file-copy approach for SQLite.
type Engine struct {
	backupDir string
}

// NewEngine creates a new backup engine.
func NewEngine() *Engine {
	dir := defaultBackupDir()
	_ = os.MkdirAll(dir, 0750)
	return &Engine{backupDir: dir}
}

// NewEngineWithDir creates a backup engine with a custom directory (for testing).
func NewEngineWithDir(dir string) *Engine {
	_ = os.MkdirAll(dir, 0750)
	return &Engine{backupDir: dir}
}

// BackupDir returns the backup storage directory.
func (e *Engine) BackupDir() string {
	return e.backupDir
}

// CreateBackup creates a backup of the SQLite database using VACUUM INTO.
// Returns file size and SHA-256 checksum.
func (e *Engine) CreateBackup(dbPath, destPath string) (int64, string, error) {
	// Ensure destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0750); err != nil {
		return 0, "", fmt.Errorf("create backup directory: %w", err)
	}

	// Use VACUUM INTO for a consistent backup (available since SQLite 3.27.0)
	srcDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", dbPath))
	if err != nil {
		return 0, "", fmt.Errorf("open source db: %w", err)
	}
	defer srcDB.Close()

	_, err = srcDB.Exec(fmt.Sprintf(`VACUUM INTO '%s'`, destPath))
	if err != nil {
		return 0, "", fmt.Errorf("vacuum into: %w", err)
	}

	// Get file info
	info, err := os.Stat(destPath)
	if err != nil {
		return 0, "", fmt.Errorf("stat backup: %w", err)
	}

	// Calculate checksum
	checksum, err := fileChecksum(destPath)
	if err != nil {
		return info.Size(), "", fmt.Errorf("checksum: %w", err)
	}

	return info.Size(), checksum, nil
}

// VerifyBackup verifies a backup file.
func (e *Engine) VerifyBackup(backupPath, expectedChecksum string) *domain.BackupVerification {
	v := &domain.BackupVerification{
		BackupID: backupPath,
		Status:   domain.BackupStatusVerified,
	}

	// Check file exists
	if _, err := os.Stat(backupPath); err != nil {
		v.FileExists = false
		v.Status = domain.BackupStatusCorrupted
		v.ErrorMessage = "file not found"
		return v
	}
	v.FileExists = true

	// Check can open as SQLite
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", backupPath))
	if err != nil {
		v.CanOpen = false
		v.Status = domain.BackupStatusCorrupted
		v.ErrorMessage = "cannot open as database"
		return v
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		v.CanOpen = false
		v.Status = domain.BackupStatusCorrupted
		v.ErrorMessage = "cannot ping database"
		return v
	}
	v.CanOpen = true

	// Integrity check
	var result string
	err = db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil || result != "ok" {
		v.IntegrityOK = false
		v.Status = domain.BackupStatusCorrupted
		v.ErrorMessage = fmt.Sprintf("integrity check failed: %s", result)
		return v
	}
	v.IntegrityOK = true

	// Checksum verification
	if expectedChecksum != "" {
		actual, err := fileChecksum(backupPath)
		if err != nil || actual != expectedChecksum {
			v.ChecksumOK = false
			v.Status = domain.BackupStatusCorrupted
			v.ErrorMessage = "checksum mismatch"
			return v
		}
	}
	v.ChecksumOK = true

	return v
}

// RestoreBackup restores a backup to the database path.
// It copies the backup file over the existing database.
func (e *Engine) RestoreBackup(backupPath, dbPath string) error {
	// Verify backup first
	v := e.VerifyBackup(backupPath, "")
	if v.Status == domain.BackupStatusCorrupted {
		return fmt.Errorf("backup is corrupted: %s", v.ErrorMessage)
	}

	// Copy backup to db path
	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("open backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy backup: %w", err)
	}

	// Remove WAL and SHM files
	os.Remove(dbPath + "-wal")
	os.Remove(dbPath + "-shm")

	return nil
}

// GenerateBackupPath creates a backup file path in the structured directory.
func GenerateBackupPath(baseDir string) string {
	now := time.Now()
	dir := filepath.Join(baseDir, now.Format("2006"), now.Format("01"))
	name := fmt.Sprintf("backup_%s.db", now.Format("20060102_150405"))
	return filepath.Join(dir, name)
}

// GenerateBackupName creates a human-friendly backup name.
func GenerateBackupName(backupType string) string {
	now := time.Now()
	return fmt.Sprintf("%s_%s", backupType, now.Format("2006-01-02_15-04-05"))
}

func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func defaultBackupDir() string {
	var base string
	if runtime.GOOS == "windows" {
		base = os.Getenv("USERPROFILE")
	} else {
		base = os.Getenv("HOME")
	}
	return filepath.Join(base, "Documents", "SalonFlowTrack", "Backups")
}
