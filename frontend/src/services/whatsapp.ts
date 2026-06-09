import { apiClient } from './api-client'
import type { WhatsAppTemplate, WhatsAppMessage, AutomationRule, WAMessageStats } from '@/types'

export async function listTemplates(): Promise<WhatsAppTemplate[]> {
  const response = await apiClient.get<WhatsAppTemplate[]>('/whatsapp/templates')
  if (!response.success) throw new Error(response.error?.message || 'Failed to list templates')
  return response.data || []
}

export async function createTemplate(data: Partial<WhatsAppTemplate>): Promise<WhatsAppTemplate> {
  const response = await apiClient.post<WhatsAppTemplate>('/whatsapp/templates', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create template')
  return response.data
}

export async function updateTemplate(id: string, data: Partial<WhatsAppTemplate>): Promise<WhatsAppTemplate> {
  const response = await apiClient.put<WhatsAppTemplate>(`/whatsapp/templates/${id}`, data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update template')
  return response.data
}

export async function deleteTemplate(id: string): Promise<void> {
  const response = await apiClient.delete(`/whatsapp/templates/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete template')
}

export async function sendMessage(data: { template_id: string; customer_id: string; phone: string; variables?: Record<string, string> }): Promise<WhatsAppMessage> {
  const response = await apiClient.post<WhatsAppMessage>('/whatsapp/send', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to send message')
  return response.data
}

export async function listMessages(page = 1, perPage = 20): Promise<{ data: WhatsAppMessage[]; total: number }> {
  const response = await apiClient.get<WhatsAppMessage[]>(`/whatsapp/messages?page=${page}&per_page=${perPage}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to list messages')
  return { data: response.data || [], total: response.meta?.total || 0 }
}

export async function getMessageStats(): Promise<WAMessageStats> {
  const response = await apiClient.get<WAMessageStats>('/whatsapp/stats')
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to get stats')
  return response.data
}

export async function listRules(): Promise<AutomationRule[]> {
  const response = await apiClient.get<AutomationRule[]>('/whatsapp/rules')
  if (!response.success) throw new Error(response.error?.message || 'Failed to list rules')
  return response.data || []
}

export async function createRule(data: Partial<AutomationRule>): Promise<AutomationRule> {
  const response = await apiClient.post<AutomationRule>('/whatsapp/rules', data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to create rule')
  return response.data
}

export async function updateRule(id: string, data: Partial<AutomationRule>): Promise<AutomationRule> {
  const response = await apiClient.put<AutomationRule>(`/whatsapp/rules/${id}`, data)
  if (!response.success || !response.data) throw new Error(response.error?.message || 'Failed to update rule')
  return response.data
}

export async function deleteRule(id: string): Promise<void> {
  const response = await apiClient.delete(`/whatsapp/rules/${id}`)
  if (!response.success) throw new Error(response.error?.message || 'Failed to delete rule')
}
