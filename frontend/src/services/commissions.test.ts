import { describe, it, expect } from 'vitest'
import { createCommissionRule, getCommissionRuleById, listCommissionRules, updateCommissionRule, deleteCommissionRule, getStaffCommission, getMonthlyCommission, getCommissionStats } from './commissions'

describe('Commissions Service', () => {
  it('creates commission rule', async () => {
    const rule = await createCommissionRule({ rule_type: 'percentage', value: 10 } as any)
    expect(rule.id).toBe('rule1')
  })

  it('gets commission rule by id', async () => {
    const rule = await getCommissionRuleById('rule1')
    expect(rule.rule_type).toBe('percentage')
  })

  it('lists commission rules', async () => {
    const result = await listCommissionRules()
    expect(result.rules).toHaveLength(1)
  })

  it('lists commission rules with params', async () => {
    const result = await listCommissionRules({ rule_type: 'percentage', target_type: 'service', is_active: true, page: 1, per_page: 10 })
    expect(result.rules).toHaveLength(1)
  })

  it('updates commission rule', async () => {
    const rule = await updateCommissionRule('rule1', { value: 15 } as any)
    expect(rule.id).toBe('rule1')
  })

  it('deletes commission rule', async () => {
    await expect(deleteCommissionRule('rule1')).resolves.toBeUndefined()
  })

  it('gets staff commission', async () => {
    const result = await getStaffCommission('staff1')
    expect(result.total).toBe(5000)
  })

  it('gets staff commission with params', async () => {
    const result = await getStaffCommission('staff1', { date_from: '2024-12-01', date_to: '2024-12-31' })
    expect(result.total).toBe(5000)
  })

  it('gets monthly commission', async () => {
    const result = await getMonthlyCommission()
    expect(result.total).toBe(15000)
  })

  it('gets monthly commission with month', async () => {
    const result = await getMonthlyCommission('2024-12')
    expect(result.total).toBe(15000)
  })

  it('gets commission stats', async () => {
    const stats = await getCommissionStats()
    expect(stats.total_commission).toBe(50000)
  })
})
