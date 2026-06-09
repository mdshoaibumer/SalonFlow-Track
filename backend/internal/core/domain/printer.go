package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Paper widths.
const (
	PaperWidth58mm = "58mm"
	PaperWidth80mm = "80mm"
	PaperWidthA4   = "A4"
)

// Document types for printing.
const (
	PrintDocInvoice    = "invoice"
	PrintDocReceipt    = "receipt"
	PrintDocExpense    = "expense"
	PrintDocSalarySlip = "salary_slip"
	PrintDocPurchase   = "purchase"
)

// Print job statuses.
const (
	PrintStatusQueued    = "queued"
	PrintStatusPrinting  = "printing"
	PrintStatusCompleted = "completed"
	PrintStatusFailed    = "failed"
)

// PrinterSettings represents printer configuration.
type PrinterSettings struct {
	ID             uuid.UUID `json:"id"`
	DefaultPrinter string    `json:"default_printer"`
	PaperWidth     string    `json:"paper_width"`
	MarginTop      int       `json:"margin_top"`
	MarginBottom   int       `json:"margin_bottom"`
	MarginLeft     int       `json:"margin_left"`
	MarginRight    int       `json:"margin_right"`
	HeaderText     string    `json:"header_text"`
	FooterText     string    `json:"footer_text"`
	ShowLogo       bool      `json:"show_logo"`
	ShowQR         bool      `json:"show_qr"`
	UPIID          string    `json:"upi_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewPrinterSettings creates default printer settings.
func NewPrinterSettings() *PrinterSettings {
	now := time.Now().UTC()
	return &PrinterSettings{
		ID:           uid.New(),
		PaperWidth:   PaperWidth80mm,
		MarginTop:    5,
		MarginBottom: 5,
		MarginLeft:   5,
		MarginRight:  5,
		FooterText:   "Thank you for visiting!",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// PrintJob represents a print job record.
type PrintJob struct {
	ID           uuid.UUID `json:"id"`
	DocumentType string    `json:"document_type"`
	DocumentID   string    `json:"document_id"`
	PrinterName  string    `json:"printer_name"`
	PaperWidth   string    `json:"paper_width"`
	Status       string    `json:"status"`
	Copies       int       `json:"copies"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewPrintJob creates a new print job.
func NewPrintJob(docType, docID, printerName, paperWidth string, copies int) *PrintJob {
	return &PrintJob{
		ID:           uid.New(),
		DocumentType: docType,
		DocumentID:   docID,
		PrinterName:  printerName,
		PaperWidth:   paperWidth,
		Status:       PrintStatusQueued,
		Copies:       copies,
		CreatedAt:    time.Now().UTC(),
	}
}

// ReceiptData holds the data needed to format a receipt.
type ReceiptData struct {
	SalonName     string        `json:"salon_name"`
	GSTIN         string        `json:"gstin"`
	Address       string        `json:"address"`
	InvoiceNumber string        `json:"invoice_number"`
	Date          string        `json:"date"`
	CustomerName  string        `json:"customer_name"`
	CustomerPhone string        `json:"customer_phone"`
	Items         []ReceiptItem `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	CGST          float64       `json:"cgst"`
	SGST          float64       `json:"sgst"`
	IGST          float64       `json:"igst"`
	Discount      float64       `json:"discount"`
	GrandTotal    float64       `json:"grand_total"`
	PaymentMethod string        `json:"payment_method"`
	FooterText    string        `json:"footer_text"`
	UPIID         string        `json:"upi_id"`
}

// ReceiptItem is a line item on a receipt.
type ReceiptItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}

// ESCPOSCommand represents raw ESC/POS commands for thermal printers.
type ESCPOSCommand struct {
	Commands []byte `json:"commands"`
}
