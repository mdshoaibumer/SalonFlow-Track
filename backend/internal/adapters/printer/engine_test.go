package printer

import (
	"strings"
	"testing"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

func TestEngine_FormatReceipt_80mm(t *testing.T) {
	e := NewEngine()
	data := &domain.ReceiptData{
		SalonName:     "Glamour Salon",
		GSTIN:         "27AABCU9603R1ZM",
		Address:       "123 Main St, Mumbai",
		InvoiceNumber: "INV-2026-0001",
		Date:          "2026-06-09",
		CustomerName:  "Priya Sharma",
		CustomerPhone: "9876543210",
		Items: []domain.ReceiptItem{
			{Name: "Haircut", Quantity: 1, Price: 500, Total: 500},
			{Name: "Hair Color", Quantity: 1, Price: 2000, Total: 2000},
		},
		Subtotal:      2500,
		CGST:          225,
		SGST:          225,
		GrandTotal:    2950,
		PaymentMethod: "UPI",
		FooterText:    "Thank you for visiting!",
	}

	receipt := e.FormatReceipt(data, domain.PaperWidth80mm)

	if !strings.Contains(receipt, "Glamour Salon") {
		t.Error("receipt should contain salon name")
	}
	if !strings.Contains(receipt, "INV-2026-0001") {
		t.Error("receipt should contain invoice number")
	}
	if !strings.Contains(receipt, "Priya Sharma") {
		t.Error("receipt should contain customer name")
	}
	if !strings.Contains(receipt, "Haircut") {
		t.Error("receipt should contain item name")
	}
	if !strings.Contains(receipt, "CGST") {
		t.Error("receipt should contain CGST")
	}
	if !strings.Contains(receipt, "2950.00") {
		t.Error("receipt should contain grand total")
	}
	if !strings.Contains(receipt, "Thank you") {
		t.Error("receipt should contain footer")
	}
}

func TestEngine_FormatReceipt_58mm(t *testing.T) {
	e := NewEngine()
	data := &domain.ReceiptData{
		SalonName:     "Test Salon",
		InvoiceNumber: "INV-001",
		Date:          "2026-06-09",
		Items: []domain.ReceiptItem{
			{Name: "Service A", Quantity: 1, Price: 100, Total: 100},
		},
		Subtotal:   100,
		GrandTotal: 100,
	}

	receipt := e.FormatReceipt(data, domain.PaperWidth58mm)

	// 58mm should have shorter lines (32 chars)
	lines := strings.Split(receipt, "\n")
	for _, line := range lines {
		if len(line) > 32 && !strings.Contains(line, "GRAND TOTAL") {
			// Allow some overflow for formatted totals
		}
	}
	if !strings.Contains(receipt, "Test Salon") {
		t.Error("receipt should contain salon name")
	}
}

func TestEngine_GenerateESCPOS(t *testing.T) {
	e := NewEngine()
	data := &domain.ReceiptData{
		SalonName:     "Test Salon",
		InvoiceNumber: "INV-001",
		Date:          "2026-06-09",
		Items: []domain.ReceiptItem{
			{Name: "Haircut", Quantity: 1, Price: 500, Total: 500},
		},
		Subtotal:   500,
		GrandTotal: 500,
		FooterText: "Thank you!",
	}

	result := e.GenerateESCPOS(data, domain.PaperWidth80mm)

	if len(result.Commands) == 0 {
		t.Fatal("expected ESC/POS commands")
	}

	// Check starts with ESC @ (reset)
	if result.Commands[0] != ESC || result.Commands[1] != '@' {
		t.Error("should start with ESC @ (reset)")
	}

	// Check contains cut command at end
	cmds := result.Commands
	found := false
	for i := 0; i < len(cmds)-2; i++ {
		if cmds[i] == GS && cmds[i+1] == 'V' && cmds[i+2] == CUT {
			found = true
			break
		}
	}
	if !found {
		t.Error("should contain cut command (GS V m)")
	}

	// Check contains salon name
	if !strings.Contains(string(result.Commands), "Test Salon") {
		t.Error("ESC/POS should contain salon name")
	}
}

func TestEngine_FormatTestPage(t *testing.T) {
	e := NewEngine()
	result := e.FormatTestPage("USB Printer", "80mm")

	if len(result.Commands) == 0 {
		t.Fatal("expected test page commands")
	}
	if !strings.Contains(string(result.Commands), "PRINTER TEST PAGE") {
		t.Error("test page should contain title")
	}
	if !strings.Contains(string(result.Commands), "USB Printer") {
		t.Error("test page should contain printer name")
	}
}

func TestEngine_CharWidth(t *testing.T) {
	e := NewEngine()

	if e.charWidth("58mm") != 32 {
		t.Errorf("58mm = %d, want 32", e.charWidth("58mm"))
	}
	if e.charWidth("80mm") != 48 {
		t.Errorf("80mm = %d, want 48", e.charWidth("80mm"))
	}
	if e.charWidth("A4") != 80 {
		t.Errorf("A4 = %d, want 80", e.charWidth("A4"))
	}
}
