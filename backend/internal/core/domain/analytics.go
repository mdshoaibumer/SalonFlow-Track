package domain

// ---- Analytics Domain Types ----

// DashboardStats holds the executive dashboard KPIs.
type DashboardStats struct {
	TodayRevenue       float64 `json:"today_revenue"`
	TodayCustomers     int     `json:"today_customers"`
	TodayInvoices      int     `json:"today_invoices"`
	MonthlyRevenue     float64 `json:"monthly_revenue"`
	MonthlyExpenses    float64 `json:"monthly_expenses"`
	MonthlyProfit      float64 `json:"monthly_profit"`
	InventoryValue     float64 `json:"inventory_value"`
	OutstandingSalary  float64 `json:"outstanding_salary"`
	OutstandingAdvance float64 `json:"outstanding_advances"`
	LowStockCount      int     `json:"low_stock_count"`
}

// KPIMetrics holds business KPI calculations.
type KPIMetrics struct {
	RevenueGrowthPct  float64 `json:"revenue_growth_pct"`
	CustomerGrowthPct float64 `json:"customer_growth_pct"`
	ProfitMarginPct   float64 `json:"profit_margin_pct"`
	AverageBillValue  float64 `json:"average_bill_value"`
	RepeatCustomerPct float64 `json:"repeat_customer_pct"`
	StaffProductivity float64 `json:"staff_productivity_pct"`
}

// RevenueReport holds revenue analytics data.
type RevenueReport struct {
	Trend        []RevenueTrendPoint `json:"trend"`
	ByService    []NameValuePair     `json:"by_service"`
	ByStaff      []NameValuePair     `json:"by_staff"`
	ByCustomer   []NameValuePair     `json:"by_customer"`
	TotalRevenue float64             `json:"total_revenue"`
	InvoiceCount int                 `json:"invoice_count"`
}

// CustomerReport holds customer analytics data.
type CustomerReport struct {
	TotalCustomers  int             `json:"total_customers"`
	NewCustomers    int             `json:"new_customers"`
	RepeatCustomers int             `json:"repeat_customers"`
	BirthdayToday   int             `json:"birthday_today"`
	InactiveCount   int             `json:"inactive_count"`
	TopCustomers    []NameValuePair `json:"top_customers"`
	GrowthTrend     []TrendPoint    `json:"growth_trend"`
}

// StaffReport holds staff analytics data.
type StaffReport struct {
	TopPerformers    []NameValuePair `json:"top_performers"`
	RevenueByStaff   []NameValuePair `json:"revenue_by_staff"`
	CustomersByStaff []NameValuePair `json:"customers_by_staff"`
	CommissionEarned []NameValuePair `json:"commission_earned"`
	SalaryCost       float64         `json:"salary_cost"`
}

// ServiceReport holds service analytics data.
type ServiceReport struct {
	TopServices      []NameValuePair `json:"top_services"`
	LeastUsed        []NameValuePair `json:"least_used"`
	RevenueByService []NameValuePair `json:"revenue_by_service"`
	AvgServiceValue  float64         `json:"avg_service_value"`
	TotalBookings    int             `json:"total_bookings"`
}

// ExpenseReport holds expense analytics data.
type ExpenseReport struct {
	TotalExpenses    float64          `json:"total_expenses"`
	ByCategory       []NameValuePair  `json:"by_category"`
	MonthlyTrend     []TrendPoint     `json:"monthly_trend"`
	RevenueVsExpense []DualTrendPoint `json:"revenue_vs_expense"`
}

// InventoryReport holds inventory analytics data.
type InventoryReport struct {
	TotalValue       float64         `json:"total_value"`
	LowStockCount    int             `json:"low_stock_count"`
	FastMoving       []NameValuePair `json:"fast_moving"`
	SlowMoving       []NameValuePair `json:"slow_moving"`
	PurchaseTrend    []TrendPoint    `json:"purchase_trend"`
	ConsumptionTrend []TrendPoint    `json:"consumption_trend"`
}

// ProfitLossReport holds P&L data.
type ProfitLossReport struct {
	Revenue    float64        `json:"revenue"`
	Expenses   float64        `json:"expenses"`
	SalaryCost float64        `json:"salary_cost"`
	NetProfit  float64        `json:"net_profit"`
	Trend      []PLTrendPoint `json:"trend"`
}

// PLTrendPoint represents a P&L data point.
type PLTrendPoint struct {
	Period   string  `json:"period"`
	Revenue  float64 `json:"revenue"`
	Expenses float64 `json:"expenses"`
	Salary   float64 `json:"salary"`
	Profit   float64 `json:"profit"`
}

// NameValuePair is a generic label-value pair.
type NameValuePair struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// TrendPoint represents a time-series data point.
type TrendPoint struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
}

// DualTrendPoint represents two series on same timeline.
type DualTrendPoint struct {
	Period string  `json:"period"`
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
}
