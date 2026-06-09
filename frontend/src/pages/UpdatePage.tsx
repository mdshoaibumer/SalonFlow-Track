import { useUpdateStatus, useUpdateHistory, useCheckForUpdate, useDownloadUpdate, useInstallUpdate } from '@/hooks/useUpdate'
import { Download, RefreshCw, CheckCircle2, XCircle, ArrowUpCircle, Clock } from 'lucide-react'

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short' })
}

function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    completed: 'bg-green-100 text-green-800',
    installed: 'bg-green-100 text-green-800',
    downloaded: 'bg-blue-100 text-blue-800',
    downloading: 'bg-blue-100 text-blue-800',
    pending: 'bg-yellow-100 text-yellow-800',
    installing: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
    rolled_back: 'bg-red-100 text-red-800',
  }
  const icons: Record<string, React.ReactNode> = {
    completed: <CheckCircle2 className="h-3 w-3" />,
    installed: <CheckCircle2 className="h-3 w-3" />,
    failed: <XCircle className="h-3 w-3" />,
    rolled_back: <XCircle className="h-3 w-3" />,
    downloading: <Download className="h-3 w-3" />,
    pending: <Clock className="h-3 w-3" />,
  }
  return (
    <span className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${styles[status] || 'bg-gray-100 text-gray-800'}`}>
      {icons[status]}
      {status.replace('_', ' ')}
    </span>
  )
}

export function UpdatePage() {
  const { data: status } = useUpdateStatus()
  const { data: historyData } = useUpdateHistory()
  const checkMut = useCheckForUpdate()
  const downloadMut = useDownloadUpdate()
  const installMut = useInstallUpdate()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Updates</h1>
          <p className="text-muted-foreground">Manage application updates and releases</p>
        </div>
        <button
          onClick={() => checkMut.mutate()}
          disabled={checkMut.isPending}
          className="inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted"
        >
          <RefreshCw className={`h-4 w-4 ${checkMut.isPending ? 'animate-spin' : ''}`} />
          Check for Updates
        </button>
      </div>

      {/* Current Version Card */}
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-start gap-4">
          <ArrowUpCircle className="h-8 w-8 text-primary" />
          <div className="flex-1">
            <div className="flex items-center gap-3">
              <h2 className="text-xl font-semibold">v{status?.current_version || '1.0.0'}</h2>
              {status?.status === 'up_to_date' && (
                <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-800">
                  <CheckCircle2 className="h-3 w-3" /> Up to date
                </span>
              )}
            </div>
            <p className="mt-1 text-sm text-muted-foreground">
              {status?.update_available
                ? `New version v${status.latest_version} is available`
                : 'You are running the latest version'}
            </p>
          </div>
        </div>
      </div>

      {/* Update Available Banner */}
      {checkMut.data?.update_available && (
        <div className="rounded-lg border border-blue-200 bg-blue-50 p-6">
          <div className="flex items-start justify-between">
            <div>
              <h3 className="text-lg font-semibold text-blue-900">v{checkMut.data.latest_version} Available</h3>
              {checkMut.data.release_notes && (
                <p className="mt-2 whitespace-pre-wrap text-sm text-blue-800">{checkMut.data.release_notes}</p>
              )}
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => downloadMut.mutate()}
                disabled={downloadMut.isPending}
                className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
              >
                <Download className="h-4 w-4" />
                {downloadMut.isPending ? 'Downloading...' : 'Download'}
              </button>
              <button
                onClick={() => installMut.mutate()}
                disabled={installMut.isPending}
                className="inline-flex items-center gap-2 rounded-lg border border-primary px-4 py-2 text-sm font-medium text-primary hover:bg-primary/10 disabled:opacity-50"
              >
                {installMut.isPending ? 'Installing...' : 'Install'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Update History */}
      <div>
        <h3 className="mb-3 text-lg font-semibold">Update History</h3>
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">From</th>
                <th className="px-4 py-3 text-left font-medium">To</th>
                <th className="px-4 py-3 text-left font-medium">Date</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Error</th>
              </tr>
            </thead>
            <tbody>
              {!historyData?.records?.length ? (
                <tr><td colSpan={5} className="px-4 py-8 text-center text-muted-foreground">No update history yet.</td></tr>
              ) : (
                historyData.records.map((rec) => (
                  <tr key={rec.id} className="border-b last:border-0">
                    <td className="px-4 py-3 font-mono text-xs">v{rec.from_version}</td>
                    <td className="px-4 py-3 font-mono text-xs">v{rec.to_version}</td>
                    <td className="px-4 py-3">{formatDate(rec.update_date)}</td>
                    <td className="px-4 py-3"><StatusBadge status={rec.status} /></td>
                    <td className="px-4 py-3 text-muted-foreground">{rec.error_message || '—'}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
