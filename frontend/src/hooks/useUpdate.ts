import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as updateService from '@/services/update'

const UPDATE_KEYS = {
  all: ['update'] as const,
  status: ['update', 'status'] as const,
  history: (page: number) => ['update', 'history', page] as const,
}

export function useUpdateStatus() {
  return useQuery({
    queryKey: UPDATE_KEYS.status,
    queryFn: () => updateService.getUpdateStatus(),
  })
}

export function useUpdateHistory(page = 1) {
  return useQuery({
    queryKey: UPDATE_KEYS.history(page),
    queryFn: () => updateService.listUpdateHistory(page),
  })
}

export function useCheckForUpdate() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => updateService.checkForUpdate(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: UPDATE_KEYS.all })
    },
  })
}

export function useDownloadUpdate() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => updateService.downloadUpdate(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: UPDATE_KEYS.all })
    },
  })
}

export function useInstallUpdate() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => updateService.installUpdate(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: UPDATE_KEYS.all })
    },
  })
}
