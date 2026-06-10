package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEngine_ParseCSV(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "test.csv")
	content := `Name,Phone,Email,Designation
John Doe,9876543210,john@example.com,Stylist
Jane Smith,9876543211,jane@example.com,Manager
`
	os.WriteFile(csvPath, []byte(content), 0644)

	e := NewEngineWithDir(dir)
	headers, rows, err := e.ParseFile(csvPath)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(headers) != 4 {
		t.Errorf("headers len = %d, want 4", len(headers))
	}
	if len(rows) != 2 {
		t.Errorf("rows len = %d, want 2", len(rows))
	}
	if headers[0] != "Name" {
		t.Errorf("headers[0] = %q, want Name", headers[0])
	}
}

func TestEngine_DetectEntity(t *testing.T) {
	e := NewEngine()

	tests := []struct {
		headers  []string
		expected string
	}{
		{[]string{"Staff Name", "Designation", "Joining Date"}, "staff"},
		{[]string{"Customer Name", "Phone", "DOB"}, "customers"},
		{[]string{"Service Name", "Price", "Duration"}, "services"},
		{[]string{"Product Name", "SKU", "MRP", "Stock"}, "products"},
		{[]string{"Expense Date", "Amount", "Vendor"}, "expenses"},
		{[]string{"Staff", "Advance Amount", "Date"}, "advances"},
		{[]string{"Staff", "Basic Salary", "Deduction", "Net Pay"}, "salary"},
	}

	for _, tt := range tests {
		got := e.DetectEntity(tt.headers)
		if got != tt.expected {
			t.Errorf("DetectEntity(%v) = %q, want %q", tt.headers, got, tt.expected)
		}
	}
}

func TestEngine_SuggestMapping(t *testing.T) {
	e := NewEngine()

	headers := []string{"Staff Name", "Phone Number", "Email ID", "Designation"}
	mappings := e.SuggestMapping(headers, "staff")

	if len(mappings) == 0 {
		t.Fatal("expected some mappings")
	}

	// Check that at least name is mapped
	found := false
	for _, m := range mappings {
		if m.TargetField == "name" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'name' field to be mapped")
	}
}

func TestEngine_ValidateRow(t *testing.T) {
	e := NewEngine()

	// Valid staff row
	row := map[string]string{"name": "John Doe", "phone": "9876543210", "email": "john@test.com"}
	valid, errs := e.ValidateRow(row, "staff")
	if !valid {
		t.Errorf("valid staff row should pass: %v", errs)
	}

	// Invalid - missing required field
	row2 := map[string]string{"phone": "9876543210"}
	valid, errs = e.ValidateRow(row2, "staff")
	if valid {
		t.Error("row without name should fail")
	}
	if len(errs) == 0 {
		t.Error("expected errors")
	}

	// Invalid phone
	row3 := map[string]string{"name": "Test", "phone": "abc"}
	valid, _ = e.ValidateRow(row3, "staff")
	if valid {
		t.Error("invalid phone should fail")
	}

	// Invalid email
	row4 := map[string]string{"name": "Test", "email": "not-an-email"}
	valid, _ = e.ValidateRow(row4, "staff")
	if valid {
		t.Error("invalid email should fail")
	}
}

func TestEngine_ValidateNumeric(t *testing.T) {
	e := NewEngine()

	// Valid service row
	row := map[string]string{"name": "Haircut", "price": "500"}
	valid, _ := e.ValidateRow(row, "services")
	if !valid {
		t.Error("valid service row should pass")
	}

	// Invalid number
	row2 := map[string]string{"name": "Haircut", "price": "abc"}
	valid, errs := e.ValidateRow(row2, "services")
	if valid {
		t.Error("invalid price should fail")
	}
	if len(errs) == 0 {
		t.Error("expected validation errors")
	}
}

func TestEngine_UnsupportedFormat(t *testing.T) {
	e := NewEngine()
	_, _, err := e.ParseFile("/tmp/test.pdf")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
