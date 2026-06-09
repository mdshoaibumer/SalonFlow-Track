import { useState } from 'react'
import { useBackups, useBackupStats, useRestores, useCreateBackup, useVerifyBackup, useRestoreBackup, useDeleteBackup } from '@/hooks/useBackup'
import { HardDrive, Plus, ShieldCheck, RotateCcw, Trash2, CheckCircle2, XCircle, AlertCircle, Clock } from 'lucide-react'

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short' })
}

function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    completed: 'bg-green-100 text-green-800',
    verified: 'bg-blue-100 text-blue-800',
    pending: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
    corrupted: 'bg-red-100 text-red-800',
  }
  const icons: Record<string, React.ReactNode> = {
    completed: <CheckCircle2 className="h-3 w-3" />,
    verified: <ShieldCheck className="h-3 w-3" />,
    pending: <Clock className="h-3 w-3" />,
    failed: <XCircle className="h-3 w-3" />,
    corrupted: <AlertCircle className="h-3 w-3" />,
  }
  return (
    <span className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${styles[status] || 'bg-gray-100 text-gray-800'}`}>
      {icons[status]}
      {status}
    </span>
  )
}

export function BackupPage() {
  const [tab, setTab] = useState<'backups' | 'restores'>('backups')
  const [page, setPage] = useState(1)
  const [restorePage, setRestorePage] = useState(1)

  const { data: stats } = useBackupStats()
  const { data: backupData, isLoading } = useBackups(page)
  const { data: restoreData } = useRestores(restorePage)
  const createBackup = useCreateBackup()
  const verifyBackup = useVerifyBackup()
  const restoreBackup = useRestoreBackup()
  const deleteBackup = useDeleteBackup()

  const handleCreate = () => {
    createBackup.mutate('manual')
  }

  const handleVerify = (id: string) => {
    verifyBackup.mutate(id)
  }

  const handleRestore = (id: string, name: string) => {
    if (confirm(`Restore from backup "${name}"? This will replace the current database.`)) {
      restoreBackup.mutate({ id })
    }
  }

  const handleDelete = (id: string, name: string) => {
    if (confirm(`Delete backup "${name}"? This cannot be undone.`)) {
      deleteBackup.mutate(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Backup & Restore</h1>
          <p className="text-muted-foreground">Manage database backups and restore points</p>
        </div>
        <button
          onClick={handleCreate}
          disabled={createBackup.isPending}
          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
        >
          <Plus className="h-4 w-4" />
          {createBackup.isPending ? 'Creating...' : 'Create Backup'}
        </button>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid gap-4 md:grid-cols-4">
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <HardDrive className="h-4 w-4" />
              Total Backups
            </div>
            <p className="mt-1 text-2xl font-bold">{stats.total_backups}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="text-sm text-muted-foreground">Last Backup</div>
            <p className="mt-1 text-sm font-medium">{stats.last_backup_name || '—'}</p>
            <p className="text-xs text-muted-foreground">{formatDate(stats.last_backup_date)}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="text-sm text-muted-foreground">Last Backup Size</div>
            <p className="mt-1 text-2xl font-bold">{formatBytes(stats.last_backup_size)}</p>
          </div>
          <div className="rounded-lg border bg-card p-4">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <RotateCcw className="h-4 w-4" />
              Total Restores
            </div>
            <p className="mt-1 text-2xl font-bold">{stats.total_restores}</p>
          </div>
        </div>
      )}

      {/* Tabs */}
      <div className="border-b">
        <div className="flex gap-4">
          <button
            onClick={() => setTab('backups')}
            className={`border-b-2 px-4 py-2 text-sm font-medium ${tab === 'backups' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'}`}
          >
            Backup History
          </button>
          <button
            onClick={() => setTab('restores')}
            className={`border-b-2 px-4 py-2 text-sm font-medium ${tab === 'restores' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'}`}
          >
            Restore History
          </button>
        </div>
      </div>

      {/* Backup History Tab */}
      {tab === 'backups' && (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Name</th>
                <th className="px-4 py-3 text-left font-medium">Type</th>
                <th className="px-4 py-3 text-left font-medium">Size</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Created</th>
                <th className="px-4 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">Loading...</td></tr>
              ) : !backupData?.backups?.length ? (
                <tr><td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">No backups yet. Create your first backup above.</td></tr>
              ) : (
                backupData.backups.map((b) => (
                  <tr key={b.id} className="border-b last:border-0">
                    <td className="px-4 py-3 font-medium">{b.backup_name}</td>
                    <td className="px-4 py-3 capitalize">{b.backup_type}</td>
                    <td className="px-4 py-3">{formatBytes(b.file_size)}</td>
                    <td className="px-4 py-3"><StatusBadge status={b.status} /></td>
                    <td className="px-4 py-3">{formatDate(b.created_at)}</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        <button
                          onClick={() => handleVerify(b.id)}
                          disabled={verifyBackup.isPending}
                          title="Verify"
                          className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                        >
                          <ShieldCheck className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleRestore(b.id, b.backup_name)}
                          disabled={restoreBackup.isPending}
                          title="Restore"
                          className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                        >
                          <RotateCcw className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDelete(b.id, b.backup_name)}
                          disabled={deleteBackup.isPending}
                          title="Delete"
                          className="rounded p-1.5 text-muted-foreground hover:bg-red-50 hover:text-red-600"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
          {backupData && backupData.meta.total_pages > 1 && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <span className="text-sm text-muted-foreground">
                Page {backupData.meta.page} of {backupData.meta.total_pages}
              </span>
              <div className="flex gap-2">
                <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Prev</button>
                <button onClick={() => setPage((p) => p + 1)} disabled={page >= backupData.meta.total_pages} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Restore History Tab */}
      {tab === 'restores' && (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Backup Name</th>
                <th className="px-4 py-3 text-left font-medium">Restore Date</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Notes</th>
              </tr>
            </thead>
            <tbody>
              {!restoreData?.restores?.length ? (
                <tr><td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">No restores performed yet.</td></tr>
              ) : (
                restoreData.restores.map((r) => (
                  <tr key={r.id} className="border-b last:border-0">
                    <td className="px-4 py-3 font-medium">{r.backup_name}</td>
                    <td className="px-4 py-3">{formatDate(r.restore_date)}</td>
                    <td className="px-4 py-3"><StatusBadge status={r.status} /></td>
                    <td className="px-4 py-3 text-muted-foreground">{r.notes || '—'}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
          {restoreData && restoreData.meta.total_pages > 1 && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <span className="text-sm text-muted-foreground">
                Page {restoreData.meta.page} of {restoreData.meta.total_pages}
              </span>
              <div className="flex gap-2">
                <button onClick={() => setRestorePage((p) => Math.max(1, p - 1))} disabled={restorePage === 1} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Prev</button>
                <button onClick={() => setRestorePage((p) => p + 1)} disabled={restorePage >= restoreData.meta.total_pages} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
