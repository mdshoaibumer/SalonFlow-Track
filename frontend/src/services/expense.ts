import { apiClient } from './api-client'
import type {
  Expense,
  ExpenseCategory,
  CreateExpenseInput,
  UpdateExpenseInput,
  ExpenseStats,
  ProfitLoss,
  MonthlyTrend,
  ExpenseReport,
} from '@/types'

export interface ListExpensesParams {
  page?: number
  per_page?: number
  category_id?: string
  status?: string
  payment_method?: string
  date_from?: string
  date_to?: string
  search?: string
}

export interface ListExpensesResponse {
  expenses: Expense[]
  meta: { page: number; per_page: number; total: number; total_pages: number }
}

// --- Categories ---

export async function listCategories(activeOnly = true): Promise<ExpenseCategory[]> {
  const response = await apiClient.get<ExpenseCategory[]>(`/expenses/categories?active_only=${activeOnly}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch categories')
  return response.data || []
}

export async function createCategory(name: string, description: string): Promise<ExpenseCategory> {
  const response = await apiClient.post<ExpenseCategory>('/expenses/categories', { name, description })
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create category')
  return response.data
}

export async function updateCategory(id: string, data: { name: string; description: string; is_active: boolean }): Promise<ExpenseCategory> {
  const response = await apiClient.put<ExpenseCategory>(`/expenses/categories/${id}`, data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update category')
  return response.data
}

// --- Expenses ---

export async function createExpense(input: CreateExpenseInput): Promise<Expense> {
  const response = await apiClient.post<Expense>('/expenses', input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create expense')
  return response.data
}

export async function listExpenses(params: ListExpensesParams = {}): Promise<ListExpensesResponse> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.per_page) query.set('per_page', String(params.per_page))
  if (params.category_id) query.set('category_id', params.category_id)
  if (params.status) query.set('status', params.status)
  if (params.payment_method) query.set('payment_method', params.payment_method)
  if (params.date_from) query.set('date_from', params.date_from)
  if (params.date_to) query.set('date_to', params.date_to)
  if (params.search) query.set('search', params.search)

  const response = await apiClient.get<Expense[]>(`/expenses?${query.toString()}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch expenses')
  return {
    expenses: response.data || [],
    meta: {
      page: response.meta?.page || 1,
      per_page: response.meta?.per_page || 20,
      total: response.meta?.total || 0,
      total_pages: response.meta?.total_pages || 0,
    },
  }
}

export async function getExpenseById(id: string): Promise<Expense> {
  const response = await apiClient.get<Expense>(`/expenses/${id}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch expense')
  return response.data
}

export async function updateExpense(id: string, input: UpdateExpenseInput): Promise<Expense> {
  const response = await apiClient.put<Expense>(`/expenses/${id}`, input)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update expense')
  return response.data
}

export async function deleteExpense(id: string): Promise<void> {
  const response = await apiClient.delete(`/expenses/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete expense')
}

// --- Stats & Reports ---

export async function getExpenseStats(): Promise<ExpenseStats> {
  const response = await apiClient.get<ExpenseStats>('/expenses/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch expense stats')
  return response.data
}

export async function getProfitLoss(dateFrom?: string, dateTo?: string): Promise<ProfitLoss> {
  const query = new URLSearchParams()
  if (dateFrom) query.set('date_from', dateFrom)
  if (dateTo) query.set('date_to', dateTo)
  const response = await apiClient.get<ProfitLoss>(`/expenses/profit-loss?${query.toString()}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch P&L')
  return response.data
}

export async function getMonthlyTrend(months = 6): Promise<MonthlyTrend[]> {
  const response = await apiClient.get<MonthlyTrend[]>(`/expenses/trend?months=${months}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to fetch trend')
  return response.data || []
}

export async function getExpenseReport(dateFrom?: string, dateTo?: string): Promise<ExpenseReport> {
  const query = new URLSearchParams()
  if (dateFrom) query.set('date_from', dateFrom)
  if (dateTo) query.set('date_to', dateTo)
  const response = await apiClient.get<ExpenseReport>(`/expenses/report?${query.toString()}`)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to fetch report')
  return response.data
}
