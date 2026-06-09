import { useState } from 'react'
import { usePrinterSettings, useSavePrinterSettings, usePrintTest, usePrintJobs } from '@/hooks/usePrinter'
import type { PrinterSettings } from '@/types'

export function PrinterPage() {
  const [tab, setTab] = useState<'settings' | 'history' | 'preview'>('settings')

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Printer & Receipt</h1>
        <p className="text-muted-foreground">Configure printers, manage receipt templates, and view print history</p>
      </div>

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('settings')} className={`px-4 py-2 rounded-t ${tab === 'settings' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Settings</button>
        <button onClick={() => setTab('history')} className={`px-4 py-2 rounded-t ${tab === 'history' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Print History</button>
        <button onClick={() => setTab('preview')} className={`px-4 py-2 rounded-t ${tab === 'preview' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Preview</button>
      </div>

      {tab === 'settings' && <PrinterSettingsTab />}
      {tab === 'history' && <PrintHistoryTab />}
      {tab === 'preview' && <PrintPreviewTab />}
    </div>
  )
}

function PrinterSettingsTab() {
  const { data: settings, isLoading } = usePrinterSettings()
  const saveMutation = useSavePrinterSettings()
  const testMutation = usePrintTest()
  const [form, setForm] = useState<Partial<PrinterSettings>>({})

  const handleSave = () => {
    saveMutation.mutate({ ...settings, ...form })
  }

  if (isLoading) return <p>Loading...</p>

  const s = { ...settings, ...form }

  return (
    <div className="space-y-4 max-w-2xl">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-sm font-medium">Default Printer</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.default_printer || ''} onChange={e => setForm({ ...form, default_printer: e.target.value })} placeholder="Printer name" />
        </div>
        <div>
          <label className="text-sm font-medium">Paper Width</label>
          <select className="w-full border rounded px-3 py-2 mt-1" value={s.paper_width || '80mm'} onChange={e => setForm({ ...form, paper_width: e.target.value })}>
            <option value="58mm">58mm (Thermal)</option>
            <option value="80mm">80mm (Thermal)</option>
            <option value="A4">A4 (Full Page)</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Header Text</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.header_text || ''} onChange={e => setForm({ ...form, header_text: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">Footer Text</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.footer_text || ''} onChange={e => setForm({ ...form, footer_text: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">UPI ID (for QR)</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.upi_id || ''} onChange={e => setForm({ ...form, upi_id: e.target.value })} placeholder="example@upi" />
        </div>
        <div className="space-y-2 pt-4">
          <label className="flex items-center gap-2">
            <input type="checkbox" checked={s.show_logo ?? false} onChange={e => setForm({ ...form, show_logo: e.target.checked })} />
            <span className="text-sm">Show Logo</span>
          </label>
          <label className="flex items-center gap-2">
            <input type="checkbox" checked={s.show_qr ?? false} onChange={e => setForm({ ...form, show_qr: e.target.checked })} />
            <span className="text-sm">Show QR Code</span>
          </label>
        </div>
      </div>

      <h3 className="font-semibold">Margins (mm)</h3>
      <div className="grid grid-cols-4 gap-3">
        <div>
          <label className="text-xs">Top</label>
          <input type="number" className="w-full border rounded px-2 py-1" value={s.margin_top ?? 5} onChange={e => setForm({ ...form, margin_top: +e.target.value })} />
        </div>
        <div>
          <label className="text-xs">Bottom</label>
          <input type="number" className="w-full border rounded px-2 py-1" value={s.margin_bottom ?? 5} onChange={e => setForm({ ...form, margin_bottom: +e.target.value })} />
        </div>
        <div>
          <label className="text-xs">Left</label>
          <input type="number" className="w-full border rounded px-2 py-1" value={s.margin_left ?? 5} onChange={e => setForm({ ...form, margin_left: +e.target.value })} />
        </div>
        <div>
          <label className="text-xs">Right</label>
          <input type="number" className="w-full border rounded px-2 py-1" value={s.margin_right ?? 5} onChange={e => setForm({ ...form, margin_right: +e.target.value })} />
        </div>
      </div>

      <div className="flex gap-3">
        <button onClick={handleSave} disabled={saveMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
          {saveMutation.isPending ? 'Saving...' : 'Save Settings'}
        </button>
        <button onClick={() => testMutation.mutate()} disabled={testMutation.isPending} className="px-4 py-2 bg-muted rounded hover:bg-muted/80">
          {testMutation.isPending ? 'Printing...' : 'Print Test Page'}
        </button>
      </div>
      {testMutation.isSuccess && <p className="text-green-600 text-sm">Test page sent to printer</p>}
    </div>
  )
}

function PrintHistoryTab() {
  const [page, setPage] = useState(1)
  const { data, isLoading } = usePrintJobs(page)

  if (isLoading) return <p>Loading...</p>

  const jobs = data?.jobs || []
  const meta = data?.meta

  return (
    <div className="space-y-4">
      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b">
            <th className="text-left p-2">Type</th>
            <th className="text-left p-2">Document</th>
            <th className="text-left p-2">Printer</th>
            <th className="text-left p-2">Paper</th>
            <th className="text-left p-2">Status</th>
            <th className="text-left p-2">Date</th>
          </tr>
        </thead>
        <tbody>
          {jobs.map(job => (
            <tr key={job.id} className="border-b">
              <td className="p-2 capitalize">{job.document_type}</td>
              <td className="p-2">{job.document_id}</td>
              <td className="p-2">{job.printer_name || '—'}</td>
              <td className="p-2">{job.paper_width}</td>
              <td className="p-2">
                <span className={`px-2 py-0.5 rounded text-xs ${job.status === 'completed' ? 'bg-green-100 text-green-700' : job.status === 'failed' ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'}`}>
                  {job.status}
                </span>
              </td>
              <td className="p-2 text-sm">{new Date(job.created_at).toLocaleString()}</td>
            </tr>
          ))}
          {jobs.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-muted-foreground">No print history</td></tr>}
        </tbody>
      </table>

      {meta && meta.total_pages > 1 && (
        <div className="flex gap-2 justify-center">
          <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="px-3 py-1 border rounded disabled:opacity-50">Prev</button>
          <span className="px-3 py-1">Page {page} of {meta.total_pages}</span>
          <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)} className="px-3 py-1 border rounded disabled:opacity-50">Next</button>
        </div>
      )}
    </div>
  )
}

function PrintPreviewTab() {
  const sampleReceipt = `
            Glamour Salon
     123 Main St, Mumbai
      GSTIN: 27AABCU9603R1ZM
------------------------------------------------
Invoice: INV-2026-0001
Date: 2026-06-09
Customer: Priya Sharma
Phone: 9876543210
------------------------------------------------
Item                    Qty   Price   Total
------------------------------------------------
Haircut                  1     500      500
Hair Color               1    2000     2000
Head Massage             1     300      300
------------------------------------------------
                      Subtotal:       2800.00
                      CGST 9%:        252.00
                      SGST 9%:        252.00
================================================
                      GRAND TOTAL:   3304.00
================================================
Paid by: UPI

         Thank you for visiting!
  `.trim()

  return (
    <div className="space-y-4">
      <h3 className="font-semibold">Receipt Preview (80mm)</h3>
      <div className="bg-white border rounded p-4 max-w-md mx-auto font-mono text-xs whitespace-pre leading-tight">
        {sampleReceipt}
      </div>
    </div>
  )
}
