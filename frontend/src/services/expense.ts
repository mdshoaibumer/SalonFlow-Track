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

export async function listCategories(activeOnly = true): Promise<ExpenseCategory[]> {
  return window.go.main.ExpenseService.ListCategories(activeOnly)
}

export async function createCategory(name: string, description: string): Promise<ExpenseCategory> {
  return window.go.main.ExpenseService.CreateCategory(name, description)
}

export async function updateCategory(id: string, data: { name: string; description: string; is_active: boolean }): Promise<ExpenseCategory> {
  return window.go.main.ExpenseService.UpdateCategory(id, data.name, data.description, data.is_active)
}

export async function createExpense(input: CreateExpenseInput): Promise<Expense> {
  return window.go.main.ExpenseService.CreateExpense(input)
}

export async function listExpenses(params: ListExpensesParams = {}): Promise<ListExpensesResponse> {
  const result = await window.go.main.ExpenseService.ListExpenses({
    CategoryID: params.category_id || '',
    Status: params.status || '',
    PaymentMethod: params.payment_method || '',
    DateFrom: params.date_from || '',
    DateTo: params.date_to || '',
    Search: params.search || '',
    Page: params.page || 1,
    PerPage: params.per_page || 20,
  })
  return {
    expenses: result.expenses || [],
    meta: {
      page: result.page,
      per_page: result.per_page,
      total: result.total,
      total_pages: result.total_pages,
    },
  }
}

export async function getExpenseById(id: string): Promise<Expense> {
  return window.go.main.ExpenseService.GetExpense(id)
}

export async function updateExpense(id: string, input: UpdateExpenseInput): Promise<Expense> {
  return window.go.main.ExpenseService.UpdateExpense(id, input)
}

export async function deleteExpense(id: string): Promise<void> {
  await window.go.main.ExpenseService.DeleteExpense(id)
}

export async function getExpenseStats(): Promise<ExpenseStats> {
  return window.go.main.ExpenseService.GetExpenseStats()
}

export async function getExpenseReport(dateFrom?: string, dateTo?: string): Promise<ExpenseReport> {
  return window.go.main.ExpenseService.GetExpenseReport({ DateFrom: dateFrom || '', DateTo: dateTo || '' })
}

export async function getProfitLoss(dateFrom?: string, dateTo?: string): Promise<ProfitLoss> {
  return window.go.main.ExpenseService.GetProfitLoss({ DateFrom: dateFrom || '', DateTo: dateTo || '' })
}

export async function getMonthlyTrend(months: number): Promise<MonthlyTrend[]> {
  return window.go.main.ExpenseService.GetMonthlyTrend(months)
}
