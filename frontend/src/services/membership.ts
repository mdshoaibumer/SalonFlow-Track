import { apiClient } from './api-client'
import type { MembershipPlan, MemberSubscription, MembershipStats } from '@/types'

export async function listPlans(): Promise<MembershipPlan[]> {
  const response = await apiClient.get<MembershipPlan[]>('/memberships/plans')
  if (!response.success) throw new Error(response.error?.message || 'Failed to list plans')
  return response.data || []
}

export async function getPlan(id: string): Promise<MembershipPlan> {
  const response = await apiClient.get<MembershipPlan>(`/memberships/plans/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get plan')
  return response.data
}

export async function createPlan(data: Partial<MembershipPlan>): Promise<MembershipPlan> {
  const response = await apiClient.post<MembershipPlan>('/memberships/plans', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create plan')
  return response.data
}

export async function updatePlan(id: string, data: Partial<MembershipPlan>): Promise<MembershipPlan> {
  const response = await apiClient.put<MembershipPlan>(`/memberships/plans/${id}`, data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update plan')
  return response.data
}

export async function deletePlan(id: string): Promise<void> {
  const response = await apiClient.delete(`/memberships/plans/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete plan')
}

export async function sellPlan(data: { plan_id: string; customer_id: string; amount_paid: number }): Promise<MemberSubscription> {
  const response = await apiClient.post<MemberSubscription>('/memberships/sell', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to sell plan')
  return response.data
}

export async function useSession(subscriptionId: string): Promise<MemberSubscription> {
  const response = await apiClient.post<MemberSubscription>(`/memberships/use-session/${subscriptionId}`, {})
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to use session')
  return response.data
}

export async function listSubscriptions(page = 1, perPage = 20): Promise<{ data: MemberSubscription[]; total: number }> {
  const response = await apiClient.get<MemberSubscription[]>(`/memberships/subscriptions?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list subscriptions')
  return { data: response.data || [], total: response.meta?.total || 0 }
}

export async function getMembershipStats(): Promise<MembershipStats> {
  const response = await apiClient.get<MembershipStats>('/memberships/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get stats')
  return response.data
}
