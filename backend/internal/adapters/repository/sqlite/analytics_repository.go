package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// AnalyticsRepository implements read-only analytics queries.
type AnalyticsRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewAnalyticsRepository creates a new AnalyticsRepository.
func NewAnalyticsRepository(db *sql.DB, log *slog.Logger) *AnalyticsRepository {
	return &AnalyticsRepository{db: db, log: log}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}
	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"

	// Today's revenue, customers, invoices
	row := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0), COUNT(*), COUNT(DISTINCT customer_id)
		FROM invoices WHERE DATE(created_at) = ? AND payment_status = 'paid'`, today)
	row.Scan(&stats.TodayRevenue, &stats.TodayInvoices, &stats.TodayCustomers)

	// Monthly revenue
	row = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0) FROM invoices
		WHERE DATE(created_at) >= ? AND payment_status = 'paid'`, monthStart)
	row.Scan(&stats.MonthlyRevenue)

	// Monthly expenses
	row = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM expenses
		WHERE expense_date >= ? AND status = 'approved'`, monthStart)
	row.Scan(&stats.MonthlyExpenses)

	stats.MonthlyProfit = stats.MonthlyRevenue - stats.MonthlyExpenses

	// Inventory value
	row = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(current_stock * purchase_price),0) FROM products WHERE status = 'active'`)
	row.Scan(&stats.InventoryValue)

	// Outstanding salary (unpaid records in current cycle)
	now := time.Now()
	row = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(net_salary),0) FROM salary_records sr
		JOIN salary_cycles sc ON sr.cycle_id = sc.id
		WHERE sc.month = ? AND sc.year = ? AND sr.payment_status = 'pending'`, int(now.Month()), now.Year())
	row.Scan(&stats.OutstandingSalary)

	// Outstanding advances
	row = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount - repaid_amount),0) FROM advances WHERE status IN ('approved','partially_repaid')`)
	row.Scan(&stats.OutstandingAdvance)

	// Low stock count
	row = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE status = 'active' AND current_stock < minimum_stock`)
	row.Scan(&stats.LowStockCount)

	return stats, nil
}

func (r *AnalyticsRepository) GetKPIMetrics(ctx context.Context, dateFrom, dateTo string) (*domain.KPIMetrics, error) {
	kpi := &domain.KPIMetrics{}

	// Calculate date range for previous period (same duration)
	from, _ := time.Parse("2006-01-02", dateFrom)
	to, _ := time.Parse("2006-01-02", dateTo)
	duration := to.Sub(from)
	prevFrom := from.Add(-duration).Format("2006-01-02")
	prevTo := from.Add(-24 * time.Hour).Format("2006-01-02")

	// Current period revenue
	var curRevenue, prevRevenue float64
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'`, dateFrom, dateTo).Scan(&curRevenue)
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'`, prevFrom, prevTo).Scan(&prevRevenue)

	if prevRevenue > 0 {
		kpi.RevenueGrowthPct = ((curRevenue - prevRevenue) / prevRevenue) * 100
	}

	// Customer growth
	var curCust, prevCust int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE DATE(created_at) >= ? AND DATE(created_at) <= ?`, dateFrom, dateTo).Scan(&curCust)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE DATE(created_at) >= ? AND DATE(created_at) <= ?`, prevFrom, prevTo).Scan(&prevCust)
	if prevCust > 0 {
		kpi.CustomerGrowthPct = (float64(curCust-prevCust) / float64(prevCust)) * 100
	}

	// Profit margin
	var expenses float64
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM expenses
		WHERE expense_date >= ? AND expense_date <= ? AND status = 'approved'`, dateFrom, dateTo).Scan(&expenses)
	if curRevenue > 0 {
		kpi.ProfitMarginPct = ((curRevenue - expenses) / curRevenue) * 100
	}

	// Average bill value
	var invoiceCount int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'`, dateFrom, dateTo).Scan(&invoiceCount)
	if invoiceCount > 0 {
		kpi.AverageBillValue = curRevenue / float64(invoiceCount)
	}

	// Repeat customer %
	var totalCusts, repeatCusts int
	r.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT customer_id) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ?`, dateFrom, dateTo).Scan(&totalCusts)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM (SELECT customer_id FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? GROUP BY customer_id HAVING COUNT(*) > 1)`, dateFrom, dateTo).Scan(&repeatCusts)
	if totalCusts > 0 {
		kpi.RepeatCustomerPct = (float64(repeatCusts) / float64(totalCusts)) * 100
	}

	// Staff productivity (% of active staff who generated revenue)
	var activeStaff, productiveStaff int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM staff WHERE status = 'active'`).Scan(&activeStaff)
	r.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT staff_id) FROM invoice_items ii
		JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ?`, dateFrom, dateTo).Scan(&productiveStaff)
	if activeStaff > 0 {
		kpi.StaffProductivity = (float64(productiveStaff) / float64(activeStaff)) * 100
	}

	return kpi, nil
}

func (r *AnalyticsRepository) GetRevenueReport(ctx context.Context, dateFrom, dateTo string, groupBy string) (*domain.RevenueReport, error) {
	report := &domain.RevenueReport{}

	// Total
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'`, dateFrom, dateTo).Scan(&report.TotalRevenue, &report.InvoiceCount)

	// Trend by groupBy (day/week/month)
	var dateExpr string
	switch groupBy {
	case "week":
		dateExpr = "strftime('%Y-W%W', created_at)"
	case "month":
		dateExpr = "strftime('%Y-%m', created_at)"
	case "year":
		dateExpr = "strftime('%Y', created_at)"
	default:
		dateExpr = "DATE(created_at)"
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s as period, COALESCE(SUM(total_amount),0)
		FROM invoices WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'
		GROUP BY period ORDER BY period`, dateExpr), dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p domain.RevenueTrendPoint
			rows.Scan(&p.Date, &p.Revenue)
			report.Trend = append(report.Trend, p)
		}
	}

	// By service (top 10)
	rows, err = r.db.QueryContext(ctx, `SELECT ii.service_name, COALESCE(SUM(ii.total),0) as rev
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY ii.service_name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.ByService = append(report.ByService, nv)
		}
	}

	// By staff (top 10)
	rows, err = r.db.QueryContext(ctx, `SELECT s.name, COALESCE(SUM(ii.total),0) as rev
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		JOIN staff s ON ii.staff_id = s.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY s.name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.ByStaff = append(report.ByStaff, nv)
		}
	}

	// By customer (top 10)
	rows, err = r.db.QueryContext(ctx, `SELECT c.name, COALESCE(SUM(i.total_amount),0) as rev
		FROM invoices i JOIN customers c ON i.customer_id = c.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY c.name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.ByCustomer = append(report.ByCustomer, nv)
		}
	}

	return report, nil
}

func (r *AnalyticsRepository) GetCustomerReport(ctx context.Context, dateFrom, dateTo string) (*domain.CustomerReport, error) {
	report := &domain.CustomerReport{}

	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&report.TotalCustomers)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE DATE(created_at) >= ? AND DATE(created_at) <= ?`, dateFrom, dateTo).Scan(&report.NewCustomers)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE strftime('%m-%d', date_of_birth) = strftime('%m-%d', 'now')`).Scan(&report.BirthdayToday)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE status = 'inactive'`).Scan(&report.InactiveCount)

	// Repeat customers (>1 invoice in period)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM (SELECT customer_id FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? GROUP BY customer_id HAVING COUNT(*) > 1)`, dateFrom, dateTo).Scan(&report.RepeatCustomers)

	// Top customers by revenue (top 10)
	rows, err := r.db.QueryContext(ctx, `SELECT c.name, COALESCE(SUM(i.total_amount),0) as rev
		FROM invoices i JOIN customers c ON i.customer_id = c.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY c.name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.TopCustomers = append(report.TopCustomers, nv)
		}
	}

	// Growth trend (monthly new customers, last 6 months)
	rows, err = r.db.QueryContext(ctx, `SELECT strftime('%Y-%m', created_at) as period, COUNT(*) as cnt
		FROM customers WHERE DATE(created_at) >= date('now', '-6 months')
		GROUP BY period ORDER BY period`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp domain.TrendPoint
			rows.Scan(&tp.Period, &tp.Value)
			report.GrowthTrend = append(report.GrowthTrend, tp)
		}
	}

	return report, nil
}

func (r *AnalyticsRepository) GetStaffReport(ctx context.Context, dateFrom, dateTo string) (*domain.StaffReport, error) {
	report := &domain.StaffReport{}

	// Top performers by revenue
	rows, err := r.db.QueryContext(ctx, `SELECT s.name, COALESCE(SUM(ii.total),0) as rev
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		JOIN staff s ON ii.staff_id = s.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY s.name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.TopPerformers = append(report.TopPerformers, nv)
		}
	}

	// Revenue by staff (same query, all)
	report.RevenueByStaff = report.TopPerformers

	// Customer count by staff
	rows, err = r.db.QueryContext(ctx, `SELECT s.name, COUNT(DISTINCT i.customer_id) as cnt
		FROM invoices i JOIN invoice_items ii ON ii.invoice_id = i.id
		JOIN staff s ON ii.staff_id = s.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ?
		GROUP BY s.name ORDER BY cnt DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.CustomersByStaff = append(report.CustomersByStaff, nv)
		}
	}

	// Commission earned by staff
	rows, err = r.db.QueryContext(ctx, `SELECT s.name, COALESCE(SUM(ct.amount),0) as comm
		FROM commission_transactions ct JOIN staff s ON ct.staff_id = s.id
		WHERE DATE(ct.created_at) >= ? AND DATE(ct.created_at) <= ?
		GROUP BY s.name ORDER BY comm DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.CommissionEarned = append(report.CommissionEarned, nv)
		}
	}

	// Salary cost in period
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(sr.net_salary),0) FROM salary_records sr
		JOIN salary_cycles sc ON sr.cycle_id = sc.id
		WHERE sc.start_date >= ? AND sc.end_date <= ?`, dateFrom, dateTo).Scan(&report.SalaryCost)

	return report, nil
}

func (r *AnalyticsRepository) GetServiceReport(ctx context.Context, dateFrom, dateTo string) (*domain.ServiceReport, error) {
	report := &domain.ServiceReport{}

	// Total bookings & avg value
	r.db.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(AVG(ii.total),0) FROM invoice_items ii
		JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'`, dateFrom, dateTo).Scan(&report.TotalBookings, &report.AvgServiceValue)

	// Top services
	rows, err := r.db.QueryContext(ctx, `SELECT ii.service_name, COUNT(*) as cnt
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY ii.service_name ORDER BY cnt DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.TopServices = append(report.TopServices, nv)
		}
	}

	// Least used
	rows, err = r.db.QueryContext(ctx, `SELECT ii.service_name, COUNT(*) as cnt
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY ii.service_name ORDER BY cnt ASC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.LeastUsed = append(report.LeastUsed, nv)
		}
	}

	// Revenue by service
	rows, err = r.db.QueryContext(ctx, `SELECT ii.service_name, COALESCE(SUM(ii.total),0) as rev
		FROM invoice_items ii JOIN invoices i ON ii.invoice_id = i.id
		WHERE DATE(i.created_at) >= ? AND DATE(i.created_at) <= ? AND i.payment_status = 'paid'
		GROUP BY ii.service_name ORDER BY rev DESC LIMIT 10`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.RevenueByService = append(report.RevenueByService, nv)
		}
	}

	return report, nil
}

func (r *AnalyticsRepository) GetExpenseReport(ctx context.Context, dateFrom, dateTo string) (*domain.ExpenseReport, error) {
	report := &domain.ExpenseReport{}

	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM expenses
		WHERE expense_date >= ? AND expense_date <= ? AND status = 'approved'`, dateFrom, dateTo).Scan(&report.TotalExpenses)

	// By category
	rows, err := r.db.QueryContext(ctx, `SELECT ec.name, COALESCE(SUM(e.amount),0) as total
		FROM expenses e JOIN expense_categories ec ON e.category_id = ec.id
		WHERE e.expense_date >= ? AND e.expense_date <= ? AND e.status = 'approved'
		GROUP BY ec.name ORDER BY total DESC`, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.ByCategory = append(report.ByCategory, nv)
		}
	}

	// Monthly trend (last 6 months)
	rows, err = r.db.QueryContext(ctx, `SELECT strftime('%Y-%m', expense_date) as period, COALESCE(SUM(amount),0)
		FROM expenses WHERE expense_date >= date('now', '-6 months') AND status = 'approved'
		GROUP BY period ORDER BY period`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp domain.TrendPoint
			rows.Scan(&tp.Period, &tp.Value)
			report.MonthlyTrend = append(report.MonthlyTrend, tp)
		}
	}

	// Revenue vs expense (monthly, last 6 months)
	rows, err = r.db.QueryContext(ctx, `SELECT m.period,
		COALESCE((SELECT SUM(total_amount) FROM invoices WHERE strftime('%Y-%m', created_at) = m.period AND payment_status = 'paid'), 0),
		COALESCE((SELECT SUM(amount) FROM expenses WHERE strftime('%Y-%m', expense_date) = m.period AND status = 'approved'), 0)
		FROM (SELECT DISTINCT strftime('%Y-%m', expense_date) as period FROM expenses WHERE expense_date >= date('now', '-6 months')
		UNION SELECT DISTINCT strftime('%Y-%m', created_at) FROM invoices WHERE created_at >= date('now', '-6 months')) m
		ORDER BY m.period`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dp domain.DualTrendPoint
			rows.Scan(&dp.Period, &dp.Value1, &dp.Value2)
			report.RevenueVsExpense = append(report.RevenueVsExpense, dp)
		}
	}

	return report, nil
}

func (r *AnalyticsRepository) GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	report := &domain.InventoryReport{}

	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(current_stock * purchase_price),0) FROM products WHERE status = 'active'`).Scan(&report.TotalValue)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE status = 'active' AND current_stock < minimum_stock`).Scan(&report.LowStockCount)

	// Fast moving (most consumption in last 30 days)
	rows, err := r.db.QueryContext(ctx, `SELECT p.name, COALESCE(SUM(st.quantity),0) as qty
		FROM stock_transactions st JOIN products p ON st.product_id = p.id
		WHERE st.transaction_type IN ('consumption','sale') AND st.transaction_date >= date('now', '-30 days')
		GROUP BY p.name ORDER BY qty DESC LIMIT 10`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.FastMoving = append(report.FastMoving, nv)
		}
	}

	// Slow moving (least consumption in last 30 days)
	rows, err = r.db.QueryContext(ctx, `SELECT p.name, COALESCE(SUM(st.quantity),0) as qty
		FROM products p LEFT JOIN stock_transactions st ON st.product_id = p.id
		AND st.transaction_type IN ('consumption','sale') AND st.transaction_date >= date('now', '-30 days')
		WHERE p.status = 'active'
		GROUP BY p.name ORDER BY qty ASC LIMIT 10`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nv domain.NameValuePair
			rows.Scan(&nv.Name, &nv.Value)
			report.SlowMoving = append(report.SlowMoving, nv)
		}
	}

	// Purchase trend (monthly)
	rows, err = r.db.QueryContext(ctx, `SELECT strftime('%Y-%m', purchase_date) as period, COALESCE(SUM(total_amount),0)
		FROM purchase_entries WHERE purchase_date >= date('now', '-6 months')
		GROUP BY period ORDER BY period`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp domain.TrendPoint
			rows.Scan(&tp.Period, &tp.Value)
			report.PurchaseTrend = append(report.PurchaseTrend, tp)
		}
	}

	// Consumption trend (monthly)
	rows, err = r.db.QueryContext(ctx, `SELECT strftime('%Y-%m', transaction_date) as period, COALESCE(SUM(quantity),0)
		FROM stock_transactions WHERE transaction_type IN ('consumption','sale') AND transaction_date >= date('now', '-6 months')
		GROUP BY period ORDER BY period`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp domain.TrendPoint
			rows.Scan(&tp.Period, &tp.Value)
			report.ConsumptionTrend = append(report.ConsumptionTrend, tp)
		}
	}

	return report, nil
}

func (r *AnalyticsRepository) GetProfitLossReport(ctx context.Context, dateFrom, dateTo string, groupBy string) (*domain.ProfitLossReport, error) {
	report := &domain.ProfitLossReport{}

	// Totals
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_amount),0) FROM invoices
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? AND payment_status = 'paid'`, dateFrom, dateTo).Scan(&report.Revenue)
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM expenses
		WHERE expense_date >= ? AND expense_date <= ? AND status = 'approved'`, dateFrom, dateTo).Scan(&report.Expenses)
	r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(sr.net_salary),0) FROM salary_records sr
		JOIN salary_cycles sc ON sr.cycle_id = sc.id
		WHERE sc.start_date >= ? AND sc.end_date <= ?`, dateFrom, dateTo).Scan(&report.SalaryCost)
	report.NetProfit = report.Revenue - report.Expenses - report.SalaryCost

	// Trend
	var dateExpr string
	switch groupBy {
	case "month":
		dateExpr = "strftime('%Y-%m', %s)"
	case "year":
		dateExpr = "strftime('%Y', %s)"
	default:
		dateExpr = "DATE(%s)"
	}

	revExpr := fmt.Sprintf(dateExpr, "created_at")
	expExpr := fmt.Sprintf(dateExpr, "expense_date")

	// Build periods from both revenue and expenses
	query := fmt.Sprintf(`SELECT period, 
		COALESCE((SELECT SUM(total_amount) FROM invoices WHERE %s = period AND payment_status = 'paid' AND DATE(created_at) >= ? AND DATE(created_at) <= ?), 0) as revenue,
		COALESCE((SELECT SUM(amount) FROM expenses WHERE %s = period AND status = 'approved' AND expense_date >= ? AND expense_date <= ?), 0) as expenses,
		0 as salary
		FROM (
			SELECT DISTINCT %s as period FROM invoices WHERE DATE(created_at) >= ? AND DATE(created_at) <= ?
			UNION 
			SELECT DISTINCT %s as period FROM expenses WHERE expense_date >= ? AND expense_date <= ?
		) periods ORDER BY period`, revExpr, expExpr, revExpr, expExpr)

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo, dateFrom, dateTo, dateFrom, dateTo, dateFrom, dateTo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pt domain.PLTrendPoint
			rows.Scan(&pt.Period, &pt.Revenue, &pt.Expenses, &pt.Salary)
			pt.Profit = pt.Revenue - pt.Expenses - pt.Salary
			report.Trend = append(report.Trend, pt)
		}
	}

	return report, nil
}
