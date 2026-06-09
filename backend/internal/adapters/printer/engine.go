package printer

import (
	"fmt"
	"strings"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ESC/POS command constants.
const (
	ESC          = 0x1B
	GS           = 0x1D
	LF           = 0x0A
	CUT          = 0x6D // full cut
	ALIGN_LEFT   = 0x00
	ALIGN_CENTER = 0x01
	ALIGN_RIGHT  = 0x02
)

// Engine handles receipt formatting and ESC/POS generation.
type Engine struct{}

// NewEngine creates a new print engine.
func NewEngine() *Engine {
	return &Engine{}
}

// FormatReceipt generates a plain text receipt for a given paper width.
func (e *Engine) FormatReceipt(data *domain.ReceiptData, paperWidth string) string {
	width := e.charWidth(paperWidth)
	var sb strings.Builder

	// Header
	sb.WriteString(e.center(data.SalonName, width) + "\n")
	if data.Address != "" {
		sb.WriteString(e.center(data.Address, width) + "\n")
	}
	if data.GSTIN != "" {
		sb.WriteString(e.center("GSTIN: "+data.GSTIN, width) + "\n")
	}
	sb.WriteString(strings.Repeat("-", width) + "\n")

	// Invoice info
	sb.WriteString(fmt.Sprintf("Invoice: %s\n", data.InvoiceNumber))
	sb.WriteString(fmt.Sprintf("Date: %s\n", data.Date))
	if data.CustomerName != "" {
		sb.WriteString(fmt.Sprintf("Customer: %s\n", data.CustomerName))
	}
	if data.CustomerPhone != "" {
		sb.WriteString(fmt.Sprintf("Phone: %s\n", data.CustomerPhone))
	}
	sb.WriteString(strings.Repeat("-", width) + "\n")

	// Items header
	sb.WriteString(e.formatLine("Item", "Qty", "Price", "Total", width) + "\n")
	sb.WriteString(strings.Repeat("-", width) + "\n")

	// Items
	for _, item := range data.Items {
		sb.WriteString(e.formatLine(
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			fmt.Sprintf("%.0f", item.Price),
			fmt.Sprintf("%.0f", item.Total),
			width,
		) + "\n")
	}
	sb.WriteString(strings.Repeat("-", width) + "\n")

	// Totals
	sb.WriteString(e.rightAlign("Subtotal:", fmt.Sprintf("%.2f", data.Subtotal), width) + "\n")
	if data.CGST > 0 {
		sb.WriteString(e.rightAlign("CGST:", fmt.Sprintf("%.2f", data.CGST), width) + "\n")
	}
	if data.SGST > 0 {
		sb.WriteString(e.rightAlign("SGST:", fmt.Sprintf("%.2f", data.SGST), width) + "\n")
	}
	if data.IGST > 0 {
		sb.WriteString(e.rightAlign("IGST:", fmt.Sprintf("%.2f", data.IGST), width) + "\n")
	}
	if data.Discount > 0 {
		sb.WriteString(e.rightAlign("Discount:", fmt.Sprintf("-%.2f", data.Discount), width) + "\n")
	}
	sb.WriteString(strings.Repeat("=", width) + "\n")
	sb.WriteString(e.rightAlign("GRAND TOTAL:", fmt.Sprintf("%.2f", data.GrandTotal), width) + "\n")
	sb.WriteString(strings.Repeat("=", width) + "\n")

	// Payment
	if data.PaymentMethod != "" {
		sb.WriteString(fmt.Sprintf("Paid by: %s\n", data.PaymentMethod))
	}

	// Footer
	sb.WriteString("\n")
	if data.FooterText != "" {
		sb.WriteString(e.center(data.FooterText, width) + "\n")
	}

	return sb.String()
}

// GenerateESCPOS generates ESC/POS byte commands for thermal printers.
func (e *Engine) GenerateESCPOS(data *domain.ReceiptData, paperWidth string) *domain.ESCPOSCommand {
	var cmd []byte

	// Initialize printer
	cmd = append(cmd, ESC, '@') // Reset

	// Center alignment for header
	cmd = append(cmd, ESC, 'a', ALIGN_CENTER)

	// Bold + double size for salon name
	cmd = append(cmd, ESC, 'E', 1)   // Bold on
	cmd = append(cmd, GS, '!', 0x11) // Double width + height
	cmd = append(cmd, []byte(data.SalonName)...)
	cmd = append(cmd, LF)
	cmd = append(cmd, GS, '!', 0x00) // Normal size
	cmd = append(cmd, ESC, 'E', 0)   // Bold off

	// Address and GSTIN
	if data.Address != "" {
		cmd = append(cmd, []byte(data.Address)...)
		cmd = append(cmd, LF)
	}
	if data.GSTIN != "" {
		cmd = append(cmd, []byte("GSTIN: "+data.GSTIN)...)
		cmd = append(cmd, LF)
	}

	// Left alignment for body
	cmd = append(cmd, ESC, 'a', ALIGN_LEFT)

	// Separator
	width := e.charWidth(paperWidth)
	cmd = append(cmd, []byte(strings.Repeat("-", width))...)
	cmd = append(cmd, LF)

	// Invoice details
	cmd = append(cmd, []byte(fmt.Sprintf("Invoice: %s", data.InvoiceNumber))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte(fmt.Sprintf("Date: %s", data.Date))...)
	cmd = append(cmd, LF)
	if data.CustomerName != "" {
		cmd = append(cmd, []byte(fmt.Sprintf("Customer: %s", data.CustomerName))...)
		cmd = append(cmd, LF)
	}

	cmd = append(cmd, []byte(strings.Repeat("-", width))...)
	cmd = append(cmd, LF)

	// Items
	for _, item := range data.Items {
		line := e.formatLine(item.Name, fmt.Sprintf("%d", item.Quantity), fmt.Sprintf("%.0f", item.Price), fmt.Sprintf("%.0f", item.Total), width)
		cmd = append(cmd, []byte(line)...)
		cmd = append(cmd, LF)
	}

	cmd = append(cmd, []byte(strings.Repeat("-", width))...)
	cmd = append(cmd, LF)

	// Totals
	cmd = append(cmd, ESC, 'E', 1) // Bold
	cmd = append(cmd, []byte(e.rightAlign("TOTAL:", fmt.Sprintf("%.2f", data.GrandTotal), width))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, ESC, 'E', 0) // Bold off

	// Footer centered
	cmd = append(cmd, LF)
	cmd = append(cmd, ESC, 'a', ALIGN_CENTER)
	if data.FooterText != "" {
		cmd = append(cmd, []byte(data.FooterText)...)
		cmd = append(cmd, LF)
	}

	// Feed and cut
	cmd = append(cmd, LF, LF, LF)
	cmd = append(cmd, GS, 'V', CUT)

	return &domain.ESCPOSCommand{Commands: cmd}
}

// FormatTestPage generates a test page for printer verification.
func (e *Engine) FormatTestPage(printerName, paperWidth string) *domain.ESCPOSCommand {
	var cmd []byte
	cmd = append(cmd, ESC, '@') // Reset
	cmd = append(cmd, ESC, 'a', ALIGN_CENTER)
	cmd = append(cmd, ESC, 'E', 1)
	cmd = append(cmd, []byte("PRINTER TEST PAGE")...)
	cmd = append(cmd, LF)
	cmd = append(cmd, ESC, 'E', 0)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte(fmt.Sprintf("Printer: %s", printerName))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte(fmt.Sprintf("Paper: %s", paperWidth))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte(fmt.Sprintf("Time: %s", time.Now().Format("2006-01-02 15:04:05")))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, ESC, 'a', ALIGN_LEFT)

	width := e.charWidth(paperWidth)
	cmd = append(cmd, []byte(strings.Repeat("-", width))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte("Character test:")...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte("0123456789")...)
	cmd = append(cmd, LF)
	cmd = append(cmd, []byte(strings.Repeat("-", width))...)
	cmd = append(cmd, LF)
	cmd = append(cmd, ESC, 'a', ALIGN_CENTER)
	cmd = append(cmd, []byte("TEST COMPLETE")...)
	cmd = append(cmd, LF, LF, LF)
	cmd = append(cmd, GS, 'V', CUT)

	return &domain.ESCPOSCommand{Commands: cmd}
}

// charWidth returns the number of characters per line for a paper width.
func (e *Engine) charWidth(paperWidth string) int {
	switch paperWidth {
	case domain.PaperWidth58mm:
		return 32
	case domain.PaperWidth80mm:
		return 48
	case domain.PaperWidthA4:
		return 80
	default:
		return 48
	}
}

func (e *Engine) center(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	pad := (width - len(text)) / 2
	return strings.Repeat(" ", pad) + text
}

func (e *Engine) rightAlign(label, value string, width int) string {
	gap := width - len(label) - len(value)
	if gap < 1 {
		gap = 1
	}
	return label + strings.Repeat(" ", gap) + value
}

func (e *Engine) formatLine(name, qty, price, total string, width int) string {
	// name takes remaining space, qty/price/total are fixed width
	qtyW := 4
	priceW := 8
	totalW := 8
	nameW := width - qtyW - priceW - totalW
	if nameW < 4 {
		nameW = 4
	}
	if len(name) > nameW {
		name = name[:nameW]
	}

	return fmt.Sprintf("%-*s%*s%*s%*s", nameW, name, qtyW, qty, priceW, price, totalW, total)
}
