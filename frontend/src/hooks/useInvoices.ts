import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listInvoices, getInvoiceById, createInvoice, recordPayment, getInvoiceStats, type ListInvoiceParams } from '@/services/invoices'
import type { CreateInvoiceInput, RecordPaymentInput } from '@/types'

export function useInvoiceList(params: ListInvoiceParams = {}) {
  return useQuery({
    queryKey: ['invoices', params],
    queryFn: () => listInvoices(params),
  })
}

export function useInvoiceById(id: string) {
  return useQuery({
    queryKey: ['invoices', id],
    queryFn: () => getInvoiceById(id),
    enabled: !!id,
  })
}

export function useInvoiceStats() {
  return useQuery({
    queryKey: ['invoices', 'stats'],
    queryFn: getInvoiceStats,
  })
}

export function useCreateInvoice() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateInvoiceInput) => createInvoice(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] })
      queryClient.invalidateQueries({ queryKey: ['customers'] })
    },
  })
}

export function useRecordPayment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ invoiceId, input }: { invoiceId: string; input: RecordPaymentInput }) => recordPayment(invoiceId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] })
    },
  })
}
