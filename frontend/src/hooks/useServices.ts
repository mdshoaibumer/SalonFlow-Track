import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listServices, getServiceById, createService, updateService, deleteService, getServiceStats, type ListServiceParams } from '@/services/services'
import type { CreateServiceInput, UpdateServiceInput } from '@/types'

export function useServiceList(params: ListServiceParams = {}) {
  return useQuery({
    queryKey: ['services', params],
    queryFn: () => listServices(params),
  })
}

export function useServiceById(id: string) {
  return useQuery({
    queryKey: ['services', id],
    queryFn: () => getServiceById(id),
    enabled: !!id,
  })
}

export function useServiceStats() {
  return useQuery({
    queryKey: ['services', 'stats'],
    queryFn: getServiceStats,
  })
}

export function useCreateService() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateServiceInput) => createService(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services'] })
    },
  })
}

export function useUpdateService() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateServiceInput }) => updateService(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services'] })
    },
  })
}

export function useDeleteService() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteService(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services'] })
    },
  })
}
