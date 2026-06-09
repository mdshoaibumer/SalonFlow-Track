import { useState } from 'react'
import { useGSTSettings, useSaveGSTSettings, useTaxRates, useCreateTaxRate, useDeleteTaxRate, useGSTReport } from '@/hooks/useGST'
import type { GSTSettings } from '@/types'

export function GSTPage() {
  const [tab, setTab] = useState<'settings' | 'rates' | 'reports'>('settings')
  const [reportStart, setReportStart] = useState('')
  const [reportEnd, setReportEnd] = useState('')

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">GST & Tax Management</h1>
        <p className="text-muted-foreground">Configure GST settings, manage tax rates, and view reports</p>
      </div>

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('settings')} className={`px-4 py-2 rounded-t ${tab === 'settings' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Settings</button>
        <button onClick={() => setTab('rates')} className={`px-4 py-2 rounded-t ${tab === 'rates' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Tax Rates</button>
        <button onClick={() => setTab('reports')} className={`px-4 py-2 rounded-t ${tab === 'reports' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Reports</button>
      </div>

      {tab === 'settings' && <GSTSettingsTab />}
      {tab === 'rates' && <TaxRatesTab />}
      {tab === 'reports' && <GSTReportsTab startDate={reportStart} endDate={reportEnd} onStartChange={setReportStart} onEndChange={setReportEnd} />}
    </div>
  )
}

function GSTSettingsTab() {
  const { data: settings, isLoading } = useGSTSettings()
  const saveMutation = useSaveGSTSettings()
  const [form, setForm] = useState<Partial<GSTSettings>>({})

  const handleSave = () => {
    saveMutation.mutate({ ...settings, ...form })
  }

  if (isLoading) return (
    <div className="space-y-4 max-w-2xl">
      <div className="flex items-center gap-2 pt-6">
        <input type="checkbox" id="gst-enabled" disabled />
        <label htmlFor="gst-enabled" className="text-sm font-medium">Enable GST Billing</label>
      </div>
      <p className="text-muted-foreground">Loading GST settings...</p>
    </div>
  )

  const s = { ...settings, ...form }

  return (
    <div className="space-y-4 max-w-2xl">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-sm font-medium">Business Name</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.business_name || ''} onChange={e => setForm({ ...form, business_name: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">GSTIN</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.gstin || ''} onChange={e => setForm({ ...form, gstin: e.target.value })} placeholder="22AAAAA0000A1Z5" />
        </div>
        <div>
          <label className="text-sm font-medium">State</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.state || ''} onChange={e => setForm({ ...form, state: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">HSN Code</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.hsn_code || ''} onChange={e => setForm({ ...form, hsn_code: e.target.value })} />
        </div>
        <div className="col-span-2">
          <label className="text-sm font-medium">Address</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={s.address || ''} onChange={e => setForm({ ...form, address: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">CGST Rate (%)</label>
          <input type="number" step="0.5" className="w-full border rounded px-3 py-2 mt-1" value={s.cgst_rate ?? 9} onChange={e => setForm({ ...form, cgst_rate: parseFloat(e.target.value) })} />
        </div>
        <div>
          <label className="text-sm font-medium">SGST Rate (%)</label>
          <input type="number" step="0.5" className="w-full border rounded px-3 py-2 mt-1" value={s.sgst_rate ?? 9} onChange={e => setForm({ ...form, sgst_rate: parseFloat(e.target.value) })} />
        </div>
        <div>
          <label className="text-sm font-medium">IGST Rate (%)</label>
          <input type="number" step="0.5" className="w-full border rounded px-3 py-2 mt-1" value={s.igst_rate ?? 18} onChange={e => setForm({ ...form, igst_rate: parseFloat(e.target.value) })} />
        </div>
        <div className="flex items-center gap-2 pt-6">
          <input type="checkbox" id="gst-enabled" checked={s.is_gst_enabled ?? false} onChange={e => setForm({ ...form, is_gst_enabled: e.target.checked })} />
          <label htmlFor="gst-enabled" className="text-sm font-medium">Enable GST Billing</label>
        </div>
      </div>
      <button onClick={handleSave} disabled={saveMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
        {saveMutation.isPending ? 'Saving...' : 'Save Settings'}
      </button>
    </div>
  )
}

function TaxRatesTab() {
  const { data: rates = [], isLoading } = useTaxRates()
  const createMutation = useCreateTaxRate()
  const deleteMutation = useDeleteTaxRate()
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [hsn, setHsn] = useState('')
  const [category, setCategory] = useState('')
  const [cgst, setCgst] = useState(9)
  const [sgst, setSgst] = useState(9)
  const [igst, setIgst] = useState(18)

  const handleCreate = () => {
    createMutation.mutate({ name, hsn_code: hsn, category, cgst_rate: cgst, sgst_rate: sgst, igst_rate: igst }, {
      onSuccess: () => { setShowForm(false); setName(''); setHsn(''); setCategory('') },
    })
  }

  if (isLoading) return <p>Loading...</p>

  return (
    <div className="space-y-4">
      <button onClick={() => setShowForm(!showForm)} className="px-4 py-2 bg-primary text-primary-foreground rounded">
        {showForm ? 'Cancel' : '+ Add Tax Rate'}
      </button>

      {showForm && (
        <div className="grid grid-cols-3 gap-3 p-4 border rounded">
          <input placeholder="Name" value={name} onChange={e => setName(e.target.value)} className="border rounded px-3 py-2" />
          <input placeholder="HSN Code" value={hsn} onChange={e => setHsn(e.target.value)} className="border rounded px-3 py-2" />
          <input placeholder="Category" value={category} onChange={e => setCategory(e.target.value)} className="border rounded px-3 py-2" />
          <input type="number" placeholder="CGST %" value={cgst} onChange={e => setCgst(+e.target.value)} className="border rounded px-3 py-2" />
          <input type="number" placeholder="SGST %" value={sgst} onChange={e => setSgst(+e.target.value)} className="border rounded px-3 py-2" />
          <input type="number" placeholder="IGST %" value={igst} onChange={e => setIgst(+e.target.value)} className="border rounded px-3 py-2" />
          <button onClick={handleCreate} className="px-4 py-2 bg-green-600 text-white rounded col-span-3">Create</button>
        </div>
      )}

      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b">
            <th className="text-left p-2">Name</th>
            <th className="text-left p-2">HSN</th>
            <th className="text-left p-2">Category</th>
            <th className="text-right p-2">CGST %</th>
            <th className="text-right p-2">SGST %</th>
            <th className="text-right p-2">IGST %</th>
            <th className="text-right p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {rates.map(rate => (
            <tr key={rate.id} className="border-b">
              <td className="p-2">{rate.name}</td>
              <td className="p-2">{rate.hsn_code}</td>
              <td className="p-2">{rate.category}</td>
              <td className="p-2 text-right">{rate.cgst_rate}</td>
              <td className="p-2 text-right">{rate.sgst_rate}</td>
              <td className="p-2 text-right">{rate.igst_rate}</td>
              <td className="p-2 text-right">
                <button onClick={() => deleteMutation.mutate(rate.id)} className="text-red-600 hover:underline text-sm">Delete</button>
              </td>
            </tr>
          ))}
          {rates.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-muted-foreground">No tax rates configured</td></tr>}
        </tbody>
      </table>
    </div>
  )
}

function GSTReportsTab({ startDate, endDate, onStartChange, onEndChange }: { startDate: string; endDate: string; onStartChange: (v: string) => void; onEndChange: (v: string) => void }) {
  const { data: report, isLoading } = useGSTReport(startDate, endDate)

  return (
    <div className="space-y-4 max-w-2xl">
      <div className="flex gap-4 items-end">
        <div>
          <label className="text-sm font-medium">Start Date</label>
          <input type="date" className="border rounded px-3 py-2 mt-1 block" value={startDate} onChange={e => onStartChange(e.target.value)} />
        </div>
        <div>
          <label className="text-sm font-medium">End Date</label>
          <input type="date" className="border rounded px-3 py-2 mt-1 block" value={endDate} onChange={e => onEndChange(e.target.value)} />
        </div>
      </div>

      {isLoading && <p>Loading report...</p>}

      {report && (
        <div className="border rounded p-4 space-y-3">
          <h3 className="font-semibold">GST Summary: {report.start_date} to {report.end_date}</h3>
          <div className="grid grid-cols-2 gap-3">
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">Total Invoices</p>
              <p className="text-xl font-bold">{report.total_invoices}</p>
            </div>
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">Taxable Amount</p>
              <p className="text-xl font-bold">₹{report.taxable_amount.toFixed(2)}</p>
            </div>
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">CGST Collected</p>
              <p className="text-xl font-bold">₹{report.total_cgst.toFixed(2)}</p>
            </div>
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">SGST Collected</p>
              <p className="text-xl font-bold">₹{report.total_sgst.toFixed(2)}</p>
            </div>
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">IGST Collected</p>
              <p className="text-xl font-bold">₹{report.total_igst.toFixed(2)}</p>
            </div>
            <div className="p-3 bg-muted rounded">
              <p className="text-sm text-muted-foreground">Total Tax</p>
              <p className="text-xl font-bold">₹{report.total_tax.toFixed(2)}</p>
            </div>
            <div className="p-3 bg-primary/10 rounded col-span-2">
              <p className="text-sm text-muted-foreground">Grand Total</p>
              <p className="text-2xl font-bold">₹{report.grand_total.toFixed(2)}</p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
