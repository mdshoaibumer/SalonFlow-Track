import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getGSTSettings, saveGSTSettings, listTaxRates, createTaxRate, updateTaxRate, deleteTaxRate, getGSTReport } from '@/services/gst'
import type { GSTSettings, TaxRate } from '@/types'

export function useGSTSettings() {
  return useQuery({ queryKey: ['gst-settings'], queryFn: getGSTSettings })
}

export function useSaveGSTSettings() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (settings: Partial<GSTSettings>) => saveGSTSettings(settings),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['gst-settings'] }),
  })
}

export function useTaxRates(category?: string) {
  return useQuery({ queryKey: ['tax-rates', category], queryFn: () => listTaxRates(category) })
}

export function useCreateTaxRate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (rate: Partial<TaxRate>) => createTaxRate(rate),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tax-rates'] }),
  })
}

export function useUpdateTaxRate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, rate }: { id: string; rate: Partial<TaxRate> }) => updateTaxRate(id, rate),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tax-rates'] }),
  })
}

export function useDeleteTaxRate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteTaxRate(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tax-rates'] }),
  })
}

export function useGSTReport(startDate: string, endDate: string, period = 'daily') {
  return useQuery({
    queryKey: ['gst-report', startDate, endDate, period],
    queryFn: () => getGSTReport(startDate, endDate, period),
    enabled: !!startDate && !!endDate,
  })
}
