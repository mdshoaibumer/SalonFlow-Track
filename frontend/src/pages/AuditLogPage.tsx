import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { ScrollText, Filter, ChevronLeft, ChevronRight } from 'lucide-react'
import { useAuth } from '@/app/providers/AuthProvider'
import * as authService from '@/services/auth'
import type { AuditFilter } from '@/types/auth'

export function AuditLogPage() {
  const { hasPermission } = useAuth()
  const [filter, setFilter] = useState<AuditFilter>({ page: 1, per_page: 50 })

  const { data, isLoading } = useQuery({
    queryKey: ['audit-logs', filter],
    queryFn: () => authService.getAuditLogs(filter),
  })

  if (!hasPermission('audit.view')) {
    return <div className="p-6 text-muted-foreground">You do not have permission to view audit logs.</div>
  }

  const logs = data?.logs || []
  const total = data?.total || 0
  const totalPages = Math.ceil(total / (filter.per_page || 50))

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <ScrollText className="h-6 w-6 text-violet-600" />
        <h1 className="text-xl font-bold">Audit Logs</h1>
        <span className="ml-auto text-sm text-muted-foreground">{total} entries</span>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-3 items-center">
        <Filter className="h-4 w-4 text-muted-foreground" />
        <select
          value={filter.module || ''}
          onChange={(e) => setFilter({ ...filter, module: e.target.value, page: 1 })}
          className="rounded-lg border border-border px-3 py-1.5 text-sm bg-background"
        >
          <option value="">All Modules</option>
          {['auth', 'customers', 'billing', 'staff', 'services', 'salary', 'expenses', 'inventory', 'backup', 'users', 'system'].map((m) => (
            <option key={m} value={m}>{m}</option>
          ))}
        </select>
        <select
          value={filter.severity || ''}
          onChange={(e) => setFilter({ ...filter, severity: e.target.value, page: 1 })}
          className="rounded-lg border border-border px-3 py-1.5 text-sm bg-background"
        >
          <option value="">All Severities</option>
          <option value="info">Info</option>
          <option value="warning">Warning</option>
          <option value="critical">Critical</option>
        </select>
        <input
          type="date"
          value={filter.from_date || ''}
          onChange={(e) => setFilter({ ...filter, from_date: e.target.value, page: 1 })}
          className="rounded-lg border border-border px-3 py-1.5 text-sm bg-background"
          placeholder="From"
        />
        <input
          type="date"
          value={filter.to_date || ''}
          onChange={(e) => setFilter({ ...filter, to_date: e.target.value, page: 1 })}
          className="rounded-lg border border-border px-3 py-1.5 text-sm bg-background"
          placeholder="To"
        />
      </div>

      {/* Table */}
      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading...</div>
      ) : (
        <div className="rounded-lg border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-muted/50">
                <tr>
                  <th className="px-3 py-2.5 text-left font-medium">Timestamp</th>
                  <th className="px-3 py-2.5 text-left font-medium">User</th>
                  <th className="px-3 py-2.5 text-left font-medium">Action</th>
                  <th className="px-3 py-2.5 text-left font-medium">Module</th>
                  <th className="px-3 py-2.5 text-left font-medium">Description</th>
                  <th className="px-3 py-2.5 text-left font-medium">Severity</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {logs.map((log) => (
                  <tr key={log.id} className="hover:bg-muted/30">
                    <td className="px-3 py-2 text-xs text-muted-foreground whitespace-nowrap">
                      {new Date(log.timestamp).toLocaleString()}
                    </td>
                    <td className="px-3 py-2 text-xs">{log.username || '—'}</td>
                    <td className="px-3 py-2 text-xs font-medium">{log.action}</td>
                    <td className="px-3 py-2 text-xs">{log.module}</td>
                    <td className="px-3 py-2 text-xs max-w-[300px] truncate">{log.description}</td>
                    <td className="px-3 py-2">
                      <SeverityBadge severity={log.severity} />
                    </td>
                  </tr>
                ))}
                {logs.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-3 py-8 text-center text-muted-foreground">
                      No audit logs found
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">
            Page {filter.page} of {totalPages}
          </span>
          <div className="flex gap-1">
            <button
              disabled={filter.page === 1}
              onClick={() => setFilter({ ...filter, page: (filter.page || 1) - 1 })}
              className="rounded p-1.5 hover:bg-muted disabled:opacity-50"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              disabled={filter.page === totalPages}
              onClick={() => setFilter({ ...filter, page: (filter.page || 1) + 1 })}
              className="rounded p-1.5 hover:bg-muted disabled:opacity-50"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

function SeverityBadge({ severity }: { severity: string }) {
  const styles = {
    info: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
    warning: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
    critical: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300',
  }
  return (
    <span className={`inline-flex rounded-full px-2 py-0.5 text-[11px] font-medium ${styles[severity as keyof typeof styles] || styles.info}`}>
      {severity}
    </span>
  )
}
