import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listCommissionRules,
  getCommissionRuleById,
  createCommissionRule,
  updateCommissionRule,
  deleteCommissionRule,
  getStaffCommission,
  getMonthlyCommission,
  getCommissionStats,
  type ListRulesParams,
} from '@/services/commissions'
import type { CreateRuleInput, UpdateRuleInput } from '@/types'

export function useCommissionRules(params: ListRulesParams = {}) {
  return useQuery({
    queryKey: ['commissions', 'rules', params],
    queryFn: () => listCommissionRules(params),
  })
}

export function useCommissionRuleById(id: string) {
  return useQuery({
    queryKey: ['commissions', 'rules', id],
    queryFn: () => getCommissionRuleById(id),
    enabled: !!id,
  })
}

export function useCreateCommissionRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateRuleInput) => createCommissionRule(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commissions'] })
    },
  })
}

export function useUpdateCommissionRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateRuleInput }) => updateCommissionRule(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commissions'] })
    },
  })
}

export function useDeleteCommissionRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteCommissionRule(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commissions'] })
    },
  })
}

export function useStaffCommission(staffId: string, params: { date_from?: string; date_to?: string } = {}) {
  return useQuery({
    queryKey: ['commissions', 'staff', staffId, params],
    queryFn: () => getStaffCommission(staffId, params),
    enabled: !!staffId,
  })
}

export function useMonthlyCommission(month?: string) {
  return useQuery({
    queryKey: ['commissions', 'monthly', month],
    queryFn: () => getMonthlyCommission(month),
  })
}

export function useCommissionStats() {
  return useQuery({
    queryKey: ['commissions', 'stats'],
    queryFn: getCommissionStats,
  })
}
