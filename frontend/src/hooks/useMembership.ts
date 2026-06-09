import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as api from '@/services/membership'
import type { MembershipPlan } from '@/types'

export function useMembershipPlans() {
  return useQuery({ queryKey: ['membership-plans'], queryFn: api.listPlans })
}

export function useMembershipPlan(id: string) {
  return useQuery({ queryKey: ['membership-plans', id], queryFn: () => api.getPlan(id), enabled: !!id })
}

export function useCreatePlan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<MembershipPlan>) => api.createPlan(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['membership-plans'] }),
  })
}

export function useUpdatePlan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: Partial<MembershipPlan> & { id: string }) => api.updatePlan(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['membership-plans'] }),
  })
}

export function useDeletePlan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.deletePlan(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['membership-plans'] }),
  })
}

export function useSellPlan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { plan_id: string; customer_id: string; amount_paid: number }) => api.sellPlan(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['subscriptions'] }),
  })
}

export function useUseSession() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (subscriptionId: string) => api.useSession(subscriptionId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['subscriptions'] }),
  })
}

export function useSubscriptions(page = 1) {
  return useQuery({ queryKey: ['subscriptions', page], queryFn: () => api.listSubscriptions(page) })
}

export function useMembershipStats() {
  return useQuery({ queryKey: ['membership-stats'], queryFn: api.getMembershipStats })
}
