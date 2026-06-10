import type { WhatsAppTemplate, WhatsAppMessage, WAMessageStats, AutomationRule } from '@/types'

export async function createTemplate(data: any): Promise<void> {
  await window.go.main.WhatsAppService.CreateTemplate(data)
}

export async function updateTemplate(id: string, data: any): Promise<void> {
  await window.go.main.WhatsAppService.UpdateTemplate({ ...data, id })
}

export async function deleteTemplate(id: string): Promise<void> {
  await window.go.main.WhatsAppService.DeleteTemplate(id)
}

export async function listTemplates(): Promise<WhatsAppTemplate[]> {
  return window.go.main.WhatsAppService.ListTemplates('')
}

export async function sendMessage(data: { template_id: string; customer_id?: string; phone: string; variables?: Record<string, string> }): Promise<WhatsAppMessage> {
  return window.go.main.WhatsAppService.SendMessage(data.template_id, data.phone, '', data.variables || {})
}

export async function listMessages(page = 1, perPage = 20) {
  const offset = (page - 1) * perPage
  const [messages, total] = await window.go.main.WhatsAppService.ListMessages(perPage, offset, '')
  return { data: messages || [], meta: { page, per_page: perPage, total, total_pages: Math.ceil(total / perPage) } }
}

export async function getMessageStats(): Promise<WAMessageStats> {
  return window.go.main.WhatsAppService.GetWhatsAppStats()
}

export async function createRule(data: any): Promise<void> {
  await window.go.main.WhatsAppService.CreateRule(data)
}

export async function updateRule(id: string, data: any): Promise<void> {
  await window.go.main.WhatsAppService.UpdateRule({ ...data, id })
}

export async function deleteRule(id: string): Promise<void> {
  await window.go.main.WhatsAppService.DeleteRule(id)
}

export async function listRules(): Promise<AutomationRule[]> {
  return window.go.main.WhatsAppService.ListRules()
}
