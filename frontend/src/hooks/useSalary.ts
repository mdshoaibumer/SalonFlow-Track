import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  generateSalary,
  listSalaries,
  getSalaryById,
  paySalary,
  listSalaryCycles,
  getSalaryStats,
  createAdvance,
  listAdvances,
  approveAdvance,
  rejectAdvance,
  type ListAdvancesParams,
} from '@/services/salary'
import type { GenerateSalaryInput, CreateAdvanceInput } from '@/types'

export function useSalaryList(month: number, year: number) {
  return useQuery({
    queryKey: ['salary', 'records', month, year],
    queryFn: () => listSalaries(month, year),
    enabled: month > 0 && year > 0,
  })
}

export function useSalaryById(id: string) {
  return useQuery({
    queryKey: ['salary', id],
    queryFn: () => getSalaryById(id),
    enabled: !!id,
  })
}

export function useSalaryCycles(year?: number) {
  return useQuery({
    queryKey: ['salary', 'cycles', year],
    queryFn: () => listSalaryCycles(year),
  })
}

export function useSalaryStats() {
  return useQuery({
    queryKey: ['salary', 'stats'],
    queryFn: getSalaryStats,
  })
}

export function useGenerateSalary() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: GenerateSalaryInput) => generateSalary(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['salary'] })
    },
  })
}

export function usePaySalary() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => paySalary(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['salary'] })
    },
  })
}

export function useAdvanceList(params: ListAdvancesParams = {}) {
  return useQuery({
    queryKey: ['salary', 'advances', params],
    queryFn: () => listAdvances(params),
  })
}

export function useCreateAdvance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateAdvanceInput) => createAdvance(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['salary', 'advances'] })
    },
  })
}

export function useApproveAdvance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => approveAdvance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['salary', 'advances'] })
    },
  })
}

export function useRejectAdvance() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => rejectAdvance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['salary', 'advances'] })
    },
  })
}
