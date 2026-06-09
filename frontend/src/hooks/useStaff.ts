import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listStaff,
  getStaffById,
  createStaff,
  updateStaff,
  deleteStaff,
  getStaffStats,
  type ListStaffParams,
} from '@/services/staff'
import type { CreateStaffInput, UpdateStaffInput } from '@/types'

export function useStaffList(params: ListStaffParams = {}) {
  return useQuery({
    queryKey: ['staff', 'list', params],
    queryFn: () => listStaff(params),
  })
}

export function useStaffById(id: string | undefined) {
  return useQuery({
    queryKey: ['staff', 'detail', id],
    queryFn: () => getStaffById(id!),
    enabled: !!id,
  })
}

export function useStaffStats() {
  return useQuery({
    queryKey: ['staff', 'stats'],
    queryFn: getStaffStats,
  })
}

export function useCreateStaff() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateStaffInput) => createStaff(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] })
    },
  })
}

export function useUpdateStaff() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateStaffInput }) =>
      updateStaff(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] })
    },
  })
}

export function useDeleteStaff() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteStaff(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] })
    },
  })
}
