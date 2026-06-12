package usecase

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DiagnosticsUseCase provides system diagnostics and log export.
type DiagnosticsUseCase struct {
	db         *sql.DB
	log        *slog.Logger
	appVersion string
	dbPath     string
	logDir     string
}

// NewDiagnosticsUseCase creates a new DiagnosticsUseCase.
func NewDiagnosticsUseCase(db *sql.DB, log *slog.Logger, appVersion, dbPath, logDir string) *DiagnosticsUseCase {
	return &DiagnosticsUseCase{
		db:         db,
		log:        log,
		appVersion: appVersion,
		dbPath:     dbPath,
		logDir:     logDir,
	}
}

// DiagnosticsInfo holds system health and configuration data.
type DiagnosticsInfo struct {
	AppVersion     string  `json:"app_version"`
	GoVersion      string  `json:"go_version"`
	OS             string  `json:"os"`
	Arch           string  `json:"arch"`
	DatabasePath   string  `json:"database_path"`
	DatabaseSize   int64   `json:"database_size_bytes"`
	DBVersion      string  `json:"db_version"`
	LogDirectory   string  `json:"log_directory"`
	NumCPU         int     `json:"num_cpu"`
	NumGoroutine   int     `json:"num_goroutine"`
	MemAllocMB     float64 `json:"mem_alloc_mb"`
	MemTotalMB     float64 `json:"mem_total_alloc_mb"`
	Uptime         string  `json:"uptime"`
	LastBackup     string  `json:"last_backup"`
	TotalUsers     int     `json:"total_users"`
	TotalInvoices  int     `json:"total_invoices"`
	TotalCustomers int     `json:"total_customers"`
}

var startTime = time.Now()

// GetDiagnostics returns current system diagnostics.
func (uc *DiagnosticsUseCase) GetDiagnostics(ctx context.Context) (*DiagnosticsInfo, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	info := &DiagnosticsInfo{
		AppVersion:   uc.appVersion,
		GoVersion:    runtime.Version(),
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		DatabasePath: uc.dbPath,
		LogDirectory: uc.logDir,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemAllocMB:   float64(memStats.Alloc) / 1024 / 1024,
		MemTotalMB:   float64(memStats.TotalAlloc) / 1024 / 1024,
		Uptime:       time.Since(startTime).Round(time.Second).String(),
	}

	// Database size
	if fi, err := os.Stat(uc.dbPath); err == nil {
		info.DatabaseSize = fi.Size()
	}

	// DB migration version
	var dbVersion string
	row := uc.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), '0') FROM schema_migrations")
	if err := row.Scan(&dbVersion); err == nil {
		info.DBVersion = dbVersion
	}

	// Counts
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&info.TotalUsers)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM invoices").Scan(&info.TotalInvoices)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM customers").Scan(&info.TotalCustomers)

	// Last backup
	var lastBackup sql.NullString
	uc.db.QueryRowContext(ctx, "SELECT MAX(created_at) FROM backups WHERE status = 'completed'").Scan(&lastBackup)
	if lastBackup.Valid {
		info.LastBackup = lastBackup.String
	}

	return info, nil
}

// ExportDiagnosticsBundle creates a ZIP file with all diagnostic data.
func (uc *DiagnosticsUseCase) ExportDiagnosticsBundle(ctx context.Context) (string, error) {
	// Create output file in temp directory
	timestamp := time.Now().Format("20060102-150405")
	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("SalonFlowTrack-Diagnostics-%s.zip", timestamp))

	zipFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("create zip file: %w", err)
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	// 1. Add diagnostics info
	diag, err := uc.GetDiagnostics(ctx)
	if err != nil {
		return "", fmt.Errorf("get diagnostics: %w", err)
	}

	diagJSON, _ := json.MarshalIndent(diag, "", "  ")
	if err := addToZip(w, "diagnostics.json", diagJSON); err != nil {
		return "", err
	}

	// 2. Add log files
	if uc.logDir != "" {
		logFiles := []string{"app.log", "error.log", "audit.log", "security.log", "performance.log"}
		for _, lf := range logFiles {
			path := filepath.Join(uc.logDir, lf)
			if data, err := readLastNBytes(path, 5*1024*1024); err == nil { // Last 5MB
				if err := addToZip(w, "logs/"+lf, data); err != nil {
					uc.log.Warn("failed to add log to bundle", "file", lf, "error", err)
				}
			}
		}
	}

	// 3. Add database metadata (not the actual DB)
	dbMeta := map[string]interface{}{
		"path":              uc.dbPath,
		"size_bytes":        diag.DatabaseSize,
		"migration_version": diag.DBVersion,
	}

	// Table row counts
	tables := []string{"users", "staff", "customers", "services", "invoices", "payments",
		"expenses", "products", "salary_slips", "audit_logs", "sessions"}
	tableCounts := make(map[string]int)
	for _, t := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", t)
		if err := uc.db.QueryRowContext(ctx, query).Scan(&count); err == nil {
			tableCounts[t] = count
		}
	}
	dbMeta["table_counts"] = tableCounts

	dbMetaJSON, _ := json.MarshalIndent(dbMeta, "", "  ")
	if err := addToZip(w, "database-metadata.json", dbMetaJSON); err != nil {
		return "", err
	}

	// 4. Add recent audit logs
	rows, err := uc.db.QueryContext(ctx, `
		SELECT id, timestamp, username, action, module, entity_type, entity_id, description, severity
		FROM audit_logs ORDER BY timestamp DESC LIMIT 500`)
	if err == nil {
		defer rows.Close()
		var auditEntries []map[string]string
		for rows.Next() {
			var id, ts, username, action, module, entityType, entityID, desc, severity string
			if err := rows.Scan(&id, &ts, &username, &action, &module, &entityType, &entityID, &desc, &severity); err == nil {
				auditEntries = append(auditEntries, map[string]string{
					"id": id, "timestamp": ts, "username": username, "action": action,
					"module": module, "entity_type": entityType, "entity_id": entityID,
					"description": desc, "severity": severity,
				})
			}
		}
		auditJSON, _ := json.MarshalIndent(auditEntries, "", "  ")
		addToZip(w, "recent-audit-logs.json", auditJSON)
	}

	// 5. System info
	sysInfo := map[string]interface{}{
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"go_version":  runtime.Version(),
		"num_cpu":     runtime.NumCPU(),
		"app_version": uc.appVersion,
		"export_time": time.Now().Format(time.RFC3339),
		"hostname":    getHostname(),
	}
	sysJSON, _ := json.MarshalIndent(sysInfo, "", "  ")
	addToZip(w, "system-info.json", sysJSON)

	uc.log.Info("diagnostics bundle exported", "path", outputPath)
	return outputPath, nil
}

func addToZip(w *zip.Writer, name string, data []byte) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func readLastNBytes(path string, maxBytes int64) ([]byte, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	size := fi.Size()
	if size <= maxBytes {
		return io.ReadAll(f)
	}

	// Seek to the last N bytes
	if _, err := f.Seek(size-maxBytes, 0); err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}
