import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getPrinterSettings, savePrinterSettings, printInvoice, printReceipt, printTest, listPrintJobs } from '@/services/printer'
import type { PrinterSettings, ReceiptData } from '@/types'

export function usePrinterSettings() {
  return useQuery({ queryKey: ['printer-settings'], queryFn: getPrinterSettings })
}

export function useSavePrinterSettings() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (settings: Partial<PrinterSettings>) => savePrinterSettings(settings),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['printer-settings'] }),
  })
}

export function usePrintInvoice() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: ReceiptData) => printInvoice(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['print-jobs'] }),
  })
}

export function usePrintReceipt() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: ReceiptData) => printReceipt(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['print-jobs'] }),
  })
}

export function usePrintTest() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => printTest(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['print-jobs'] }),
  })
}

export function usePrintJobs(page = 1) {
  return useQuery({ queryKey: ['print-jobs', page], queryFn: () => listPrintJobs(page) })
}
