import { describe, it, expect, vi } from 'vitest'
import { createTemplate, updateTemplate, deleteTemplate, listTemplates, sendMessage, listMessages, getMessageStats, createRule, updateRule, deleteRule, listRules } from './whatsapp'

describe('WhatsApp Service', () => {
  it('creates template', async () => {
    await expect(createTemplate({ name: 'Welcome', body: 'Hello!' })).resolves.toBeUndefined()
  })

  it('updates template', async () => {
    await expect(updateTemplate('tpl1', { name: 'Welcome Updated' })).resolves.toBeUndefined()
  })

  it('deletes template', async () => {
    await expect(deleteTemplate('tpl1')).resolves.toBeUndefined()
  })

  it('lists templates', async () => {
    const templates = await listTemplates()
    expect(templates).toHaveLength(1)
    expect(templates[0].name).toBe('Welcome')
  })

  it('sends message', async () => {
    const msg = await sendMessage({ template_id: 'tpl1', phone: '9876543210' })
    expect(msg.id).toBe('msg1')
  })

  it('sends message with variables', async () => {
    const msg = await sendMessage({ template_id: 'tpl1', phone: '9876543210', variables: { name: 'Anjali' } })
    expect(msg.id).toBe('msg1')
  })

  it('lists messages', async () => {
    const result = await listMessages()
    expect(result.data).toHaveLength(1)
    expect(result.meta.total).toBe(1)
  })

  it('lists messages with pagination', async () => {
    const result = await listMessages(2, 10)
    expect(result.data).toHaveLength(1)
  })

  it('gets message stats', async () => {
    const stats = await getMessageStats()
    expect(stats.total_sent).toBe(100)
  })

  it('creates rule', async () => {
    await expect(createRule({ name: 'Birthday', trigger: 'birthday' })).resolves.toBeUndefined()
  })

  it('updates rule', async () => {
    await expect(updateRule('rule1', { name: 'Birthday Updated' })).resolves.toBeUndefined()
  })

  it('deletes rule', async () => {
    await expect(deleteRule('rule1')).resolves.toBeUndefined()
  })

  it('lists rules', async () => {
    const rules = await listRules()
    expect(rules).toHaveLength(1)
  })

  it('listMessages returns empty when API returns undefined', async () => {
    vi.spyOn(window.go.main.WhatsAppService, 'ListMessages').mockResolvedValueOnce([undefined as any, 0])
    const r = await listMessages()
    expect(r.data).toEqual([])
  })
})
