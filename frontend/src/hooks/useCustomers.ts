import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listCustomers, getCustomerById, createCustomer, updateCustomer, deleteCustomer, getCustomerStats, type ListCustomerParams } from '@/services/customers'
import type { CreateCustomerInput, UpdateCustomerInput } from '@/types'

export function useCustomerList(params: ListCustomerParams = {}) {
  return useQuery({
    queryKey: ['customers', params],
    queryFn: () => listCustomers(params),
  })
}

export function useCustomerById(id: string) {
  return useQuery({
    queryKey: ['customers', id],
    queryFn: () => getCustomerById(id),
    enabled: !!id,
  })
}

export function useCustomerStats() {
  return useQuery({
    queryKey: ['customers', 'stats'],
    queryFn: getCustomerStats,
  })
}

export function useCreateCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateCustomerInput) => createCustomer(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
    },
  })
}

export function useUpdateCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateCustomerInput }) => updateCustomer(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
    },
  })
}

export function useDeleteCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteCustomer(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
    },
  })
}
