import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as licenseService from '@/services/license'

const LICENSE_KEYS = {
  all: ['license'] as const,
  status: ['license', 'status'] as const,
  events: (page: number) => ['license', 'events', page] as const,
  notifications: ['license', 'notifications'] as const,
  deviceID: ['license', 'deviceID'] as const,
}

export function useLicenseStatus() {
  return useQuery({
    queryKey: LICENSE_KEYS.status,
    queryFn: () => licenseService.getLicenseStatus(),
  })
}

export function useLicenseEvents(page = 1) {
  return useQuery({
    queryKey: LICENSE_KEYS.events(page),
    queryFn: () => licenseService.listLicenseEvents(page),
  })
}

export function useLicenseNotifications(unreadOnly = true) {
  return useQuery({
    queryKey: LICENSE_KEYS.notifications,
    queryFn: () => licenseService.getNotifications(unreadOnly),
    refetchInterval: 60000,
  })
}

export function useDeviceID() {
  return useQuery({
    queryKey: LICENSE_KEYS.deviceID,
    queryFn: () => licenseService.getDeviceID(),
    staleTime: Infinity,
  })
}

export function useValidateLicense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => licenseService.validateLicense(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.all })
    },
  })
}

export function useActivateLicense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: { licenseKey: string; customerName: string; salonName: string }) =>
      licenseService.activateLicense(data.licenseKey, data.customerName, data.salonName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.all })
    },
  })
}

export function useImportLicenseFile() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (fileData: number[]) => licenseService.importLicenseFile(fileData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.all })
    },
  })
}

export function useRenewLicense() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (licenseKey?: string) => licenseService.renewLicense(licenseKey),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.all })
    },
  })
}

export function useDismissNotification() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => licenseService.dismissNotification(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.notifications })
    },
  })
}

export function useMarkNotificationRead() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => licenseService.markNotificationRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LICENSE_KEYS.notifications })
    },
  })
}
