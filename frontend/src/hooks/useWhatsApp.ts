import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as api from '@/services/whatsapp'
import type { WhatsAppTemplate, AutomationRule } from '@/types'

export function useWhatsAppTemplates() {
  return useQuery({ queryKey: ['wa-templates'], queryFn: api.listTemplates })
}

export function useCreateTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<WhatsAppTemplate>) => api.createTemplate(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-templates'] }),
  })
}

export function useUpdateTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: Partial<WhatsAppTemplate> & { id: string }) => api.updateTemplate(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-templates'] }),
  })
}

export function useDeleteTemplate() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.deleteTemplate(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-templates'] }),
  })
}

export function useSendMessage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { template_id: string; customer_id: string; phone: string; variables?: Record<string, string> }) => api.sendMessage(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-messages'] }),
  })
}

export function useWhatsAppMessages(page = 1) {
  return useQuery({ queryKey: ['wa-messages', page], queryFn: () => api.listMessages(page) })
}

export function useWAMessageStats() {
  return useQuery({ queryKey: ['wa-stats'], queryFn: api.getMessageStats })
}

export function useAutomationRules() {
  return useQuery({ queryKey: ['wa-rules'], queryFn: api.listRules })
}

export function useCreateRule() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<AutomationRule>) => api.createRule(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-rules'] }),
  })
}

export function useUpdateRule() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: Partial<AutomationRule> & { id: string }) => api.updateRule(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-rules'] }),
  })
}

export function useDeleteRule() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.deleteRule(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['wa-rules'] }),
  })
}
