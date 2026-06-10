import { useState } from 'react'
import { useSalaryList, useSalaryStats, useGenerateSalary, usePaySalary, useSalaryCycles } from '@/hooks/useSalary'
import { Clock, CheckCircle, AlertCircle, Wallet } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { PageHeader } from '@/components/shared/PageHeader'
import { KPICard } from '@/components/shared/KPICard'
import { DataTable, type ColumnDef } from '@/components/shared/DataTable'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorState } from '@/components/shared/ErrorState'

const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']

interface SalaryRecord {
  id: string
  staff_name: string
  base_salary: number
  commission_amount: number
  bonus_amount: number
  advance_amount: number
  deduction_amount: number
  net_salary: number
  payment_status: string
}

export function SalaryPage() {
  const now = new Date()
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [year, setYear] = useState(now.getFullYear())

  const { data: stats } = useSalaryStats()
  const { data: records, isLoading, error, refetch } = useSalaryList(month, year)
  const { data: cycles } = useSalaryCycles(year)
  const generateSalary = useGenerateSalary()
  const paySalary = usePaySalary()

  const currentCycle = cycles?.find((c: { month: number; year: number }) => c.month === month && c.year === year)

  const handleGenerate = () => {
    generateSalary.mutate({ month, year })
  }

  const handlePay = (id: string) => {
    paySalary.mutate(id)
  }

  const handleExport = () => {
    if (!records) return
    const csv = [
      ['Staff', 'Base Salary', 'Commission', 'Bonus', 'Advance', 'Deduction', 'Net Salary', 'Status'].join(','),
      ...records.map((r: SalaryRecord) =>
        [r.staff_name, r.base_salary, Math.round(r.commission_amount), r.bonus_amount, r.advance_amount, r.deduction_amount, Math.round(r.net_salary), r.payment_status].join(',')
      ),
    ].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `salary-${MONTHS[month - 1]}-${year}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<SalaryRecord, unknown>[] = [
    {
      accessorKey: 'staff_name',
      header: 'Staff',
      cell: ({ row }) => <span className="font-medium">{row.original.staff_name}</span>,
    },
    {
      accessorKey: 'base_salary',
      header: 'Base Salary',
      cell: ({ row }) => `₹${row.original.base_salary.toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'commission_amount',
      header: 'Commission',
      cell: ({ row }) => `₹${Math.round(row.original.commission_amount).toLocaleString('en-IN')}`,
    },
    {
      accessorKey: 'bonus_amount',
      header: 'Bonus',
      cell: ({ row }) => row.original.bonus_amount > 0 ? `₹${row.original.bonus_amount.toLocaleString('en-IN')}` : '-',
    },
    {
      accessorKey: 'advance_amount',
      header: 'Advance',
      cell: ({ row }) => row.original.advance_amount > 0 ? <span className="text-orange-600">-₹{row.original.advance_amount.toLocaleString('en-IN')}</span> : '-',
    },
    {
      accessorKey: 'deduction_amount',
      header: 'Deduction',
      cell: ({ row }) => row.original.deduction_amount > 0 ? <span className="text-red-600">-₹{row.original.deduction_amount.toLocaleString('en-IN')}</span> : '-',
    },
    {
      accessorKey: 'net_salary',
      header: 'Net Salary',
      cell: ({ row }) => <span className="font-bold">₹{Math.round(row.original.net_salary).toLocaleString('en-IN')}</span>,
    },
    {
      accessorKey: 'payment_status',
      header: 'Status',
      cell: ({ row }) => (
        <Badge variant={row.original.payment_status === 'paid' ? 'default' : 'secondary'}>
          {row.original.payment_status}
        </Badge>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) =>
        row.original.payment_status !== 'paid' ? (
          <Button size="sm" variant="outline" onClick={() => handlePay(row.original.id)} className="text-green-700 border-green-300 hover:bg-green-50">
            Mark Paid
          </Button>
        ) : null,
    },
  ]

  if (error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Salary Management" description="Generate and manage monthly payroll" />
        <ErrorState title="Failed to load salary data" onRetry={() => refetch()} />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Salary Management"
        description="Generate and manage monthly payroll"
        actions={
          <div className="flex items-center gap-2">
            {!currentCycle && (
              <Button onClick={handleGenerate} disabled={generateSalary.isPending} size="sm">
                {generateSalary.isPending ? 'Generating...' : 'Generate Salary'}
              </Button>
            )}
          </div>
        }
      />

      {/* KPIs */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <KPICard title="Total Payroll" value={`₹${Math.round(stats?.total_payroll ?? 0).toLocaleString('en-IN')}`} icon={Wallet} />
        <KPICard title="Pending" value={String(stats?.pending_payments ?? 0)} icon={Clock} />
        <KPICard title="Paid" value={String(stats?.paid_salaries ?? 0)} icon={CheckCircle} />
        <KPICard title="Outstanding Advances" value={`₹${Math.round(stats?.outstanding_advances ?? 0).toLocaleString('en-IN')}`} icon={AlertCircle} />
      </div>

      {/* Month/Year Selector */}
      <div className="flex items-center gap-3">
        <Select value={String(month)} onValueChange={(v) => setMonth(Number(v))}>
          <SelectTrigger className="w-[140px] h-9">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {MONTHS.map((m, i) => (
              <SelectItem key={i} value={String(i + 1)}>{m}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={String(year)} onValueChange={(v) => setYear(Number(v))}>
          <SelectTrigger className="w-[100px] h-9">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {[2024, 2025, 2026, 2027].map((y) => (
              <SelectItem key={y} value={String(y)}>{y}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        {currentCycle && (
          <Badge variant="outline" className="ml-2">
            {currentCycle.status}
          </Badge>
        )}
      </div>

      {/* Data Table */}
      {isLoading ? (
        <LoadingState variant="table" />
      ) : (
        <DataTable
          columns={columns}
          data={records || []}
          searchPlaceholder="Search staff..."
          emptyTitle="No salary records"
          emptyDescription="Generate salary for this month to get started."
          emptyAction={!currentCycle ? { label: 'Generate Salary', onClick: handleGenerate } : undefined}
          onExport={records && records.length > 0 ? handleExport : undefined}
          exportLabel="Export CSV"
        />
      )}
    </div>
  )
}
