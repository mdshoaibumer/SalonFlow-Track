import { describe, it, expect, vi } from 'vitest'
import { createPlan, updatePlan, deletePlan, getPlan, listPlans, sellPlan, useSession, listSubscriptions, getMembershipStats } from './membership'

describe('Membership Service', () => {
  it('creates plan', async () => {
    await expect(createPlan({ name: 'Gold', price: 5000, services: ['svc1'] })).resolves.toBeUndefined()
  })

  it('creates plan without services', async () => {
    await expect(createPlan({ name: 'Gold', price: 5000 })).resolves.toBeUndefined()
  })

  it('updates plan', async () => {
    await expect(updatePlan('plan1', { name: 'Gold Pro', services: ['svc1'] })).resolves.toBeUndefined()
  })

  it('updates plan without services', async () => {
    await expect(updatePlan('plan1', { name: 'Gold Pro' })).resolves.toBeUndefined()
  })

  it('deletes plan', async () => {
    await expect(deletePlan('plan1')).resolves.toBeUndefined()
  })

  it('gets plan', async () => {
    const plan = await getPlan('plan1')
    expect(plan.name).toBe('Gold')
  })

  it('lists plans', async () => {
    const plans = await listPlans()
    expect(plans).toHaveLength(1)
  })

  it('sells plan', async () => {
    const sub = await sellPlan({ plan_id: 'plan1', customer_id: 'cust1', amount_paid: 5000 })
    expect(sub.id).toBe('sub1')
  })

  it('uses session', async () => {
    await expect(useSession('sub1')).resolves.toBeUndefined()
  })

  it('lists subscriptions', async () => {
    const result = await listSubscriptions()
    expect(result.data).toEqual([])
    expect(result.meta.total).toBe(0)
  })

  it('lists subscriptions with pagination', async () => {
    const result = await listSubscriptions(2, 10)
    expect(result.data).toEqual([])
  })

  it('gets membership stats', async () => {
    const stats = await getMembershipStats()
    expect(stats.total_plans).toBe(3)
  })

  it('listSubscriptions returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.MembershipService, 'ListSubscriptions').mockResolvedValueOnce([undefined as any, 0])
    const r = await listSubscriptions()
    expect(r.data).toEqual([])
  })
})
