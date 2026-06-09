import { apiClient } from './api-client'
import type { CommissionRule, CommissionStats, CommissionStaffSummary, StaffCommissionOutput, CreateRuleInput, UpdateRuleInput } from '@/types'

export interface ListRulesParams {
  page?: number
  per_page?: number
  rule_type?: string
  target_type?: string
  is_active?: string
}

export interface ListRulesResponse {
  rules: CommissionRule[]
  meta: { page: number; per_page: number; total: number; total_pages: number }
}

export async function listCommissionRules(params: ListRulesParams = {}): Promise<ListRulesResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.rule_type) query.set('rule_type', params.rule_type)
  if (params.target_type) query.set('target_type', params.target_type)
  if (params.is_active) query.set('is_active', params.is_active)
  const qs = query.toString()
  const path = qs ? `/commissions/rules?${qs}` : '/commissions/rules'
  const response = await apiClient.get<CommissionRule[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch rules')
  return {
    rules: response.data || [],
    meta: response.meta as { page: number; per_page: number; total: number; total_pages: number } || { page: 1, per_page: 20, total: 0, total_pages: 0 },
  }
}

export async function getCommissionRuleById(id: string): Promise<CommissionRule> {
  const response = await apiClient.get<CommissionRule>(`/commissions/rules/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch rule')
  return response.data
}

export async function createCommissionRule(input: CreateRuleInput): Promise<CommissionRule> {
  const response = await apiClient.post<CommissionRule>('/commissions/rules', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create rule')
  return response.data
}

export async function updateCommissionRule(id: string, input: UpdateRuleInput): Promise<CommissionRule> {
  const response = await apiClient.put<CommissionRule>(`/commissions/rules/${id}`, input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update rule')
  return response.data
}

export async function deleteCommissionRule(id: string): Promise<void> {
  const response = await apiClient.delete(`/commissions/rules/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete rule')
}

export async function getStaffCommission(staffId: string, params: { date_from?: string; date_to?: string } = {}): Promise<StaffCommissionOutput> {
  const query = new URLSearchParams()
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  const qs = query.toString()
  const path = qs ? `/commissions/staff/${staffId}?${qs}` : `/commissions/staff/${staffId}`
  const response = await apiClient.get<StaffCommissionOutput>(path)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch staff commission')
  return response.data
}

export async function getMonthlyCommission(month?: string): Promise<CommissionStaffSummary[]> {
  const query = new URLSearchParams()
  if (month) query.set('month', month)
  const qs = query.toString()
  const path = qs ? `/commissions/monthly?${qs}` : '/commissions/monthly'
  const response = await apiClient.get<CommissionStaffSummary[]>(path)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch monthly commission')
  return response.data || []
}

export async function getCommissionStats(): Promise<CommissionStats> {
  const response = await apiClient.get<CommissionStats>('/commissions/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch commission stats')
  return response.data
}
