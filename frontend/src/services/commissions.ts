import type { CommissionRule, CommissionStats, CreateRuleInput, UpdateRuleInput } from '@/types'

export interface ListRulesParams {
  rule_type?: string
  target_type?: string
  is_active?: boolean
  page?: number
  per_page?: number
}

export async function createCommissionRule(input: CreateRuleInput): Promise<CommissionRule> {
  return window.go.main.CommissionService.CreateRule(input)
}

export async function getCommissionRuleById(id: string): Promise<CommissionRule> {
  return window.go.main.CommissionService.GetRule(id)
}

export async function listCommissionRules(params: ListRulesParams = {}) {
  return window.go.main.CommissionService.ListRules({
    rule_type: params.rule_type || '',
    target_type: params.target_type || '',
    is_active: params.is_active ?? null,
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
}

export async function updateCommissionRule(id: string, input: UpdateRuleInput): Promise<CommissionRule> {
  return window.go.main.CommissionService.UpdateRule(id, input)
}

export async function deleteCommissionRule(id: string): Promise<void> {
  await window.go.main.CommissionService.DeleteRule(id)
}

export async function getStaffCommission(staffId: string, params: { date_from?: string; date_to?: string } = {}) {
  return window.go.main.CommissionService.GetStaffCommission({
    staff_id: staffId,
    date_from: params.date_from || '',
    date_to: params.date_to || '',
  })
}

export async function getMonthlyCommission(month?: string) {
  return window.go.main.CommissionService.GetMonthlyCommission({ month: month || '' })
}

export async function getCommissionStats(): Promise<CommissionStats> {
  return window.go.main.CommissionService.GetCommissionStats()
}
