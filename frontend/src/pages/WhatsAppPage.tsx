import { useState } from 'react'
import { useWhatsAppTemplates, useCreateTemplate, useDeleteTemplate, useSendMessage, useWhatsAppMessages, useWAMessageStats, useAutomationRules, useCreateRule, useDeleteRule } from '@/hooks/useWhatsApp'
import type { WhatsAppTemplate, WhatsAppMessage, AutomationRule } from '@/types'

export function WhatsAppPage() {
  const [tab, setTab] = useState<'templates' | 'messages' | 'rules' | 'send'>('templates')

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">WhatsApp Automation</h1>
        <p className="text-muted-foreground">Manage templates, send messages, and configure automation rules</p>
      </div>

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('templates')} className={`px-4 py-2 rounded-t ${tab === 'templates' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Templates</button>
        <button onClick={() => setTab('send')} className={`px-4 py-2 rounded-t ${tab === 'send' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Send Message</button>
        <button onClick={() => setTab('messages')} className={`px-4 py-2 rounded-t ${tab === 'messages' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Message Log</button>
        <button onClick={() => setTab('rules')} className={`px-4 py-2 rounded-t ${tab === 'rules' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Automation Rules</button>
      </div>

      {tab === 'templates' && <TemplatesTab />}
      {tab === 'send' && <SendMessageTab />}
      {tab === 'messages' && <MessagesTab />}
      {tab === 'rules' && <RulesTab />}
    </div>
  )
}

function TemplatesTab() {
  const { data: templates, isLoading } = useWhatsAppTemplates()
  const createMutation = useCreateTemplate()
  const deleteMutation = useDeleteTemplate()
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', category: 'appointment', body: '' })

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(form, { onSuccess: () => { setShowForm(false); setForm({ name: '', category: 'appointment', body: '' }) } })
  }

  return (
    <div className="space-y-4">
      <button onClick={() => setShowForm(!showForm)} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">
        {showForm ? 'Cancel' : 'New Template'}
      </button>

      {showForm && (
        <form onSubmit={handleCreate} className="border rounded p-4 space-y-3 max-w-lg">
          <div>
            <label className="text-sm font-medium">Name</label>
            <input className="w-full border rounded px-3 py-2 mt-1" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="text-sm font-medium">Category</label>
            <select className="w-full border rounded px-3 py-2 mt-1" value={form.category} onChange={e => setForm({ ...form, category: e.target.value })}>
              <option value="appointment">Appointment</option>
              <option value="promotion">Promotion</option>
              <option value="reminder">Reminder</option>
              <option value="feedback">Feedback</option>
            </select>
          </div>
          <div>
            <label className="text-sm font-medium">Body (use {'{{variable}}'} for placeholders)</label>
            <textarea className="w-full border rounded px-3 py-2 mt-1" value={form.body} onChange={e => setForm({ ...form, body: e.target.value })} rows={4} required />
          </div>
          <button type="submit" disabled={createMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Create</button>
        </form>
      )}

      {isLoading && <p>Loading templates...</p>}

      <div className="grid gap-3">
        {(templates || []).map((t: WhatsAppTemplate) => (
          <div key={t.id} className="border rounded p-4 flex justify-between items-start">
            <div>
              <h3 className="font-medium">{t.name}</h3>
              <p className="text-xs text-muted-foreground mt-1">{t.category} &middot; {t.is_active ? 'Active' : 'Inactive'}</p>
              <p className="text-sm mt-2 bg-muted p-2 rounded">{t.body}</p>
            </div>
            <button onClick={() => deleteMutation.mutate(t.id)} className="text-xs px-2 py-1 bg-red-100 text-red-800 rounded">Delete</button>
          </div>
        ))}
      </div>
    </div>
  )
}

function SendMessageTab() {
  const { data: templates } = useWhatsAppTemplates()
  const sendMutation = useSendMessage()
  const [form, setForm] = useState({ template_id: '', customer_id: '', phone: '' })

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault()
    sendMutation.mutate(form)
  }

  return (
    <form onSubmit={handleSend} className="space-y-4 max-w-lg">
      <div>
        <label className="text-sm font-medium">Template</label>
        <select className="w-full border rounded px-3 py-2 mt-1" value={form.template_id} onChange={e => setForm({ ...form, template_id: e.target.value })} required>
          <option value="">Select template...</option>
          {(templates || []).map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
        </select>
      </div>
      <div>
        <label className="text-sm font-medium">Customer ID</label>
        <input className="w-full border rounded px-3 py-2 mt-1" value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} required />
      </div>
      <div>
        <label className="text-sm font-medium">Phone Number</label>
        <input className="w-full border rounded px-3 py-2 mt-1" value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} placeholder="+91..." required />
      </div>
      <button type="submit" disabled={sendMutation.isPending} className="px-4 py-2 bg-green-600 text-white rounded">
        {sendMutation.isPending ? 'Sending...' : 'Send Message'}
      </button>
      {sendMutation.isSuccess && <p className="text-green-600 text-sm">Message sent!</p>}
      {sendMutation.isError && <p className="text-red-600 text-sm">{(sendMutation.error as Error).message}</p>}
    </form>
  )
}

function MessagesTab() {
  const { data: stats } = useWAMessageStats()
  const { data: messages, isLoading } = useWhatsAppMessages()

  if (isLoading) return <p>Loading...</p>

  return (
    <div className="space-y-4">
      {stats && (
        <div className="grid grid-cols-4 gap-4">
          <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.total_sent}</p><p className="text-xs text-muted-foreground">Sent</p></div>
          <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.total_delivered}</p><p className="text-xs text-muted-foreground">Delivered</p></div>
          <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.total_failed}</p><p className="text-xs text-muted-foreground">Failed</p></div>
          <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.total_pending}</p><p className="text-xs text-muted-foreground">Pending</p></div>
        </div>
      )}
      <div className="border rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-muted"><tr><th className="px-4 py-2 text-left">Phone</th><th className="px-4 py-2 text-left">Body</th><th className="px-4 py-2 text-left">Status</th><th className="px-4 py-2 text-left">Sent</th></tr></thead>
          <tbody>
            {(messages?.data || []).map((m: WhatsAppMessage) => (
              <tr key={m.id} className="border-t">
                <td className="px-4 py-2">{m.phone}</td>
                <td className="px-4 py-2 max-w-xs truncate">{m.body}</td>
                <td className="px-4 py-2"><span className={`px-2 py-0.5 rounded text-xs ${m.status === 'delivered' ? 'bg-green-100' : m.status === 'failed' ? 'bg-red-100' : 'bg-yellow-100'}`}>{m.status}</span></td>
                <td className="px-4 py-2">{m.sent_at ? new Date(m.sent_at).toLocaleString() : '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function RulesTab() {
  const { data: rules, isLoading } = useAutomationRules()
  const { data: templates } = useWhatsAppTemplates()
  const createMutation = useCreateRule()
  const deleteMutation = useDeleteRule()
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', trigger: 'appointment_booked', template_id: '', delay_minutes: 0, is_active: true })

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(form, { onSuccess: () => setShowForm(false) })
  }

  if (isLoading) return <p>Loading...</p>

  return (
    <div className="space-y-4">
      <button onClick={() => setShowForm(!showForm)} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">
        {showForm ? 'Cancel' : 'New Rule'}
      </button>

      {showForm && (
        <form onSubmit={handleCreate} className="border rounded p-4 space-y-3 max-w-lg">
          <div>
            <label className="text-sm font-medium">Rule Name</label>
            <input className="w-full border rounded px-3 py-2 mt-1" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="text-sm font-medium">Trigger</label>
            <select className="w-full border rounded px-3 py-2 mt-1" value={form.trigger} onChange={e => setForm({ ...form, trigger: e.target.value })}>
              <option value="appointment_booked">Appointment Booked</option>
              <option value="appointment_reminder">Appointment Reminder</option>
              <option value="invoice_created">Invoice Created</option>
              <option value="birthday">Birthday</option>
              <option value="membership_expiry">Membership Expiry</option>
            </select>
          </div>
          <div>
            <label className="text-sm font-medium">Template</label>
            <select className="w-full border rounded px-3 py-2 mt-1" value={form.template_id} onChange={e => setForm({ ...form, template_id: e.target.value })} required>
              <option value="">Select template...</option>
              {(templates || []).map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
            </select>
          </div>
          <div>
            <label className="text-sm font-medium">Delay (minutes)</label>
            <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={form.delay_minutes} onChange={e => setForm({ ...form, delay_minutes: parseInt(e.target.value) || 0 })} />
          </div>
          <button type="submit" disabled={createMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Create Rule</button>
        </form>
      )}

      <div className="grid gap-3">
        {(rules || []).map((r: AutomationRule) => (
          <div key={r.id} className="border rounded p-4 flex justify-between items-center">
            <div>
              <h3 className="font-medium">{r.name}</h3>
              <p className="text-xs text-muted-foreground">Trigger: {r.trigger} &middot; Delay: {r.delay_minutes}min &middot; {r.is_active ? 'Active' : 'Inactive'}</p>
            </div>
            <button onClick={() => deleteMutation.mutate(r.id)} className="text-xs px-2 py-1 bg-red-100 text-red-800 rounded">Delete</button>
          </div>
        ))}
      </div>
    </div>
  )
}
