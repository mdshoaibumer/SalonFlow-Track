import { describe, it, expect, vi } from 'vitest'
import { listCategories, createCategory, updateCategory, createExpense, listExpenses, getExpenseById, updateExpense, deleteExpense, getExpenseStats, getExpenseReport, getProfitLoss, getMonthlyTrend } from './expense'

describe('Expense Service', () => {
  it('lists categories', async () => {
    const cats = await listCategories()
    expect(cats).toHaveLength(2)
    expect(cats[0].name).toBe('Rent')
  })

  it('lists categories with activeOnly false', async () => {
    const cats = await listCategories(false)
    expect(cats).toHaveLength(2)
  })

  it('creates category', async () => {
    const cat = await createCategory('New', 'A new category')
    expect(cat.id).toBe('cat3')
  })

  it('updates category', async () => {
    const cat = await updateCategory('cat1', { name: 'Rent', description: 'Monthly', is_active: true })
    expect(cat.id).toBe('cat1')
  })

  it('creates expense', async () => {
    const expense = await createExpense({ amount: 5000, category_id: 'cat1' } as any)
    expect(expense.id).toBe('exp1')
  })

  it('lists expenses', async () => {
    const result = await listExpenses()
    expect(result.expenses).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists expenses with params', async () => {
    const result = await listExpenses({ category_id: 'cat1', status: 'paid', payment_method: 'cash', date_from: '2024-01-01', date_to: '2024-12-31', search: 'rent', page: 1, per_page: 10 })
    expect(result.expenses).toHaveLength(1)
  })

  it('gets expense by id', async () => {
    const expense = await getExpenseById('exp1')
    expect(expense.amount).toBe(25000)
  })

  it('updates expense', async () => {
    const expense = await updateExpense('exp1', { amount: 30000 } as any)
    expect(expense.id).toBe('exp1')
  })

  it('deletes expense', async () => {
    await expect(deleteExpense('exp1')).resolves.toBeUndefined()
  })

  it('gets expense stats', async () => {
    const stats = await getExpenseStats()
    expect(stats.monthly_expenses).toBe(85000)
  })

  it('gets expense report', async () => {
    const report = await getExpenseReport('2024-01-01', '2024-12-31')
    expect(report).toBeDefined()
  })

  it('gets expense report with no params', async () => {
    const report = await getExpenseReport()
    expect(report).toBeDefined()
  })

  it('gets profit loss', async () => {
    const pl = await getProfitLoss('2024-01-01', '2024-12-31')
    expect(pl.revenue).toBe(300000)
  })

  it('gets profit loss with no params', async () => {
    const pl = await getProfitLoss()
    expect(pl.revenue).toBe(300000)
  })

  it('gets monthly trend', async () => {
    const trend = await getMonthlyTrend(6)
    expect(trend).toEqual([])
  })

  it('listExpenses returns empty when API returns undefined expenses', async () => {
    vi.spyOn(window.go.main.ExpenseService, 'ListExpenses').mockResolvedValueOnce({ expenses: undefined, total: 0 } as any)
    const r = await listExpenses()
    expect(r.expenses).toEqual([])
  })
})
