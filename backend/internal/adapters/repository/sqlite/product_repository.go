package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ProductRepository implements ports.ProductRepository.
type ProductRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *sql.DB, log *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, log: log}
}

// --- Products ---

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (id, product_code, name, category, brand, unit, sku, purchase_price,
		selling_price, current_stock, minimum_stock, maximum_stock, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID.String(), p.ProductCode, p.Name, p.Category, p.Brand, p.Unit, p.SKU,
		p.PurchasePrice, p.SellingPrice, p.CurrentStock, p.MinimumStock, p.MaximumStock,
		p.Status, p.CreatedAt.Format(time.RFC3339), p.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create product", err)
	}
	return nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, product_code, name, category, brand, unit, sku, purchase_price,
		selling_price, current_stock, minimum_stock, maximum_stock, status, created_at, updated_at
		FROM products WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())
	p, err := r.scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("product", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get product", err)
	}
	return p, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	query := `UPDATE products SET name = ?, category = ?, brand = ?, unit = ?, sku = ?,
		purchase_price = ?, selling_price = ?, minimum_stock = ?, maximum_stock = ?,
		status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		p.Name, p.Category, p.Brand, p.Unit, p.SKU,
		p.PurchasePrice, p.SellingPrice, p.MinimumStock, p.MaximumStock,
		p.Status, time.Now().UTC().Format(time.RFC3339), p.ID.String())
	if err != nil {
		return apperror.Database("update product", err)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id.String())
	if err != nil {
		return apperror.Database("delete product", err)
	}
	return nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, filter ports.ProductFilter) ([]domain.Product, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.Category != "" {
		where = append(where, "category = ?")
		args = append(args, filter.Category)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		where = append(where, "(name LIKE ? OR product_code LIKE ? OR brand LIKE ?)")
		s := "%" + filter.Search + "%"
		args = append(args, s, s, s)
	}
	if filter.LowStock {
		where = append(where, "current_stock < minimum_stock AND status = 'active'")
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count products", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	query := fmt.Sprintf(`SELECT id, product_code, name, category, brand, unit, sku, purchase_price,
		selling_price, current_stock, minimum_stock, maximum_stock, status, created_at, updated_at
		FROM products WHERE %s ORDER BY name ASC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list products", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		p, err := r.scanProductRow(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan product", err)
		}
		products = append(products, *p)
	}
	return products, total, nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID uuid.UUID, delta float64) error {
	query := `UPDATE products SET current_stock = current_stock + ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, delta, time.Now().UTC().Format(time.RFC3339), productID.String())
	if err != nil {
		return apperror.Database("update stock", err)
	}
	return nil
}

// --- Stock Transactions ---

func (r *ProductRepository) CreateStockTransaction(ctx context.Context, st *domain.StockTransaction) error {
	query := `INSERT INTO stock_transactions (id, product_id, transaction_type, quantity, unit_cost,
		reference_type, reference_id, notes, transaction_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		st.ID.String(), st.ProductID.String(), st.TransactionType, st.Quantity, st.UnitCost,
		st.ReferenceType, st.ReferenceID, st.Notes, st.TransactionDate,
		st.CreatedAt.Format(time.RFC3339), st.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create stock transaction", err)
	}
	return nil
}

func (r *ProductRepository) ListStockTransactions(ctx context.Context, filter ports.StockTransactionFilter) ([]domain.StockTransaction, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.ProductID != "" {
		where = append(where, "st.product_id = ?")
		args = append(args, filter.ProductID)
	}
	if filter.TransactionType != "" {
		where = append(where, "st.transaction_type = ?")
		args = append(args, filter.TransactionType)
	}
	if filter.DateFrom != "" {
		where = append(where, "st.transaction_date >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		where = append(where, "st.transaction_date <= ?")
		args = append(args, filter.DateTo)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stock_transactions st WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count stock transactions", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	query := fmt.Sprintf(`SELECT st.id, st.product_id, p.name, st.transaction_type, st.quantity, st.unit_cost,
		st.reference_type, st.reference_id, st.notes, st.transaction_date, st.created_at, st.updated_at
		FROM stock_transactions st
		JOIN products p ON p.id = st.product_id
		WHERE %s ORDER BY st.transaction_date DESC, st.created_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list stock transactions", err)
	}
	defer rows.Close()

	var txns []domain.StockTransaction
	for rows.Next() {
		var st domain.StockTransaction
		var idStr, prodIDStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &prodIDStr, &st.ProductName, &st.TransactionType, &st.Quantity,
			&st.UnitCost, &st.ReferenceType, &st.ReferenceID, &st.Notes, &st.TransactionDate,
			&createdStr, &updatedStr)
		if err != nil {
			return nil, 0, apperror.Database("scan stock transaction", err)
		}
		st.ID, _ = uuid.Parse(idStr)
		st.ProductID, _ = uuid.Parse(prodIDStr)
		st.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		st.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		txns = append(txns, st)
	}
	return txns, total, nil
}

// --- Purchases ---

func (r *ProductRepository) CreatePurchaseEntry(ctx context.Context, pe *domain.PurchaseEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Database("begin purchase tx", err)
	}
	defer tx.Rollback()

	entryQuery := `INSERT INTO purchase_entries (id, purchase_number, vendor_name, invoice_number,
		purchase_date, total_amount, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, entryQuery,
		pe.ID.String(), pe.PurchaseNumber, pe.VendorName, pe.InvoiceNumber,
		pe.PurchaseDate, pe.TotalAmount, pe.Notes,
		pe.CreatedAt.Format(time.RFC3339), pe.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return apperror.Database("create purchase entry", err)
	}

	itemQuery := `INSERT INTO purchase_items (id, purchase_entry_id, product_id, quantity, unit_price, line_total, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	for _, item := range pe.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID.String(), pe.ID.String(), item.ProductID.String(),
			item.Quantity, item.UnitPrice, item.LineTotal,
			item.CreatedAt.Format(time.RFC3339), item.UpdatedAt.Format(time.RFC3339))
		if err != nil {
			return apperror.Database("create purchase item", err)
		}
	}

	return tx.Commit()
}

func (r *ProductRepository) GetPurchaseByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseEntry, error) {
	query := `SELECT id, purchase_number, vendor_name, invoice_number, purchase_date, total_amount, notes, created_at, updated_at
		FROM purchase_entries WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var pe domain.PurchaseEntry
	var idStr, createdStr, updatedStr string
	err := row.Scan(&idStr, &pe.PurchaseNumber, &pe.VendorName, &pe.InvoiceNumber,
		&pe.PurchaseDate, &pe.TotalAmount, &pe.Notes, &createdStr, &updatedStr)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("purchase_entry", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get purchase entry", err)
	}
	pe.ID, _ = uuid.Parse(idStr)
	pe.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	pe.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	// Load items
	itemQuery := `SELECT pi.id, pi.purchase_entry_id, pi.product_id, p.name, pi.quantity, pi.unit_price, pi.line_total, pi.created_at, pi.updated_at
		FROM purchase_items pi JOIN products p ON p.id = pi.product_id WHERE pi.purchase_entry_id = ?`
	rows, err := r.db.QueryContext(ctx, itemQuery, id.String())
	if err != nil {
		return nil, apperror.Database("get purchase items", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.PurchaseItem
		var iIDStr, peIDStr, prodIDStr, iCreatedStr, iUpdatedStr string
		err := rows.Scan(&iIDStr, &peIDStr, &prodIDStr, &item.ProductName, &item.Quantity, &item.UnitPrice, &item.LineTotal, &iCreatedStr, &iUpdatedStr)
		if err != nil {
			return nil, apperror.Database("scan purchase item", err)
		}
		item.ID, _ = uuid.Parse(iIDStr)
		item.PurchaseEntryID, _ = uuid.Parse(peIDStr)
		item.ProductID, _ = uuid.Parse(prodIDStr)
		item.CreatedAt, _ = time.Parse(time.RFC3339, iCreatedStr)
		item.UpdatedAt, _ = time.Parse(time.RFC3339, iUpdatedStr)
		pe.Items = append(pe.Items, item)
	}

	return &pe, nil
}

func (r *ProductRepository) ListPurchases(ctx context.Context, dateFrom, dateTo string, limit, offset int) ([]domain.PurchaseEntry, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if dateFrom != "" {
		where = append(where, "purchase_date >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		where = append(where, "purchase_date <= ?")
		args = append(args, dateTo)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM purchase_entries WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count purchases", err)
	}

	if limit <= 0 {
		limit = 20
	}
	query := fmt.Sprintf(`SELECT id, purchase_number, vendor_name, invoice_number, purchase_date, total_amount, notes, created_at, updated_at
		FROM purchase_entries WHERE %s ORDER BY purchase_date DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list purchases", err)
	}
	defer rows.Close()

	var entries []domain.PurchaseEntry
	for rows.Next() {
		var pe domain.PurchaseEntry
		var idStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &pe.PurchaseNumber, &pe.VendorName, &pe.InvoiceNumber,
			&pe.PurchaseDate, &pe.TotalAmount, &pe.Notes, &createdStr, &updatedStr)
		if err != nil {
			return nil, 0, apperror.Database("scan purchase entry", err)
		}
		pe.ID, _ = uuid.Parse(idStr)
		pe.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		pe.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		entries = append(entries, pe)
	}
	return entries, total, nil
}

// --- Reporting ---

func (r *ProductRepository) GetLowStockProducts(ctx context.Context) ([]domain.LowStockItem, error) {
	query := `SELECT id, product_code, name, category, current_stock, minimum_stock, (minimum_stock - current_stock)
		FROM products WHERE current_stock < minimum_stock AND status = 'active' ORDER BY (minimum_stock - current_stock) DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.Database("get low stock", err)
	}
	defer rows.Close()

	var items []domain.LowStockItem
	for rows.Next() {
		var item domain.LowStockItem
		err := rows.Scan(&item.ProductID, &item.ProductCode, &item.ProductName, &item.Category,
			&item.CurrentStock, &item.MinimumStock, &item.Deficit)
		if err != nil {
			return nil, apperror.Database("scan low stock item", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ProductRepository) GetInventoryStats(ctx context.Context) (*domain.InventoryStats, error) {
	stats := &domain.InventoryStats{}

	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products`).Scan(&stats.TotalProducts)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE status = 'active'`).Scan(&stats.ActiveProducts)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE current_stock < minimum_stock AND status = 'active'`).Scan(&stats.LowStockCount)
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(current_stock * purchase_price), 0) FROM products WHERE status = 'active'`).Scan(&stats.TotalValue)

	// This month's purchases
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount), 0) FROM purchase_entries WHERE purchase_date >= ?`, monthStart).Scan(&stats.TotalPurchases)

	return stats, nil
}

func (r *ProductRepository) GetInventoryValue(ctx context.Context) (float64, error) {
	var value float64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(current_stock * purchase_price), 0) FROM products WHERE status = 'active'`).Scan(&value)
	if err != nil {
		return 0, apperror.Database("get inventory value", err)
	}
	return value, nil
}

// --- Code Generation ---

func (r *ProductRepository) NextProductCode(ctx context.Context) (string, error) {
	return r.nextCode(ctx, "PRD")
}

func (r *ProductRepository) NextPurchaseNumber(ctx context.Context) (string, error) {
	return r.nextCode(ctx, "PUR")
}

func (r *ProductRepository) nextCode(ctx context.Context, prefix string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", apperror.Database("begin tx for code", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE product_code_seq SET seq = seq + 1 WHERE prefix = ?`, prefix)
	if err != nil {
		return "", apperror.Database("increment seq", err)
	}

	var seq int
	err = tx.QueryRowContext(ctx, `SELECT seq FROM product_code_seq WHERE prefix = ?`, prefix).Scan(&seq)
	if err != nil {
		return "", apperror.Database("read seq", err)
	}

	if err := tx.Commit(); err != nil {
		return "", apperror.Database("commit seq", err)
	}

	return fmt.Sprintf("%s-%06d", prefix, seq), nil
}

// --- Scan helpers ---

func (r *ProductRepository) scanProduct(row *sql.Row) (*domain.Product, error) {
	var p domain.Product
	var idStr, createdStr, updatedStr string
	err := row.Scan(&idStr, &p.ProductCode, &p.Name, &p.Category, &p.Brand, &p.Unit, &p.SKU,
		&p.PurchasePrice, &p.SellingPrice, &p.CurrentStock, &p.MinimumStock, &p.MaximumStock,
		&p.Status, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	p.ID, _ = uuid.Parse(idStr)
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &p, nil
}

func (r *ProductRepository) scanProductRow(rows *sql.Rows) (*domain.Product, error) {
	var p domain.Product
	var idStr, createdStr, updatedStr string
	err := rows.Scan(&idStr, &p.ProductCode, &p.Name, &p.Category, &p.Brand, &p.Unit, &p.SKU,
		&p.PurchasePrice, &p.SellingPrice, &p.CurrentStock, &p.MinimumStock, &p.MaximumStock,
		&p.Status, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}
	p.ID, _ = uuid.Parse(idStr)
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &p, nil
}
