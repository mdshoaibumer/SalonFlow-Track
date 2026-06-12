import { useQuery, useMutation } from '@tanstack/react-query'
import { Activity, Download, Database, Cpu, HardDrive, Clock, Users, FileText } from 'lucide-react'
import { toast } from 'sonner'
import { useAuth } from '@/app/providers/AuthProvider'
import * as authService from '@/services/auth'

export function DiagnosticsPage() {
  const { hasPermission } = useAuth()

  const { data: diag, isLoading } = useQuery({
    queryKey: ['diagnostics'],
    queryFn: authService.getDiagnostics,
    refetchInterval: 30000, // Refresh every 30s
  })

  const exportMutation = useMutation({
    mutationFn: authService.exportDiagnosticsBundle,
    onSuccess: (path) => {
      toast.success(`Diagnostics exported to: ${path}`)
    },
    onError: () => toast.error('Failed to export diagnostics'),
  })

  if (!hasPermission('diagnostics.view')) {
    return <div className="p-6 text-muted-foreground">You do not have permission to view diagnostics.</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Activity className="h-6 w-6 text-violet-600" />
          <h1 className="text-xl font-bold">System Diagnostics</h1>
        </div>
        {hasPermission('diagnostics.export') && (
          <button
            onClick={() => exportMutation.mutate()}
            disabled={exportMutation.isPending}
            className="flex items-center gap-2 rounded-lg bg-violet-600 px-4 py-2 text-sm font-medium text-white hover:bg-violet-700 disabled:opacity-50 transition-colors"
          >
            <Download className="h-4 w-4" />
            {exportMutation.isPending ? 'Exporting...' : 'Export Diagnostics Bundle'}
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading diagnostics...</div>
      ) : diag ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {/* Application */}
          <DiagCard icon={Cpu} title="Application" items={[
            { label: 'Version', value: diag.app_version },
            { label: 'Go Version', value: diag.go_version },
            { label: 'OS / Arch', value: `${diag.os} / ${diag.arch}` },
            { label: 'Uptime', value: diag.uptime },
          ]} />

          {/* Database */}
          <DiagCard icon={Database} title="Database" items={[
            { label: 'Path', value: diag.database_path },
            { label: 'Size', value: formatBytes(diag.database_size_bytes) },
            { label: 'Migration Version', value: diag.db_version },
          ]} />

          {/* Resources */}
          <DiagCard icon={HardDrive} title="Resources" items={[
            { label: 'CPUs', value: String(diag.num_cpu) },
            { label: 'Goroutines', value: String(diag.num_goroutine) },
            { label: 'Memory (Alloc)', value: `${diag.mem_alloc_mb.toFixed(1)} MB` },
            { label: 'Memory (Total)', value: `${diag.mem_total_alloc_mb.toFixed(1)} MB` },
          ]} />

          {/* Counts */}
          <DiagCard icon={Users} title="Data Counts" items={[
            { label: 'Users', value: String(diag.total_users) },
            { label: 'Customers', value: String(diag.total_customers) },
            { label: 'Invoices', value: String(diag.total_invoices) },
          ]} />

          {/* Backup & Logs */}
          <DiagCard icon={Clock} title="Maintenance" items={[
            { label: 'Last Backup', value: diag.last_backup || 'Never' },
            { label: 'Log Directory', value: diag.log_directory },
          ]} />

          {/* Support */}
          <div className="rounded-lg border border-border p-4 space-y-3">
            <div className="flex items-center gap-2">
              <FileText className="h-4 w-4 text-violet-600" />
              <h3 className="font-medium text-sm">Support</h3>
            </div>
            <p className="text-xs text-muted-foreground">
              If you encounter an issue, export the diagnostics bundle and send it to the developer for investigation.
            </p>
            <p className="text-xs text-muted-foreground">
              The bundle includes: log files, database metadata, recent audit logs, and system information.
            </p>
          </div>
        </div>
      ) : null}
    </div>
  )
}

function DiagCard({ icon: Icon, title, items }: {
  icon: any
  title: string
  items: { label: string; value: string }[]
}) {
  return (
    <div className="rounded-lg border border-border p-4 space-y-3">
      <div className="flex items-center gap-2">
        <Icon className="h-4 w-4 text-violet-600" />
        <h3 className="font-medium text-sm">{title}</h3>
      </div>
      <dl className="space-y-1.5">
        {items.map((item) => (
          <div key={item.label} className="flex justify-between text-xs">
            <dt className="text-muted-foreground">{item.label}</dt>
            <dd className="font-medium truncate max-w-[180px]" title={item.value}>{item.value}</dd>
          </div>
        ))}
      </dl>
    </div>
  )
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}
