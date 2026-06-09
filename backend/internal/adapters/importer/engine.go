package importer

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/xuri/excelize/v2"
)

// Engine implements ports.ImportEngine.
type Engine struct {
	uploadDir string
}

// NewEngine creates a new import engine.
func NewEngine() *Engine {
	return &Engine{uploadDir: defaultUploadDir()}
}

// NewEngineWithDir creates an engine with a custom upload directory.
func NewEngineWithDir(dir string) *Engine {
	return &Engine{uploadDir: dir}
}

// UploadDir returns the directory where uploaded files are stored.
func (e *Engine) UploadDir() string {
	return e.uploadDir
}

// ParseFile reads an Excel (.xlsx/.xls) or CSV file and returns headers + rows.
func (e *Engine) ParseFile(filePath string) ([]string, [][]string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".xlsx", ".xls":
		return e.parseExcel(filePath)
	case ".csv":
		return e.parseCSV(filePath)
	default:
		return nil, nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// DetectEntity guesses the target entity from column headers.
func (e *Engine) DetectEntity(headers []string) string {
	normalized := make([]string, len(headers))
	for i, h := range headers {
		normalized[i] = strings.ToLower(strings.TrimSpace(h))
	}
	joined := strings.Join(normalized, " ")

	switch {
	case containsAny(joined, "salary", "basic", "deduction", "net pay", "gross"):
		return domain.ImportEntitySalary
	case containsAny(joined, "advance", "advance amount", "repayment"):
		return domain.ImportEntityAdvances
	case containsAny(joined, "expense", "vendor", "payment method", "expense date"):
		return domain.ImportEntityExpenses
	case containsAny(joined, "product", "sku", "stock", "mrp", "purchase price"):
		return domain.ImportEntityProducts
	case containsAny(joined, "service", "duration", "service name", "price"):
		return domain.ImportEntityServices
	case containsAny(joined, "customer", "client", "dob", "date of birth"):
		return domain.ImportEntityCustomers
	case containsAny(joined, "staff", "employee", "designation", "joining"):
		return domain.ImportEntityStaff
	default:
		return domain.ImportEntityCustomers
	}
}

// SuggestMapping suggests column-to-field mappings based on header names.
func (e *Engine) SuggestMapping(headers []string, targetEntity string) []domain.ColumnMapping {
	fieldMap := getFieldMap(targetEntity)
	var mappings []domain.ColumnMapping

	for _, h := range headers {
		norm := normalizeHeader(h)
		targetField := ""
		for pattern, field := range fieldMap {
			if strings.Contains(norm, pattern) {
				targetField = field
				break
			}
		}
		if targetField != "" {
			mappings = append(mappings, domain.ColumnMapping{
				SourceColumn: h,
				TargetField:  targetField,
			})
		}
	}
	return mappings
}

// ValidateRow validates a single mapped row for a target entity.
func (e *Engine) ValidateRow(row map[string]string, targetEntity string) (bool, []string) {
	var errors []string

	required := getRequiredFields(targetEntity)
	for _, field := range required {
		val := strings.TrimSpace(row[field])
		if val == "" {
			errors = append(errors, fmt.Sprintf("missing required field: %s", field))
		}
	}

	// Validate specific fields
	if phone := row["phone"]; phone != "" {
		if !isValidPhone(phone) {
			errors = append(errors, "invalid phone number")
		}
	}

	if email := row["email"]; email != "" {
		if !isValidEmail(email) {
			errors = append(errors, "invalid email format")
		}
	}

	// Validate numeric fields
	numericFields := getNumericFields(targetEntity)
	for _, field := range numericFields {
		if val := row[field]; val != "" {
			if !isNumeric(val) {
				errors = append(errors, fmt.Sprintf("invalid number for field: %s", field))
			}
		}
	}

	// Validate date fields
	dateFields := getDateFields(targetEntity)
	for _, field := range dateFields {
		if val := row[field]; val != "" {
			if !isValidDate(val) {
				errors = append(errors, fmt.Sprintf("invalid date for field: %s", field))
			}
		}
	}

	return len(errors) == 0, errors
}

func (e *Engine) parseExcel(filePath string) ([]string, [][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	// Use the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, nil, fmt.Errorf("no sheets found")
	}

	allRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("read rows: %w", err)
	}

	if len(allRows) == 0 {
		return nil, nil, fmt.Errorf("empty sheet")
	}

	headers := allRows[0]
	var rows [][]string
	if len(allRows) > 1 {
		rows = allRows[1:]
	}

	return headers, rows, nil
}

func (e *Engine) parseCSV(filePath string) ([]string, [][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("read csv: %w", err)
	}

	if len(records) == 0 {
		return nil, nil, fmt.Errorf("empty csv")
	}

	headers := records[0]
	var rows [][]string
	if len(records) > 1 {
		rows = records[1:]
	}

	return headers, rows, nil
}

func getFieldMap(entity string) map[string]string {
	switch entity {
	case domain.ImportEntityStaff:
		return map[string]string{
			"name": "name", "phone": "phone", "email": "email",
			"designation": "designation", "joining": "date_of_joining",
			"salary": "base_salary",
		}
	case domain.ImportEntityCustomers:
		return map[string]string{
			"name": "name", "phone": "phone", "email": "email",
			"dob": "date_of_birth", "birth": "date_of_birth",
		}
	case domain.ImportEntityServices:
		return map[string]string{
			"name": "name", "service": "name", "price": "price",
			"duration": "duration", "category": "category",
		}
	case domain.ImportEntityProducts:
		return map[string]string{
			"name": "name", "product": "name", "sku": "sku",
			"brand": "brand", "category": "category",
			"mrp": "mrp", "price": "selling_price",
			"purchase": "purchase_price", "stock": "stock_quantity",
		}
	case domain.ImportEntityExpenses:
		return map[string]string{
			"amount": "amount", "date": "expense_date",
			"vendor": "vendor_name", "description": "description",
			"category": "category", "payment": "payment_method",
		}
	case domain.ImportEntityAdvances:
		return map[string]string{
			"staff": "staff_name", "amount": "amount",
			"date": "advance_date", "reason": "reason",
		}
	case domain.ImportEntitySalary:
		return map[string]string{
			"staff": "staff_name", "month": "month",
			"basic": "base_salary", "incentive": "incentive",
			"advance": "advance_deduction", "deduction": "deductions",
			"net": "net_pay",
		}
	default:
		return map[string]string{}
	}
}

func getRequiredFields(entity string) []string {
	switch entity {
	case domain.ImportEntityStaff:
		return []string{"name"}
	case domain.ImportEntityCustomers:
		return []string{"name"}
	case domain.ImportEntityServices:
		return []string{"name", "price"}
	case domain.ImportEntityProducts:
		return []string{"name"}
	case domain.ImportEntityExpenses:
		return []string{"amount", "expense_date"}
	case domain.ImportEntityAdvances:
		return []string{"staff_name", "amount"}
	case domain.ImportEntitySalary:
		return []string{"staff_name"}
	default:
		return nil
	}
}

func getNumericFields(entity string) []string {
	switch entity {
	case domain.ImportEntityStaff:
		return []string{"base_salary"}
	case domain.ImportEntityServices:
		return []string{"price", "duration"}
	case domain.ImportEntityProducts:
		return []string{"mrp", "selling_price", "purchase_price", "stock_quantity"}
	case domain.ImportEntityExpenses:
		return []string{"amount"}
	case domain.ImportEntityAdvances:
		return []string{"amount"}
	case domain.ImportEntitySalary:
		return []string{"base_salary", "incentive", "advance_deduction", "deductions", "net_pay"}
	default:
		return nil
	}
}

func getDateFields(entity string) []string {
	switch entity {
	case domain.ImportEntityStaff:
		return []string{"date_of_joining"}
	case domain.ImportEntityCustomers:
		return []string{"date_of_birth"}
	case domain.ImportEntityExpenses:
		return []string{"expense_date"}
	case domain.ImportEntityAdvances:
		return []string{"advance_date"}
	default:
		return nil
	}
}

func normalizeHeader(h string) string {
	s := strings.ToLower(strings.TrimSpace(h))
	// Remove non-alphanumeric except spaces
	result := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			return r
		}
		return ' '
	}, s)
	return strings.Join(strings.Fields(result), " ")
}

func containsAny(s string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}

var phoneRegex = regexp.MustCompile(`^[+]?[\d\s\-()]{7,15}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidPhone(phone string) bool {
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func isNumeric(s string) bool {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return true
	}
	for i, r := range s {
		if r == '.' || r == '-' {
			if r == '-' && i != 0 {
				return false
			}
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isValidDate(s string) bool {
	s = strings.TrimSpace(s)
	formats := []string{
		"2006-01-02", "02-01-2006", "01/02/2006", "02/01/2006",
		"2006/01/02", "Jan 2, 2006", "2 Jan 2006", time.RFC3339,
	}
	for _, f := range formats {
		if _, err := time.Parse(f, s); err == nil {
			return true
		}
	}
	return false
}

func defaultUploadDir() string {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	default:
		base, _ = os.UserHomeDir()
	}
	return filepath.Join(base, "SalonFlowTrack", "Imports")
}
