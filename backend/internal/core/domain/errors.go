package domain

import "errors"

// Sentinel domain errors for Staff.
var (
	ErrStaffNameRequired       = errors.New("staff full name is required")
	ErrStaffPhoneRequired      = errors.New("staff phone is required")
	ErrStaffPhoneInvalid       = errors.New("staff phone must be a valid 10-digit Indian mobile number")
	ErrStaffPhoneDuplicate     = errors.New("staff phone number already exists")
	ErrStaffInvalidDesignation = errors.New("staff designation must be stylist, assistant, receptionist, or manager")
	ErrStaffInvalidGender      = errors.New("staff gender must be male, female, or other")
	ErrStaffInvalidSalary      = errors.New("staff base salary cannot be negative")
	ErrStaffInvalidCommission  = errors.New("staff commission percentage must be between 0 and 100")
)

// Sentinel domain errors for Service.
var (
	ErrServiceNameRequired           = errors.New("service name is required")
	ErrServiceInvalidCategory        = errors.New("service category is invalid")
	ErrServiceInvalidDuration        = errors.New("service duration must be greater than 0")
	ErrServiceInvalidPrice           = errors.New("service price must be greater than 0")
	ErrServiceInvalidCostPrice       = errors.New("service cost price cannot be negative")
	ErrServiceInvalidCommissionType  = errors.New("commission type must be fixed or percentage")
	ErrServiceInvalidCommissionValue = errors.New("commission value is invalid")
	ErrServiceNameDuplicate          = errors.New("service name already exists")
)

// Sentinel domain errors for Customer.
var (
	ErrCustomerNameRequired   = errors.New("customer full name is required")
	ErrCustomerPhoneRequired  = errors.New("customer phone is required")
	ErrCustomerPhoneInvalid   = errors.New("customer phone must be a valid 10-digit Indian mobile number")
	ErrCustomerPhoneDuplicate = errors.New("customer phone number already exists")
	ErrCustomerInvalidGender  = errors.New("customer gender must be male, female, or other")
)

// Sentinel domain errors for Invoice.
var (
	ErrInvoiceCustomerRequired = errors.New("invoice must have a customer")
	ErrInvoiceStaffRequired    = errors.New("invoice must have a staff member")
	ErrInvoiceItemsRequired    = errors.New("invoice must have at least one service")
	ErrInvoiceInvalidDiscount  = errors.New("invoice discount cannot be negative")
	ErrInvoiceInvalidTax       = errors.New("invoice tax cannot be negative")
	ErrInvoiceAlreadyPaid      = errors.New("invoice is already fully paid")
	ErrPaymentInvalidAmount    = errors.New("payment amount must be greater than 0")
	ErrPaymentInvalidMethod    = errors.New("payment method must be cash, upi, card, or bank_transfer")
	ErrPaymentExceedsBalance   = errors.New("payment amount exceeds remaining balance")
)

// Sentinel domain errors for Commission.
var (
	ErrCommissionRuleNameRequired  = errors.New("commission rule name is required")
	ErrCommissionInvalidRuleType   = errors.New("commission rule type must be revenue_based, service_based, or fixed")
	ErrCommissionInvalidTargetType = errors.New("commission target type must be global, staff, or service")
	ErrCommissionInvalidCalcType   = errors.New("commission calculation type must be percentage, fixed_amount, or tiered")
	ErrCommissionInvalidCalcValue  = errors.New("commission calculation value cannot be negative")
)

// Sentinel domain errors for Salary & Advances.
var (
	ErrAdvanceInvalidAmount = errors.New("advance amount must be greater than 0")
	ErrAdvanceDateRequired  = errors.New("advance date is required")
	ErrSalaryCycleExists    = errors.New("salary cycle already exists for this month")
	ErrSalaryCycleNotFound  = errors.New("salary cycle not found")
	ErrSalaryAlreadyPaid    = errors.New("salary is already paid")
	ErrSalaryInvalidMonth   = errors.New("month must be between 1 and 12")
	ErrSalaryInvalidYear    = errors.New("year must be 2020 or later")
)

// Sentinel domain errors for Expense.
var (
	ErrExpenseCategoryNameRequired = errors.New("expense category name is required")
	ErrExpenseCategoryRequired     = errors.New("expense must have a category")
	ErrExpenseInvalidAmount        = errors.New("expense amount must be greater than 0")
	ErrExpenseDateRequired         = errors.New("expense date is required")
	ErrExpenseInvalidPaymentMethod = errors.New("expense payment method must be cash, upi, bank_transfer, card, or cheque")
	ErrExpenseNotFound             = errors.New("expense not found")
	ErrExpenseCategoryNotFound     = errors.New("expense category not found")
)

// Sentinel domain errors for Product & Inventory.
var (
	ErrProductNameRequired     = errors.New("product name is required")
	ErrProductInvalidCategory  = errors.New("product category is invalid")
	ErrProductInvalidPrice     = errors.New("product price cannot be negative")
	ErrProductNotFound         = errors.New("product not found")
	ErrProductRequired         = errors.New("product ID is required")
	ErrInvalidTransactionType  = errors.New("invalid stock transaction type")
	ErrInvalidQuantity         = errors.New("quantity cannot be zero")
	ErrTransactionDateRequired = errors.New("transaction date is required")
	ErrInsufficientStock       = errors.New("insufficient stock for this operation")
	ErrPurchaseVendorRequired  = errors.New("vendor name is required for purchase")
	ErrPurchaseDateRequired    = errors.New("purchase date is required")
	ErrPurchaseItemsRequired   = errors.New("purchase must have at least one item")
)
