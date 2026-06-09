import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listCategories,
  createCategory,
  updateCategory,
  createExpense,
  listExpenses,
  getExpenseById,
  updateExpense,
  deleteExpense,
  getExpenseStats,
  getProfitLoss,
  getMonthlyTrend,
  getExpenseReport,
  type ListExpensesParams,
} from '@/services/expense'
import type { CreateExpenseInput, UpdateExpenseInput } from '@/types'

export function useExpenseCategories(activeOnly = true) {
  return useQuery({
    queryKey: ['expense-categories', activeOnly],
    queryFn: () => listCategories(activeOnly),
  })
}

export function useCreateCategory() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: { name: string; description: string }) => createCategory(input.name, input.description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expense-categories'] })
    },
  })
}

export function useUpdateCategory() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: { id: string; name: string; description: string; is_active: boolean }) =>
      updateCategory(input.id, { name: input.name, description: input.description, is_active: input.is_active }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expense-categories'] })
    },
  })
}

export function useExpenseList(params: ListExpensesParams = {}) {
  return useQuery({
    queryKey: ['expenses', params],
    queryFn: () => listExpenses(params),
  })
}

export function useExpenseById(id: string) {
  return useQuery({
    queryKey: ['expenses', id],
    queryFn: () => getExpenseById(id),
    enabled: !!id,
  })
}

export function useCreateExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateExpenseInput) => createExpense(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({ queryKey: ['expense-stats'] })
    },
  })
}

export function useUpdateExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateExpenseInput }) => updateExpense(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({ queryKey: ['expense-stats'] })
    },
  })
}

export function useDeleteExpense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteExpense(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({ queryKey: ['expense-stats'] })
    },
  })
}

export function useExpenseStats() {
  return useQuery({
    queryKey: ['expense-stats'],
    queryFn: getExpenseStats,
  })
}

export function useProfitLoss(dateFrom?: string, dateTo?: string) {
  return useQuery({
    queryKey: ['profit-loss', dateFrom, dateTo],
    queryFn: () => getProfitLoss(dateFrom, dateTo),
  })
}

export function useMonthlyTrend(months = 6) {
  return useQuery({
    queryKey: ['expense-trend', months],
    queryFn: () => getMonthlyTrend(months),
  })
}

export function useExpenseReport(dateFrom?: string, dateTo?: string) {
  return useQuery({
    queryKey: ['expense-report', dateFrom, dateTo],
    queryFn: () => getExpenseReport(dateFrom, dateTo),
  })
}
