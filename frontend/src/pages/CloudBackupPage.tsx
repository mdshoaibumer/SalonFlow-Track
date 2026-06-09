import { useState } from 'react'
import { useCloudConfig, useSaveCloudConfig, useTestCloudConnection, useBackupNow, useRestoreBackup, useCloudHistory, useCloudStats } from '@/hooks/useCloudBackup'
import type { CloudBackupConfig, CloudBackupHistory } from '@/types'

export function CloudBackupPage() {
  const [tab, setTab] = useState<'config' | 'history'>('config')

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Cloud Backup</h1>
        <p className="text-muted-foreground">Configure cloud storage, backup database, and restore from cloud</p>
      </div>

      <StatsBar />

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('config')} className={`px-4 py-2 rounded-t ${tab === 'config' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Configuration</button>
        <button onClick={() => setTab('history')} className={`px-4 py-2 rounded-t ${tab === 'history' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Backup History</button>
      </div>

      {tab === 'config' && <ConfigTab />}
      {tab === 'history' && <HistoryTab />}
    </div>
  )
}

function StatsBar() {
  const { data: stats } = useCloudStats()
  if (!stats) return null
  return (
    <div className="grid grid-cols-4 gap-4">
      <div className="border rounded p-3 text-center"><p className="text-lg font-bold">{stats.provider || 'None'}</p><p className="text-xs text-muted-foreground">Provider</p></div>
      <div className="border rounded p-3 text-center"><p className="text-2xl font-bold">{stats.total_backups}</p><p className="text-xs text-muted-foreground">Total Backups</p></div>
      <div className="border rounded p-3 text-center"><p className="text-lg font-bold">{stats.total_size_bytes ? (stats.total_size_bytes / 1024 / 1024).toFixed(1) + ' MB' : '0'}</p><p className="text-xs text-muted-foreground">Total Size</p></div>
      <div className="border rounded p-3 text-center"><p className="text-lg font-bold">{stats.auto_enabled ? 'ON' : 'OFF'}</p><p className="text-xs text-muted-foreground">Auto Backup</p></div>
    </div>
  )
}

function ConfigTab() {
  const { data: config, isLoading } = useCloudConfig()
  const saveMutation = useSaveCloudConfig()
  const testMutation = useTestCloudConnection()
  const backupMutation = useBackupNow()
  const [form, setForm] = useState<Partial<CloudBackupConfig>>({})

  if (isLoading) return <p>Loading...</p>

  const cfg = { ...config, ...form }

  const handleSave = () => saveMutation.mutate({ ...config, ...form })
  const handleTest = () => testMutation.mutate()
  const handleBackup = () => backupMutation.mutate()

  return (
    <div className="space-y-4 max-w-2xl">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-sm font-medium">Provider</label>
          <select className="w-full border rounded px-3 py-2 mt-1" value={cfg.provider || 'none'} onChange={e => setForm({ ...form, provider: e.target.value })}>
            <option value="none">None</option>
            <option value="google_drive">Google Drive</option>
            <option value="aws_s3">AWS S3</option>
            <option value="digitalocean_spaces">DigitalOcean Spaces</option>
          </select>
        </div>
        <div>
          <label className="text-sm font-medium">Bucket / Folder Name</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={cfg.bucket_name || ''} onChange={e => setForm({ ...form, bucket_name: e.target.value })} />
        </div>
        <div>
          <label className="text-sm font-medium">Region</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={cfg.region || ''} onChange={e => setForm({ ...form, region: e.target.value })} placeholder="ap-south-1" />
        </div>
        <div>
          <label className="text-sm font-medium">Access Key</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={cfg.access_key || ''} onChange={e => setForm({ ...form, access_key: e.target.value })} type="password" />
        </div>
        <div>
          <label className="text-sm font-medium">Endpoint (optional)</label>
          <input className="w-full border rounded px-3 py-2 mt-1" value={cfg.endpoint || ''} onChange={e => setForm({ ...form, endpoint: e.target.value })} placeholder="https://..." />
        </div>
        <div>
          <label className="text-sm font-medium">Auto Backup Interval (hours)</label>
          <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={cfg.auto_backup_interval_hours || 24} onChange={e => setForm({ ...form, auto_backup_interval_hours: parseInt(e.target.value) || 24 })} />
        </div>
        <div>
          <label className="text-sm font-medium">Max Versions</label>
          <input type="number" className="w-full border rounded px-3 py-2 mt-1" value={cfg.max_versions || 10} onChange={e => setForm({ ...form, max_versions: parseInt(e.target.value) || 10 })} />
        </div>
        <div className="space-y-2 pt-4">
          <label className="flex items-center gap-2">
            <input type="checkbox" checked={cfg.encrypt_backups ?? true} onChange={e => setForm({ ...form, encrypt_backups: e.target.checked })} />
            <span className="text-sm">Encrypt Backups</span>
          </label>
          <label className="flex items-center gap-2">
            <input type="checkbox" checked={cfg.auto_backup ?? false} onChange={e => setForm({ ...form, auto_backup: e.target.checked })} />
            <span className="text-sm">Auto Backup</span>
          </label>
        </div>
      </div>

      <div className="flex gap-2 pt-4">
        <button onClick={handleSave} disabled={saveMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded">
          {saveMutation.isPending ? 'Saving...' : 'Save Config'}
        </button>
        <button onClick={handleTest} disabled={testMutation.isPending} className="px-4 py-2 bg-blue-600 text-white rounded">
          {testMutation.isPending ? 'Testing...' : 'Test Connection'}
        </button>
        <button onClick={handleBackup} disabled={backupMutation.isPending} className="px-4 py-2 bg-green-600 text-white rounded">
          {backupMutation.isPending ? 'Backing up...' : 'Backup Now'}
        </button>
      </div>
      {testMutation.isSuccess && <p className="text-green-600 text-sm">Connection successful!</p>}
      {testMutation.isError && <p className="text-red-600 text-sm">{(testMutation.error as Error).message}</p>}
      {backupMutation.isSuccess && <p className="text-green-600 text-sm">Backup completed!</p>}
    </div>
  )
}

function HistoryTab() {
  const { data: history, isLoading } = useCloudHistory()
  const restoreMutation = useRestoreBackup()

  if (isLoading) return <p>Loading...</p>

  return (
    <div className="border rounded-lg overflow-hidden">
      <table className="w-full text-sm">
        <thead className="bg-muted">
          <tr>
            <th className="px-4 py-2 text-left">File</th>
            <th className="px-4 py-2 text-left">Size</th>
            <th className="px-4 py-2 text-left">Provider</th>
            <th className="px-4 py-2 text-left">Status</th>
            <th className="px-4 py-2 text-left">Date</th>
            <th className="px-4 py-2 text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          {(history?.data || []).map((h: CloudBackupHistory) => (
            <tr key={h.id} className="border-t">
              <td className="px-4 py-2">{h.file_name}</td>
              <td className="px-4 py-2">{(h.file_size / 1024).toFixed(1)} KB</td>
              <td className="px-4 py-2">{h.provider}</td>
              <td className="px-4 py-2"><span className={`px-2 py-0.5 rounded text-xs ${h.status === 'completed' ? 'bg-green-100' : h.status === 'failed' ? 'bg-red-100' : 'bg-yellow-100'}`}>{h.status}</span></td>
              <td className="px-4 py-2">{new Date(h.created_at).toLocaleString()}</td>
              <td className="px-4 py-2">
                {h.status === 'completed' && (
                  <button onClick={() => restoreMutation.mutate(h.id)} disabled={restoreMutation.isPending} className="text-xs px-2 py-1 bg-blue-100 rounded">Restore</button>
                )}
              </td>
            </tr>
          ))}
          {(!history?.data || history.data.length === 0) && <tr><td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">No backup history</td></tr>}
        </tbody>
      </table>
    </div>
  )
}
