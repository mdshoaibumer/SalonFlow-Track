package main

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/salonflow/salonflow-track/internal/adapters/backup"
	"github.com/salonflow/salonflow-track/internal/adapters/cloudbackup"
	"github.com/salonflow/salonflow-track/internal/adapters/gst"
	"github.com/salonflow/salonflow-track/internal/adapters/handler"
	"github.com/salonflow/salonflow-track/internal/adapters/importer"
	"github.com/salonflow/salonflow-track/internal/adapters/license"
	"github.com/salonflow/salonflow-track/internal/adapters/printer"
	"github.com/salonflow/salonflow-track/internal/adapters/repository/sqlite"
	"github.com/salonflow/salonflow-track/internal/adapters/update"
	"github.com/salonflow/salonflow-track/internal/config"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/internal/database"
	appmw "github.com/salonflow/salonflow-track/pkg/middleware"
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

	// Use Cases
	staffUC      *usecase.StaffUseCase
	serviceUC    *usecase.ServiceUseCase
	customerUC   *usecase.CustomerUseCase
	invoiceUC    *usecase.InvoiceUseCase
	perfUC       *usecase.PerformanceUseCase
	commissionUC *usecase.CommissionUseCase
	salaryUC     *usecase.SalaryUseCase
	expenseUC    *usecase.ExpenseUseCase
	productUC    *usecase.ProductUseCase
	analyticsUC  *usecase.AnalyticsUseCase
	backupUC     *usecase.BackupUseCase
	licenseUC    *usecase.LicenseUseCase
	updateUC     *usecase.UpdateUseCase
	importUC     *usecase.ImportUseCase
	gstUC        *usecase.GSTUseCase
	printerUC    *usecase.PrinterUseCase
	apptUC       *usecase.AppointmentUseCase
	whatsappUC   *usecase.WhatsAppUseCase
	membershipUC *usecase.MembershipUseCase
	cloudUC      *usecase.CloudBackupUseCase
}

// NewContainer builds the full dependency graph.
func NewContainer(cfg *config.Config, log *slog.Logger, db *database.DB) *Container {
	c := &Container{
		cfg: cfg,
		log: log,
		db:  db,
	}
	c.initRepositories()
	c.initUseCases()
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
}

// HTTPServer returns a configured *http.Server ready to listen.
func (c *Container) HTTPServer() *http.Server {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(appmw.RequestLogger(c.log))
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://wails.localhost:*", "http://wails.localhost", "wails://*", "null"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Handlers
	healthH := handler.NewHealthHandler(c.db, c.cfg, c.log)
	staffH := handler.NewStaffHandler(c.staffUC)
	serviceH := handler.NewServiceHandler(c.serviceUC)
	customerH := handler.NewCustomerHandler(c.customerUC)
	invoiceH := handler.NewInvoiceHandler(c.invoiceUC)
	perfH := handler.NewPerformanceHandler(c.perfUC)
	commissionH := handler.NewCommissionHandler(c.commissionUC)
	salaryH := handler.NewSalaryHandler(c.salaryUC)
	expenseH := handler.NewExpenseHandler(c.expenseUC)
	productH := handler.NewProductHandler(c.productUC)
	analyticsH := handler.NewAnalyticsHandler(c.analyticsUC)
	backupH := handler.NewBackupHandler(c.backupUC)
	licenseH := handler.NewLicenseHandler(c.licenseUC)
	updateH := handler.NewUpdateHandler(c.updateUC)
	importH := handler.NewImportHandler(c.importUC, c.importEngine.UploadDir())
	gstH := handler.NewGSTHandler(c.gstUC)
	printerH := handler.NewPrinterHandler(c.printerUC)
	apptH := handler.NewAppointmentHandler(c.apptUC)
	whatsappH := handler.NewWhatsAppHandler(c.whatsappUC)
	membershipH := handler.NewMembershipHandler(c.membershipUC)
	cloudH := handler.NewCloudBackupHandler(c.cloudUC)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthH.Check)
		r.Mount("/staff", staffH.Routes())
		r.Mount("/services", serviceH.Routes())
		r.Mount("/customers", customerH.Routes())
		r.Mount("/invoices", invoiceH.Routes())
		r.Mount("/performance", perfH.Routes())
		r.Mount("/commissions", commissionH.Routes())
		r.Mount("/salary", salaryH.Routes())
		r.Mount("/expenses", expenseH.Routes())
		r.Mount("/products", productH.Routes())
		r.Mount("/reports", analyticsH.Routes())
		r.Mount("/backups", backupH.Routes())
		r.Mount("/license", licenseH.Routes())
		r.Mount("/update", updateH.Routes())
		r.Mount("/import", importH.Routes())
		r.Mount("/gst", gstH.Routes())
		r.Mount("/print", printerH.Routes())
		r.Mount("/appointments", apptH.Routes())
		r.Mount("/whatsapp", whatsappH.Routes())
		r.Mount("/memberships", membershipH.Routes())
		r.Mount("/cloud-backup", cloudH.Routes())
	})

	addr := net.JoinHostPort(c.cfg.Server.Host, c.cfg.Server.Port)
	return &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(c.cfg.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(c.cfg.Server.WriteTimeoutSec) * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
