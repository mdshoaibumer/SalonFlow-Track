package main

import (
	"context"
	"log/slog"

	"github.com/salonflow/salonflow-track/internal/adapters/backup"
	"github.com/salonflow/salonflow-track/internal/adapters/cloudbackup"
	"github.com/salonflow/salonflow-track/internal/adapters/gst"
	"github.com/salonflow/salonflow-track/internal/adapters/importer"
	"github.com/salonflow/salonflow-track/internal/adapters/license"
	"github.com/salonflow/salonflow-track/internal/adapters/printer"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/adapters/update"
	"github.com/salonflow/salonflow-track/internal/config"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/internal/database"
	applogger "github.com/salonflow/salonflow-track/internal/logger"
)

// Container is the application's composition root.
// It owns all dependencies and wires them together.
type Container struct {
	cfg *config.Config
	log *slog.Logger
	db  *database.DB

	// Repositories
	staffRepo      ports.StaffRepository
	serviceRepo    ports.ServiceRepository
	customerRepo   ports.CustomerRepository
	invoiceRepo    ports.InvoiceRepository
	perfRepo       ports.PerformanceRepository
	commissionRepo ports.CommissionRepository
	salaryRepo     ports.SalaryRepository
	expenseRepo    ports.ExpenseRepository
	productRepo    ports.ProductRepository
	analyticsRepo  ports.AnalyticsRepository
	backupRepo     ports.BackupRepository
	backupEngine   ports.BackupEngine
	licenseRepo    ports.LicenseRepository
	licenseEngine  ports.LicenseEngine
	updateRepo     ports.UpdateRepository
	updateEngine   ports.UpdateEngine
	importRepo     ports.ImportRepository
	importEngine   ports.ImportEngine
	gstRepo        ports.GSTRepository
	gstEngine      ports.GSTEngine
	printerRepo    ports.PrinterRepository
	printerEngine  ports.PrintEngine
	apptRepo       ports.AppointmentRepository
	whatsappRepo   ports.WhatsAppRepository
	membershipRepo ports.MembershipRepository
	cloudRepo      ports.CloudBackupRepository
	cloudEngine    ports.CloudBackupEngine
	authRepo       ports.AuthRepository

	// Use Cases
	staffUC       *usecase.StaffUseCase
	serviceUC     *usecase.ServiceUseCase
	customerUC    *usecase.CustomerUseCase
	invoiceUC     *usecase.InvoiceUseCase
	perfUC        *usecase.PerformanceUseCase
	commissionUC  *usecase.CommissionUseCase
	salaryUC      *usecase.SalaryUseCase
	expenseUC     *usecase.ExpenseUseCase
	productUC     *usecase.ProductUseCase
	analyticsUC   *usecase.AnalyticsUseCase
	backupUC      *usecase.BackupUseCase
	licenseUC     *usecase.LicenseUseCase
	updateUC      *usecase.UpdateUseCase
	importUC      *usecase.ImportUseCase
	gstUC         *usecase.GSTUseCase
	printerUC     *usecase.PrinterUseCase
	apptUC        *usecase.AppointmentUseCase
	whatsappUC    *usecase.WhatsAppUseCase
	membershipUC  *usecase.MembershipUseCase
	cloudUC       *usecase.CloudBackupUseCase
	authUC        *usecase.AuthUseCase
	auditSvc      *usecase.AuditService
	diagnosticsUC *usecase.DiagnosticsUseCase

	// Wails Binding Services
	StaffSvc       *StaffService
	CustomerSvc    *CustomerService
	ServiceSvc     *ServiceService
	InvoiceSvc     *InvoiceService
	ExpenseSvc     *ExpenseService
	ProductSvc     *ProductService
	PerformanceSvc *PerformanceService
	CommissionSvc  *CommissionService
	SalarySvc      *SalaryService
	AnalyticsSvc   *AnalyticsService
	BackupSvc      *BackupService
	LicenseSvc     *LicenseService
	UpdateSvc      *UpdateService
	ImportSvc      *ImportService
	GSTSvc         *GSTService
	PrinterSvc     *PrinterService
	AppointmentSvc *AppointmentService
	WhatsAppSvc    *WhatsAppService
	MembershipSvc  *MembershipService
	CloudBackupSvc *CloudBackupService
	AuthSvc        *AuthService
	DiagnosticsSvc *DiagnosticsService

	// Multi-logger
	MultiLog *applogger.MultiLogger
}

// NewContainer builds the full dependency graph.
func NewContainer(cfg *config.Config, log *slog.Logger, db *database.DB) *Container {
	c := &Container{
		cfg: cfg,
		log: log,
		db:  db,
	}
	c.MultiLog = applogger.NewMultiLogger(cfg.Log)
	c.initRepositories()
	c.initUseCases()
	c.initBindings()
	return c
}

func (c *Container) initRepositories() {
	c.staffRepo = sqlite.NewStaffRepository(c.db.Conn(), c.log)
	c.serviceRepo = sqlite.NewServiceRepository(c.db.Conn(), c.log)
	c.customerRepo = sqlite.NewCustomerRepository(c.db.Conn(), c.log)
	c.invoiceRepo = sqlite.NewInvoiceRepository(c.db.Conn(), c.log)
	c.perfRepo = sqlite.NewPerformanceRepository(c.db.Conn(), c.log)
	c.commissionRepo = sqlite.NewCommissionRepository(c.db.Conn(), c.log)
	c.salaryRepo = sqlite.NewSalaryRepository(c.db.Conn(), c.log)
	c.expenseRepo = sqlite.NewExpenseRepository(c.db.Conn(), c.log)
	c.productRepo = sqlite.NewProductRepository(c.db.Conn(), c.log)
	c.analyticsRepo = sqlite.NewAnalyticsRepository(c.db.Conn(), c.log)
	c.backupRepo = sqlite.NewBackupRepository(c.db.Conn(), c.log)
	c.backupEngine = backup.NewEngine()
	c.licenseRepo = sqlite.NewLicenseRepository(c.db.Conn(), c.log)
	c.licenseEngine = license.NewEngine()
	c.updateRepo = sqlite.NewUpdateRepository(c.db.Conn(), c.log)
	c.updateEngine = update.NewEngine()
	c.importRepo = sqlite.NewImportRepository(c.db.Conn(), c.log)
	c.importEngine = importer.NewEngine()
	c.gstRepo = sqlite.NewGSTRepository(c.db.Conn(), c.log)
	c.gstEngine = gst.NewEngine()
	c.printerRepo = sqlite.NewPrinterRepository(c.db.Conn(), c.log)
	c.printerEngine = printer.NewEngine()
	c.apptRepo = sqlite.NewAppointmentRepository(c.db.Conn(), c.log)
	c.whatsappRepo = sqlite.NewWhatsAppRepository(c.db.Conn(), c.log)
	c.membershipRepo = sqlite.NewMembershipRepository(c.db.Conn(), c.log)
	c.cloudRepo = sqlite.NewCloudBackupRepository(c.db.Conn(), c.log)
	c.cloudEngine = cloudbackup.NewEngine()
	c.authRepo = sqlite.NewAuthRepository(c.db.Conn(), c.log)
}

func (c *Container) initUseCases() {
	c.staffUC = usecase.NewStaffUseCase(c.staffRepo, c.log)
	c.serviceUC = usecase.NewServiceUseCase(c.serviceRepo, c.log)
	c.customerUC = usecase.NewCustomerUseCase(c.customerRepo, c.log)
	c.perfUC = usecase.NewPerformanceUseCase(c.perfRepo, c.log)
	c.commissionUC = usecase.NewCommissionUseCase(c.commissionRepo, c.log)
	c.invoiceUC = usecase.NewInvoiceUseCase(c.invoiceRepo, c.customerRepo, c.serviceRepo, c.staffRepo, c.perfUC, c.commissionUC, c.log)
	c.salaryUC = usecase.NewSalaryUseCase(c.salaryRepo, c.staffRepo, c.commissionRepo, c.log)
	c.expenseUC = usecase.NewExpenseUseCase(c.expenseRepo, c.invoiceRepo, c.log)
	c.productUC = usecase.NewProductUseCase(c.productRepo, c.log)
	c.analyticsUC = usecase.NewAnalyticsUseCase(c.analyticsRepo, c.log)
	c.backupUC = usecase.NewBackupUseCase(c.backupRepo, c.backupEngine, c.cfg.Database.Path, c.log)
	c.licenseUC = usecase.NewLicenseUseCase(c.licenseRepo, c.licenseEngine, c.log)
	c.updateUC = usecase.NewUpdateUseCase(c.updateRepo, c.updateEngine, c.backupEngine, c.cfg.App.Version, c.cfg.Database.Path, c.log)
	c.importUC = usecase.NewImportUseCase(c.importRepo, c.importEngine, c.log)
	c.gstUC = usecase.NewGSTUseCase(c.gstRepo, c.gstEngine, c.log)
	c.printerUC = usecase.NewPrinterUseCase(c.printerRepo, c.printerEngine, c.log)
	c.apptUC = usecase.NewAppointmentUseCase(c.apptRepo, c.log)
	c.whatsappUC = usecase.NewWhatsAppUseCase(c.whatsappRepo, c.log)
	c.membershipUC = usecase.NewMembershipUseCase(c.membershipRepo, c.log)
	c.cloudUC = usecase.NewCloudBackupUseCase(c.cloudRepo, c.cloudEngine, c.cfg.Database.Path, c.log)
	c.authUC = usecase.NewAuthUseCase(c.authRepo, c.log, c.cfg.App.Version)
	c.auditSvc = usecase.NewAuditService(c.authRepo, c.log, c.cfg.App.Version)
	c.diagnosticsUC = usecase.NewDiagnosticsUseCase(c.db.Conn(), c.log, c.cfg.App.Version, c.cfg.Database.Path, c.MultiLog.LogDir())
}

func (c *Container) initBindings() {
	c.AuthSvc = NewAuthService(c.authUC, c.auditSvc)
	c.DiagnosticsSvc = NewDiagnosticsService(c.diagnosticsUC)
	guard := NewPermissionGuard(c.AuthSvc)
	licGuard := NewLicenseGuard(c.licenseUC)
	c.StaffSvc = NewStaffService(c.staffUC)
	c.StaffSvc.guard = guard
	c.CustomerSvc = NewCustomerService(c.customerUC)
	c.CustomerSvc.guard = guard
	c.CustomerSvc.licGuard = licGuard
	c.ServiceSvc = NewServiceService(c.serviceUC)
	c.ServiceSvc.guard = guard
	c.InvoiceSvc = NewInvoiceService(c.invoiceUC)
	c.InvoiceSvc.guard = guard
	c.InvoiceSvc.licGuard = licGuard
	c.ExpenseSvc = NewExpenseService(c.expenseUC)
	c.ExpenseSvc.guard = guard
	c.ExpenseSvc.licGuard = licGuard
	c.ProductSvc = NewProductService(c.productUC)
	c.ProductSvc.guard = guard
	c.ProductSvc.licGuard = licGuard
	c.PerformanceSvc = NewPerformanceService(c.perfUC)
	c.PerformanceSvc.guard = guard
	c.CommissionSvc = NewCommissionService(c.commissionUC)
	c.CommissionSvc.guard = guard
	c.SalarySvc = NewSalaryService(c.salaryUC)
	c.SalarySvc.guard = guard
	c.SalarySvc.licGuard = licGuard
	c.AnalyticsSvc = NewAnalyticsService(c.analyticsUC)
	c.AnalyticsSvc.guard = guard
	c.BackupSvc = NewBackupService(c.backupUC)
	c.BackupSvc.guard = guard
	c.LicenseSvc = NewLicenseService(c.licenseUC)
	c.LicenseSvc.guard = guard
	c.UpdateSvc = NewUpdateService(c.updateUC)
	c.UpdateSvc.guard = guard
	c.ImportSvc = NewImportService(c.importUC)
	c.ImportSvc.guard = guard
	c.GSTSvc = NewGSTService(c.gstUC)
	c.GSTSvc.guard = guard
	c.PrinterSvc = NewPrinterService(c.printerUC)
	c.PrinterSvc.guard = guard
	c.AppointmentSvc = NewAppointmentService(c.apptUC)
	c.AppointmentSvc.guard = guard
	c.WhatsAppSvc = NewWhatsAppService(c.whatsappUC)
	c.WhatsAppSvc.guard = guard
	c.MembershipSvc = NewMembershipService(c.membershipUC)
	c.MembershipSvc.guard = guard
	c.CloudBackupSvc = NewCloudBackupService(c.cloudUC)
	c.CloudBackupSvc.guard = guard
}

// SetContext propagates the Wails context to all binding services.
func (c *Container) SetContext(ctx context.Context) {
	licGuard := c.CustomerSvc.licGuard
	if licGuard != nil {
		licGuard.SetContext(ctx)
	}
	c.StaffSvc.SetContext(ctx)
	c.CustomerSvc.SetContext(ctx)
	c.ServiceSvc.SetContext(ctx)
	c.InvoiceSvc.SetContext(ctx)
	c.ExpenseSvc.SetContext(ctx)
	c.ProductSvc.SetContext(ctx)
	c.PerformanceSvc.SetContext(ctx)
	c.CommissionSvc.SetContext(ctx)
	c.SalarySvc.SetContext(ctx)
	c.AnalyticsSvc.SetContext(ctx)
	c.BackupSvc.SetContext(ctx)
	c.LicenseSvc.SetContext(ctx)
	c.UpdateSvc.SetContext(ctx)
	c.ImportSvc.SetContext(ctx)
	c.GSTSvc.SetContext(ctx)
	c.PrinterSvc.SetContext(ctx)
	c.AppointmentSvc.SetContext(ctx)
	c.WhatsAppSvc.SetContext(ctx)
	c.MembershipSvc.SetContext(ctx)
	c.CloudBackupSvc.SetContext(ctx)
	c.AuthSvc.SetContext(ctx)
	c.DiagnosticsSvc.SetContext(ctx)
}

// Bindings returns all service structs to register with Wails.
func (c *Container) Bindings() []interface{} {
	return []interface{}{
		c.StaffSvc,
		c.CustomerSvc,
		c.ServiceSvc,
		c.InvoiceSvc,
		c.ExpenseSvc,
		c.ProductSvc,
		c.PerformanceSvc,
		c.CommissionSvc,
		c.SalarySvc,
		c.AnalyticsSvc,
		c.BackupSvc,
		c.LicenseSvc,
		c.UpdateSvc,
		c.ImportSvc,
		c.GSTSvc,
		c.PrinterSvc,
		c.AppointmentSvc,
		c.WhatsAppSvc,
		c.MembershipSvc,
		c.CloudBackupSvc,
		c.AuthSvc,
		c.DiagnosticsSvc,
	}
}
