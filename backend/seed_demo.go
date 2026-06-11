package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// SeedDemoData populates the database with one month of realistic demo data.
// It is idempotent — if data already exists, it skips seeding.
func SeedDemoData(ctx context.Context, c *Container, log *slog.Logger) error {
	// Check if data already exists (idempotency guard)
	existing, err := c.staffUC.List(ctx, usecase.ListStaffInput{Page: 1, PerPage: 1})
	if err == nil && existing.Total > 0 {
		log.Info("demo data already exists, skipping seed", "existing_staff", existing.Total)
		return nil
	}

	log.Info("seeding demo data for one month...")

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, time.Local)

	// 1. Create Staff
	staffInputs := []usecase.CreateStaffInput{
		{FullName: "Priya Sharma", Phone: "9876543210", Email: "priya@salonflow.com", Gender: "female", Designation: "stylist", JoiningDate: "2025-01-15", BaseSalary: 25000, CommissionPercentage: 15},
		{FullName: "Rahul Verma", Phone: "9876543211", Email: "rahul@salonflow.com", Gender: "male", Designation: "stylist", JoiningDate: "2025-03-01", BaseSalary: 22000, CommissionPercentage: 12},
		{FullName: "Anita Desai", Phone: "9876543212", Email: "anita@salonflow.com", Gender: "female", Designation: "manager", JoiningDate: "2024-06-10", BaseSalary: 35000, CommissionPercentage: 10},
		{FullName: "Vikram Singh", Phone: "9876543213", Email: "vikram@salonflow.com", Gender: "male", Designation: "assistant", JoiningDate: "2025-06-01", BaseSalary: 15000, CommissionPercentage: 8},
		{FullName: "Meera Patel", Phone: "9876543214", Email: "meera@salonflow.com", Gender: "female", Designation: "stylist", JoiningDate: "2025-02-20", BaseSalary: 28000, CommissionPercentage: 18},
	}

	var staffIDs []string
	for _, input := range staffInputs {
		staff, err := c.staffUC.Create(ctx, input)
		if err != nil {
			return fmt.Errorf("create staff %s: %w", input.FullName, err)
		}
		staffIDs = append(staffIDs, staff.ID.String())
		log.Info("created staff", "name", staff.FullName)
	}

	// 2. Create Services
	serviceInputs := []usecase.CreateServiceInput{
		{Name: "Haircut - Women", Category: "hair", Description: "Professional women's haircut with styling", DurationMinutes: 45, Price: 800, CostPrice: 100, CommissionType: "percentage", CommissionValue: 15},
		{Name: "Haircut - Men", Category: "hair", Description: "Men's haircut and grooming", DurationMinutes: 30, Price: 400, CostPrice: 50, CommissionType: "percentage", CommissionValue: 12},
		{Name: "Hair Coloring", Category: "coloring", Description: "Full head coloring with premium products", DurationMinutes: 120, Price: 3500, CostPrice: 800, CommissionType: "percentage", CommissionValue: 10},
		{Name: "Hair Spa", Category: "spa", Description: "Deep conditioning hair spa treatment", DurationMinutes: 60, Price: 1500, CostPrice: 300, CommissionType: "fixed", CommissionValue: 200},
		{Name: "Classic Facial", Category: "facial", Description: "Deep cleansing facial for radiant skin", DurationMinutes: 60, Price: 1200, CostPrice: 250, CommissionType: "percentage", CommissionValue: 12},
		{Name: "Gold Facial", Category: "facial", Description: "Premium gold facial for anti-aging", DurationMinutes: 90, Price: 2500, CostPrice: 600, CommissionType: "percentage", CommissionValue: 15},
		{Name: "Full Body Massage", Category: "massage", Description: "Relaxing full body massage", DurationMinutes: 60, Price: 2000, CostPrice: 200, CommissionType: "fixed", CommissionValue: 300},
		{Name: "Manicure", Category: "skin", Description: "Professional manicure with nail art", DurationMinutes: 45, Price: 600, CostPrice: 100, CommissionType: "fixed", CommissionValue: 100},
		{Name: "Pedicure", Category: "skin", Description: "Relaxing pedicure treatment", DurationMinutes: 45, Price: 700, CostPrice: 120, CommissionType: "fixed", CommissionValue: 100},
		{Name: "Keratin Treatment", Category: "treatment", Description: "Smoothing keratin treatment", DurationMinutes: 180, Price: 5000, CostPrice: 1500, CommissionType: "percentage", CommissionValue: 10},
		{Name: "Bridal Makeup", Category: "other", Description: "Complete bridal makeup package", DurationMinutes: 120, Price: 8000, CostPrice: 1500, CommissionType: "percentage", CommissionValue: 20},
		{Name: "Threading - Eyebrows", Category: "other", Description: "Eyebrow threading and shaping", DurationMinutes: 15, Price: 100, CostPrice: 10, CommissionType: "fixed", CommissionValue: 20},
	}

	type svcRef struct {
		ID       string
		Name     string
		Price    float64
		Duration int
	}
	var services []svcRef
	for _, input := range serviceInputs {
		svc, err := c.serviceUC.Create(ctx, input)
		if err != nil {
			return fmt.Errorf("create service %s: %w", input.Name, err)
		}
		services = append(services, svcRef{ID: svc.ID.String(), Name: svc.Name, Price: svc.Price, Duration: svc.DurationMinutes})
		log.Info("created service", "name", svc.Name)
	}

	// 3. Create Customers
	customerInputs := []usecase.CreateCustomerInput{
		{FullName: "Aisha Khan", Phone: "9123456780", Email: "aisha@email.com", Gender: "female", DateOfBirth: "1992-03-15", Address: "HSR Layout, Bangalore"},
		{FullName: "Rohan Gupta", Phone: "9123456781", Email: "rohan@email.com", Gender: "male", DateOfBirth: "1988-07-22", Address: "Koramangala, Bangalore"},
		{FullName: "Sneha Reddy", Phone: "9123456782", Email: "sneha@email.com", Gender: "female", DateOfBirth: "1995-11-08", Address: "Indiranagar, Bangalore"},
		{FullName: "Amit Joshi", Phone: "9123456783", Email: "amit@email.com", Gender: "male", DateOfBirth: "1985-01-30", Address: "Whitefield, Bangalore"},
		{FullName: "Kavya Nair", Phone: "9123456784", Email: "kavya@email.com", Gender: "female", DateOfBirth: "1998-06-12", Address: "JP Nagar, Bangalore"},
		{FullName: "Deepak Menon", Phone: "9123456785", Email: "deepak@email.com", Gender: "male", DateOfBirth: "1990-09-25", Address: "Jayanagar, Bangalore"},
		{FullName: "Pooja Iyer", Phone: "9123456786", Email: "pooja@email.com", Gender: "female", DateOfBirth: "1993-04-18", Address: "BTM Layout, Bangalore"},
		{FullName: "Suresh Kumar", Phone: "9123456787", Email: "suresh@email.com", Gender: "male", DateOfBirth: "1982-12-05", Address: "MG Road, Bangalore"},
		{FullName: "Nisha Agarwal", Phone: "9123456788", Email: "nisha@email.com", Gender: "female", DateOfBirth: "1997-02-28", Address: "Electronic City, Bangalore"},
		{FullName: "Rajesh Pillai", Phone: "9123456789", Email: "rajesh@email.com", Gender: "male", DateOfBirth: "1987-08-14", Address: "Marathahalli, Bangalore"},
		{FullName: "Divya Sharma", Phone: "9123456790", Email: "divya@email.com", Gender: "female", DateOfBirth: "1994-05-20", Address: "Hebbal, Bangalore"},
		{FullName: "Kiran Rao", Phone: "9123456791", Email: "kiran@email.com", Gender: "male", DateOfBirth: "1991-10-03", Address: "Yelahanka, Bangalore"},
		{FullName: "Lakshmi Devi", Phone: "9123456792", Email: "lakshmi@email.com", Gender: "female", DateOfBirth: "1989-07-17", Address: "Malleshwaram, Bangalore"},
		{FullName: "Arjun Shetty", Phone: "9123456793", Email: "arjun@email.com", Gender: "male", DateOfBirth: "1996-01-09", Address: "Basavanagudi, Bangalore"},
		{FullName: "Fatima Begum", Phone: "9123456794", Email: "fatima@email.com", Gender: "female", DateOfBirth: "1986-11-22", Address: "Richmond Town, Bangalore"},
	}

	var customerIDs []string
	for _, input := range customerInputs {
		cust, err := c.customerUC.Create(ctx, input)
		if err != nil {
			return fmt.Errorf("create customer %s: %w", input.FullName, err)
		}
		customerIDs = append(customerIDs, cust.ID.String())
		log.Info("created customer", "name", cust.FullName)
	}

	// 4. Create Products (Inventory)
	productInputs := []usecase.CreateProductInput{
		{Name: "L'Oreal Hair Color - Black", Category: "coloring", Brand: "L'Oreal", Unit: "pcs", SKU: "LOR-HC-BLK", PurchasePrice: 450, SellingPrice: 700, MinimumStock: 10, MaximumStock: 50},
		{Name: "L'Oreal Hair Color - Brown", Category: "coloring", Brand: "L'Oreal", Unit: "pcs", SKU: "LOR-HC-BRN", PurchasePrice: 450, SellingPrice: 700, MinimumStock: 10, MaximumStock: 50},
		{Name: "Matrix Shampoo 500ml", Category: "hair_care", Brand: "Matrix", Unit: "pcs", SKU: "MTX-SH-500", PurchasePrice: 380, SellingPrice: 550, MinimumStock: 5, MaximumStock: 30},
		{Name: "Matrix Conditioner 500ml", Category: "hair_care", Brand: "Matrix", Unit: "pcs", SKU: "MTX-CN-500", PurchasePrice: 400, SellingPrice: 600, MinimumStock: 5, MaximumStock: 30},
		{Name: "Wella Hair Spa Cream 400g", Category: "spa", Brand: "Wella", Unit: "pcs", SKU: "WEL-HSP-400", PurchasePrice: 550, SellingPrice: 850, MinimumStock: 5, MaximumStock: 25},
		{Name: "VLCC Gold Facial Kit", Category: "facial", Brand: "VLCC", Unit: "pcs", SKU: "VLCC-GF-KIT", PurchasePrice: 800, SellingPrice: 1200, MinimumStock: 3, MaximumStock: 20},
		{Name: "OPI Nail Polish Set", Category: "retail", Brand: "OPI", Unit: "set", SKU: "OPI-NP-SET", PurchasePrice: 1200, SellingPrice: 1800, MinimumStock: 5, MaximumStock: 20},
		{Name: "Keratin Treatment Kit", Category: "treatment", Brand: "GK Hair", Unit: "pcs", SKU: "GKH-KT-KIT", PurchasePrice: 2500, SellingPrice: 3800, MinimumStock: 2, MaximumStock: 10},
		{Name: "Massage Oil - Lavender 1L", Category: "spa", Brand: "Kama Ayurveda", Unit: "bottle", SKU: "KAM-MO-LAV", PurchasePrice: 600, SellingPrice: 950, MinimumStock: 3, MaximumStock: 15},
		{Name: "Disposable Towels (Pack 100)", Category: "equipment", Brand: "Generic", Unit: "pack", SKU: "GEN-DT-100", PurchasePrice: 350, SellingPrice: 500, MinimumStock: 5, MaximumStock: 20},
	}

	var productIDs []string
	for _, input := range productInputs {
		prod, err := c.productUC.CreateProduct(ctx, input)
		if err != nil {
			return fmt.Errorf("create product %s: %w", input.Name, err)
		}
		productIDs = append(productIDs, prod.ID.String())
		log.Info("created product", "name", prod.Name)
	}

	// Add initial stock to all products
	for _, pid := range productIDs {
		_, err := c.productUC.AdjustStock(ctx, usecase.StockAdjustInput{
			ProductID:       pid,
			TransactionType: "purchase",
			Quantity:        20,
			Notes:           "Initial stock",
		})
		if err != nil {
			log.Warn("stock adjust failed", "product", pid, "error", err)
		}
	}

	// 5. Create Expense Categories (or use existing ones)
	expenseCategories := []struct{ Name, Desc string }{
		{"Rent", "Monthly salon rent"},
		{"Electricity", "Electricity bills"},
		{"Products & Supplies", "Salon product purchases"},
		{"Staff Welfare", "Tea, snacks, uniforms"},
		{"Marketing", "Ads, pamphlets, social media"},
		{"Maintenance", "Equipment repair, plumbing"},
		{"Miscellaneous", "Other expenses"},
	}

	var expCatIDs []string

	// Try to use existing categories first
	existingCats, err := c.expenseUC.ListCategories(ctx, false)
	if err == nil && len(existingCats) > 0 {
		for _, cat := range existingCats {
			expCatIDs = append(expCatIDs, cat.ID.String())
		}
		log.Info("using existing expense categories", "count", len(expCatIDs))
	} else {
		for _, cat := range expenseCategories {
			ec, err := c.expenseUC.CreateCategory(ctx, cat.Name, cat.Desc)
			if err != nil {
				log.Warn("create expense category failed (may already exist)", "name", cat.Name, "error", err)
				continue
			}
			expCatIDs = append(expCatIDs, ec.ID.String())
			log.Info("created expense category", "name", cat.Name)
		}
	}

	if len(expCatIDs) == 0 {
		log.Warn("no expense categories available, skipping expenses")
	}

	// 6. Create Invoices (spread across the month) - this also creates performance/commission data
	rng := rand.New(rand.NewSource(42))
	paymentMethods := []string{"cash", "upi", "card", "bank_transfer"}

	for day := 0; day < 30; day++ {
		date := monthStart.AddDate(0, 0, day)
		if date.After(now) {
			break
		}

		// 3-8 invoices per day
		numInvoices := 3 + rng.Intn(6)
		for i := 0; i < numInvoices; i++ {
			customerIdx := rng.Intn(len(customerIDs))
			staffIdx := rng.Intn(len(staffIDs))

			// 1-3 services per invoice
			numItems := 1 + rng.Intn(3)
			var items []usecase.CreateInvoiceItemInput
			usedServices := make(map[int]bool)
			for j := 0; j < numItems; j++ {
				svcIdx := rng.Intn(len(services))
				if usedServices[svcIdx] {
					continue
				}
				usedServices[svcIdx] = true
				items = append(items, usecase.CreateInvoiceItemInput{
					ServiceID: services[svcIdx].ID,
					Quantity:  1,
					Discount:  0,
				})
			}

			if len(items) == 0 {
				items = append(items, usecase.CreateInvoiceItemInput{
					ServiceID: services[0].ID,
					Quantity:  1,
					Discount:  0,
				})
			}

			discount := 0.0
			if rng.Float64() < 0.2 { // 20% chance of discount
				discount = float64(rng.Intn(3)+1) * 100 // 100-300 discount
			}

			input := usecase.CreateInvoiceInput{
				CustomerID:    customerIDs[customerIdx],
				StaffID:       staffIDs[staffIdx],
				Items:         items,
				Discount:      discount,
				Tax:           0,
				PaymentMethod: paymentMethods[rng.Intn(len(paymentMethods))],
				Notes:         fmt.Sprintf("Demo invoice - Day %d", day+1),
			}

			_, err := c.invoiceUC.Create(ctx, input)
			if err != nil {
				log.Warn("create invoice failed", "day", day, "error", err)
				continue
			}
		}

		// Log progress every 5 days
		if (day+1)%5 == 0 {
			log.Info("seeded invoices", "days_done", day+1)
		}
	}

	// 7. Create Appointments (past month + some upcoming)
	timeSlots := []string{"09:00", "10:00", "11:00", "12:00", "14:00", "15:00", "16:00", "17:00"}
	apptStatuses := []string{"completed", "completed", "completed", "completed", "cancelled", "no_show"} // weighted toward completed

	for day := 0; day < 30; day++ {
		date := monthStart.AddDate(0, 0, day)
		dateStr := date.Format("2006-01-02")

		// 4-10 appointments per day
		numAppts := 4 + rng.Intn(7)
		for i := 0; i < numAppts && i < len(timeSlots); i++ {
			customerIdx := rng.Intn(len(customerIDs))
			staffIdx := rng.Intn(len(staffIDs))
			svcIdx := rng.Intn(len(services))

			startTime := timeSlots[i]
			endHour := 9 + i + 1
			endTime := fmt.Sprintf("%02d:00", endHour)

			isWalkin := rng.Float64() < 0.3 // 30% walk-ins

			appt := domain.NewAppointment(
				customerIDs[customerIdx],
				staffIDs[staffIdx],
				dateStr,
				startTime,
				endTime,
				isWalkin,
			)

			// Set status based on whether in past or future
			if date.Before(now) {
				appt.Status = apptStatuses[rng.Intn(len(apptStatuses))]
			} else {
				appt.Status = "booked"
			}

			apptServices := []domain.AppointmentService{
				*domain.NewAppointmentService(
					appt.ID,
					services[svcIdx].ID,
					services[svcIdx].Name,
					services[svcIdx].Duration,
					services[svcIdx].Price,
				),
			}

			err := c.apptUC.Create(ctx, appt, apptServices)
			if err != nil {
				log.Warn("create appointment failed", "day", day, "error", err)
				continue
			}
		}
	}
	log.Info("seeded appointments")

	// 8. Create Expenses (throughout the month)
	expenseVendors := []string{"Property Owner", "BESCOM", "Beauty Wholesale", "Local Supplier", "Google Ads", "Plumber Ji", "General Store"}
	expenseAmounts := []float64{50000, 8000, 15000, 3000, 5000, 2000, 1500}

	if len(expCatIDs) > 0 {
		for day := 0; day < 30; day += 3 { // Every 3 days
			date := monthStart.AddDate(0, 0, day)
			if date.After(now) {
				break
			}
			dateStr := date.Format("2006-01-02")

			catIdx := rng.Intn(len(expCatIDs))
			vendorIdx := rng.Intn(len(expenseVendors))
			amount := expenseAmounts[vendorIdx] * (0.8 + rng.Float64()*0.4) // vary ±20%

			input := usecase.CreateExpenseInput{
				CategoryID:       expCatIDs[catIdx],
				Amount:           float64(int(amount)),
				ExpenseDate:      dateStr,
				PaymentMethod:    paymentMethods[rng.Intn(len(paymentMethods))],
				VendorName:       expenseVendors[vendorIdx],
				InvoiceReference: fmt.Sprintf("EXP-%s-%02d", date.Format("0601"), day),
				Description:      fmt.Sprintf("Monthly %s expense", expenseVendors[vendorIdx]),
			}

			_, err := c.expenseUC.CreateExpense(ctx, input)
			if err != nil {
				log.Warn("create expense failed", "day", day, "error", err)
			}
		}
	}
	log.Info("seeded expenses")

	// 9. Create Membership Plans
	plan1 := domain.NewMembershipPlan("Silver Package", "package", 5000, 90, 10)
	plan1.Description = "10 sessions over 3 months with 10% discount"
	plan1.DiscountPercentage = 10

	plan1Services := []domain.PackageService{
		*domain.NewPackageService(plan1.ID, services[0].ID, services[0].Name, 5),
		*domain.NewPackageService(plan1.ID, services[4].ID, services[4].Name, 5),
	}
	if err := c.membershipUC.CreatePlan(ctx, plan1, plan1Services); err != nil {
		log.Warn("create plan1 failed", "error", err)
	}

	plan2 := domain.NewMembershipPlan("Gold Membership", "membership", 12000, 180, 20)
	plan2.Description = "6-month premium membership with 15% off all services"
	plan2.DiscountPercentage = 15
	plan2.PriorityBooking = true

	plan2Services := []domain.PackageService{
		*domain.NewPackageService(plan2.ID, services[0].ID, services[0].Name, 8),
		*domain.NewPackageService(plan2.ID, services[4].ID, services[4].Name, 6),
		*domain.NewPackageService(plan2.ID, services[6].ID, services[6].Name, 6),
	}
	if err := c.membershipUC.CreatePlan(ctx, plan2, plan2Services); err != nil {
		log.Warn("create plan2 failed", "error", err)
	}

	plan3 := domain.NewMembershipPlan("Bridal Package", "package", 25000, 30, 5)
	plan3.Description = "Complete bridal preparation package"
	plan3.DiscountPercentage = 20

	plan3Services := []domain.PackageService{
		*domain.NewPackageService(plan3.ID, services[10].ID, services[10].Name, 1),
		*domain.NewPackageService(plan3.ID, services[5].ID, services[5].Name, 2),
		*domain.NewPackageService(plan3.ID, services[3].ID, services[3].Name, 2),
	}
	if err := c.membershipUC.CreatePlan(ctx, plan3, plan3Services); err != nil {
		log.Warn("create plan3 failed", "error", err)
	}

	// Sell some memberships
	_, err = c.membershipUC.SellPlan(ctx, customerIDs[0], plan1.ID, 5000)
	if err != nil {
		log.Warn("sell plan1 failed", "error", err)
	}
	_, err = c.membershipUC.SellPlan(ctx, customerIDs[2], plan2.ID, 12000)
	if err != nil {
		log.Warn("sell plan2 failed", "error", err)
	}
	_, err = c.membershipUC.SellPlan(ctx, customerIDs[4], plan3.ID, 25000)
	if err != nil {
		log.Warn("sell plan3 failed", "error", err)
	}

	log.Info("seeded memberships")

	// 10. Stock consumption throughout the month
	for day := 0; day < 30; day += 2 {
		date := monthStart.AddDate(0, 0, day)
		if date.After(now) {
			break
		}

		// Consume 1-3 random products
		numProducts := 1 + rng.Intn(3)
		for i := 0; i < numProducts; i++ {
			prodIdx := rng.Intn(len(productIDs))
			qty := float64(1 + rng.Intn(3))
			_, err := c.productUC.AdjustStock(ctx, usecase.StockAdjustInput{
				ProductID:       productIDs[prodIdx],
				TransactionType: "consumption",
				Quantity:        qty,
				Notes:           fmt.Sprintf("Daily usage - Day %d", day+1),
			})
			if err != nil {
				// ignore - may not have enough stock
				continue
			}
		}
	}
	log.Info("seeded stock transactions")

	log.Info("demo data seeding complete!",
		"staff", len(staffIDs),
		"services", len(services),
		"customers", len(customerIDs),
		"products", len(productIDs),
		"expense_categories", len(expCatIDs),
	)

	return nil
}
