import { useState } from 'react'
import { useSalaryList, useSalaryStats, useGenerateSalary, usePaySalary, useSalaryCycles } from '@/hooks/useSalary'
import { IndianRupee, Clock, CheckCircle, AlertCircle } from 'lucide-react'

const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']

export function SalaryPage() {
  const now = new Date()
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [year, setYear] = useState(now.getFullYear())

  const { data: stats } = useSalaryStats()
  const { data: records, isLoading } = useSalaryList(month, year)
  const { data: cycles } = useSalaryCycles(year)
  const generateSalary = useGenerateSalary()
  const paySalary = usePaySalary()

  const currentCycle = cycles?.find(c => c.month === month && c.year === year)

  const handleGenerate = () => {
    if (confirm(`Generate salary for ${MONTHS[month - 1]} ${year}?`)) {
      generateSalary.mutate({ month, year })
    }
  }

  const handlePay = (id: string, name: string) => {
    if (confirm(`Mark salary as paid for ${name}?`)) {
      paySalary.mutate(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Salary Management</h1>
          <p className="text-muted-foreground">Generate and manage monthly payroll</p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={month}
            onChange={(e) => setMonth(Number(e.target.value))}
            className="rounded-md border px-3 py-2 text-sm"
          >
            {MONTHS.map((m, i) => (
              <option key={i} value={i + 1}>{m}</option>
            ))}
          </select>
          <select
            value={year}
            onChange={(e) => setYear(Number(e.target.value))}
            className="rounded-md border px-3 py-2 text-sm"
          >
            {[2024, 2025, 2026, 2027].map(y => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
          {!currentCycle && (
            <button
              onClick={handleGenerate}
              disabled={generateSalary.isPending}
              className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              {generateSalary.isPending ? 'Generating...' : 'Generate Salary'}
            </button>
          )}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <StatCard title="Total Payroll" value={`₹${Math.round(stats?.total_payroll ?? 0).toLocaleString()}`} icon={<IndianRupee className="h-4 w-4" />} />
        <StatCard title="Pending" value={String(stats?.pending_payments ?? 0)} icon={<Clock className="h-4 w-4" />} />
        <StatCard title="Paid" value={String(stats?.paid_salaries ?? 0)} icon={<CheckCircle className="h-4 w-4" />} />
        <StatCard title="Outstanding Advances" value={`₹${Math.round(stats?.outstanding_advances ?? 0).toLocaleString()}`} icon={<AlertCircle className="h-4 w-4" />} />
      </div>

      {/* Salary Table */}
      <div className="rounded-lg border bg-card">
        <div className="p-6">
          <h3 className="text-lg font-semibold">
            {MONTHS[month - 1]} {year} — Salary Register
            {currentCycle && (
              <span className={`ml-3 inline-flex rounded-full px-2 py-1 text-xs font-medium ${currentCycle.status === 'generated' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                {currentCycle.status}
              </span>
            )}
          </h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-t bg-muted/50">
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Staff</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Base Salary</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Commission</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Bonus</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Advance</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Deduction</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Net Salary</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-muted-foreground">Status</th>
                <th className="px-6 py-3 text-right text-xs font-medium uppercase text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading && (
                <tr><td colSpan={9} className="px-6 py-8 text-center text-muted-foreground">Loading...</td></tr>
              )}
              {!isLoading && (!records || records.length === 0) && (
                <tr><td colSpan={9} className="px-6 py-8 text-center text-muted-foreground">No salary records. Generate salary to get started.</td></tr>
              )}
              {records?.map((rec) => (
                <tr key={rec.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm font-medium">{rec.staff_name}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{rec.base_salary.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{Math.round(rec.commission_amount).toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right">₹{rec.bonus_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right text-orange-600">-₹{rec.advance_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right text-red-600">-₹{rec.deduction_amount.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-right font-bold">₹{Math.round(rec.net_salary).toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${rec.payment_status === 'paid' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                      {rec.payment_status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    {rec.payment_status !== 'paid' && (
                      <button
                        onClick={() => handlePay(rec.id, rec.staff_name)}
                        className="rounded bg-green-600 px-3 py-1 text-xs font-medium text-white hover:bg-green-700"
                      >
                        Pay
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

function StatCard({ title, value, icon }: { title: string; value: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center gap-2 text-muted-foreground mb-2">
        {icon}
        <span className="text-sm font-medium">{title}</span>
      </div>
      <p className="text-2xl font-bold">{value}</p>
    </div>
  )
}
