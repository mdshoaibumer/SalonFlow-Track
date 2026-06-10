import type { MembershipPlan, MemberSubscription, MembershipStats } from '@/types'

export async function createPlan(data: any): Promise<void> {
  const services = data.services || []
  await window.go.main.MembershipService.CreatePlan(data, services)
}

export async function updatePlan(id: string, data: any): Promise<void> {
  const services = data.services || []
  await window.go.main.MembershipService.UpdatePlan({ ...data, id }, services)
}

export async function deletePlan(id: string): Promise<void> {
  await window.go.main.MembershipService.DeletePlan(id)
}

export async function getPlan(id: string): Promise<MembershipPlan> {
  return window.go.main.MembershipService.GetPlan(id)
}

export async function listPlans(): Promise<MembershipPlan[]> {
  return window.go.main.MembershipService.ListPlans('')
}

export async function sellPlan(data: { plan_id: string; customer_id: string; amount_paid: number }): Promise<MemberSubscription> {
  return window.go.main.MembershipService.SellPlan(data.customer_id, data.plan_id, data.amount_paid)
}

export async function useSession(subscriptionId: string): Promise<void> {
  await window.go.main.MembershipService.UseSession(subscriptionId)
}

export async function listSubscriptions(page = 1, perPage = 20) {
  const offset = (page - 1) * perPage
  const [subscriptions, total] = await window.go.main.MembershipService.ListSubscriptions('', '', perPage, offset)
  return { data: subscriptions || [], meta: { page, per_page: perPage, total, total_pages: Math.ceil(total / perPage) } }
}

export async function getMembershipStats(): Promise<MembershipStats> {
  return window.go.main.MembershipService.GetMembershipStats()
}
