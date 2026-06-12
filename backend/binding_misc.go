package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// AppointmentService exposes appointment operations to the Wails frontend.
type AppointmentService struct {
	ctx   context.Context
	uc    *usecase.AppointmentUseCase
	guard *PermissionGuard
}

func NewAppointmentService(uc *usecase.AppointmentUseCase) *AppointmentService {
	return &AppointmentService{uc: uc}
}

func (s *AppointmentService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *AppointmentService) CreateAppointment(appt *domain.Appointment, services []domain.AppointmentService) error {
	return s.uc.Create(s.ctx, appt, services)
}

func (s *AppointmentService) UpdateAppointment(appt *domain.Appointment, services []domain.AppointmentService) error {
	return s.uc.Update(s.ctx, appt, services)
}

func (s *AppointmentService) UpdateAppointmentStatus(id string, status, note string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.UpdateStatus(s.ctx, uid, status, note)
}

func (s *AppointmentService) DeleteAppointment(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.Delete(s.ctx, uid)
}

func (s *AppointmentService) GetAppointment(id string) (*domain.Appointment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetByID(s.ctx, uid)
}

func (s *AppointmentService) ListAppointments(filter domain.AppointmentFilter) ([]domain.Appointment, int, error) {
	return s.uc.List(s.ctx, filter)
}

func (s *AppointmentService) GetAppointmentHistory(id string) ([]domain.AppointmentHistory, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetHistory(s.ctx, uid)
}

// WhatsAppService exposes WhatsApp messaging to the Wails frontend.
type WhatsAppService struct {
	ctx   context.Context
	uc    *usecase.WhatsAppUseCase
	guard *PermissionGuard
}

func NewWhatsAppService(uc *usecase.WhatsAppUseCase) *WhatsAppService {
	return &WhatsAppService{uc: uc}
}

func (s *WhatsAppService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *WhatsAppService) CreateTemplate(tmpl *domain.WhatsAppTemplate) error {
	return s.uc.CreateTemplate(s.ctx, tmpl)
}

func (s *WhatsAppService) UpdateTemplate(tmpl *domain.WhatsAppTemplate) error {
	return s.uc.UpdateTemplate(s.ctx, tmpl)
}

func (s *WhatsAppService) DeleteTemplate(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteTemplate(s.ctx, uid)
}

func (s *WhatsAppService) ListTemplates(category string) ([]domain.WhatsAppTemplate, error) {
	return s.uc.ListTemplates(s.ctx, category)
}

func (s *WhatsAppService) SendMessage(templateID, phone, name string, variables map[string]string) (*domain.WhatsAppMessage, error) {
	return s.uc.SendMessage(s.ctx, templateID, phone, name, variables)
}

func (s *WhatsAppService) ListMessages(limit, offset int, status string) ([]domain.WhatsAppMessage, int, error) {
	return s.uc.ListMessages(s.ctx, limit, offset, status)
}

func (s *WhatsAppService) GetWhatsAppStats() (*domain.WAMessageStats, error) {
	return s.uc.GetStats(s.ctx)
}

func (s *WhatsAppService) CreateRule(rule *domain.AutomationRule) error {
	return s.uc.CreateRule(s.ctx, rule)
}

func (s *WhatsAppService) UpdateRule(rule *domain.AutomationRule) error {
	return s.uc.UpdateRule(s.ctx, rule)
}

func (s *WhatsAppService) DeleteRule(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteRule(s.ctx, uid)
}

func (s *WhatsAppService) ListRules() ([]domain.AutomationRule, error) {
	return s.uc.ListRules(s.ctx)
}

// MembershipService exposes membership operations to the Wails frontend.
type MembershipService struct {
	ctx   context.Context
	uc    *usecase.MembershipUseCase
	guard *PermissionGuard
}

func NewMembershipService(uc *usecase.MembershipUseCase) *MembershipService {
	return &MembershipService{uc: uc}
}

func (s *MembershipService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *MembershipService) CreatePlan(plan *domain.MembershipPlan, services []domain.PackageService) error {
	return s.uc.CreatePlan(s.ctx, plan, services)
}

func (s *MembershipService) UpdatePlan(plan *domain.MembershipPlan, services []domain.PackageService) error {
	return s.uc.UpdatePlan(s.ctx, plan, services)
}

func (s *MembershipService) DeletePlan(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeletePlan(s.ctx, uid)
}

func (s *MembershipService) GetPlan(id string) (*domain.MembershipPlan, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetPlan(s.ctx, uid)
}

func (s *MembershipService) ListPlans(planType string) ([]domain.MembershipPlan, error) {
	return s.uc.ListPlans(s.ctx, planType)
}

func (s *MembershipService) SellPlan(customerID string, planID string, amountPaid float64) (*domain.MemberSubscription, error) {
	uid, err := uuid.Parse(planID)
	if err != nil {
		return nil, err
	}
	return s.uc.SellPlan(s.ctx, customerID, uid, amountPaid)
}

func (s *MembershipService) UseSession(subscriptionID string) error {
	uid, err := uuid.Parse(subscriptionID)
	if err != nil {
		return err
	}
	return s.uc.UseSession(s.ctx, uid)
}

func (s *MembershipService) ListSubscriptions(customerID, status string, limit, offset int) ([]domain.MemberSubscription, int, error) {
	return s.uc.ListSubscriptions(s.ctx, customerID, status, limit, offset)
}

func (s *MembershipService) GetMembershipStats() (*domain.MembershipStats, error) {
	return s.uc.GetStats(s.ctx)
}

// CloudBackupService exposes cloud backup operations to the Wails frontend.
type CloudBackupService struct {
	ctx   context.Context
	uc    *usecase.CloudBackupUseCase
	guard *PermissionGuard
}

func NewCloudBackupService(uc *usecase.CloudBackupUseCase) *CloudBackupService {
	return &CloudBackupService{uc: uc}
}

func (s *CloudBackupService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *CloudBackupService) GetConfig() (*domain.CloudBackupConfig, error) {
	return s.uc.GetConfig(s.ctx)
}

func (s *CloudBackupService) SaveConfig(cfg *domain.CloudBackupConfig) error {
	return s.uc.SaveConfig(s.ctx, cfg)
}

func (s *CloudBackupService) TestConnection() error {
	return s.uc.TestConnection(s.ctx)
}

func (s *CloudBackupService) BackupNow() (*domain.CloudBackupHistory, error) {
	return s.uc.BackupNow(s.ctx)
}

func (s *CloudBackupService) Restore(historyID string) error {
	uid, err := uuid.Parse(historyID)
	if err != nil {
		return err
	}
	return s.uc.Restore(s.ctx, uid)
}

func (s *CloudBackupService) ListHistory(limit, offset int) ([]domain.CloudBackupHistory, int, error) {
	return s.uc.ListHistory(s.ctx, limit, offset)
}

func (s *CloudBackupService) GetCloudBackupStats() (*domain.CloudBackupStats, error) {
	return s.uc.GetStats(s.ctx)
}
