import { useState } from 'react'
import { useImportJobs, useUploadFile, useValidateImport, useProcessImport } from '@/hooks/useImport'
import { Upload, FileSpreadsheet, CheckCircle2, XCircle, AlertTriangle, ArrowRight, Loader2 } from 'lucide-react'
import type { ColumnMapping, ImportPreview, ImportUploadResult } from '@/types'

const TARGET_ENTITIES = [
  { value: '', label: 'Auto-detect' },
  { value: 'staff', label: 'Staff' },
  { value: 'customers', label: 'Customers' },
  { value: 'services', label: 'Services' },
  { value: 'products', label: 'Products' },
  { value: 'expenses', label: 'Expenses' },
  { value: 'advances', label: 'Advances' },
  { value: 'salary', label: 'Salary' },
]

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short' })
}

function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    completed: 'bg-green-100 text-green-800',
    validated: 'bg-blue-100 text-blue-800',
    importing: 'bg-yellow-100 text-yellow-800',
    pending: 'bg-gray-100 text-gray-800',
    validating: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
  }
  return (
    <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${styles[status] || 'bg-gray-100 text-gray-800'}`}>
      {status}
    </span>
  )
}

type WizardStep = 'upload' | 'mapping' | 'preview' | 'done'

export function ImportPage() {
  const [step, setStep] = useState<WizardStep>('upload')
  const [targetEntity, setTargetEntity] = useState('')
  const [uploadResult, setUploadResult] = useState<ImportUploadResult | null>(null)
  const [mappings, setMappings] = useState<ColumnMapping[]>([])
  const [preview, setPreview] = useState<ImportPreview | null>(null)
  const [page, setPage] = useState(1)

  const { data: jobsData, isLoading } = useImportJobs(page)
  const uploadMut = useUploadFile()
  const validateMut = useValidateImport()
  const processMut = useProcessImport()

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    uploadMut.mutate({ file, targetEntity: targetEntity || undefined }, {
      onSuccess: (result) => {
        setUploadResult(result)
        setMappings(result.mappings)
        setStep('mapping')
      },
    })
  }

  const handleValidate = () => {
    if (!uploadResult) return
    validateMut.mutate({ jobId: uploadResult.job.id, mappings }, {
      onSuccess: (prev) => {
        setPreview(prev)
        setStep('preview')
      },
    })
  }

  const handleProcess = () => {
    if (!uploadResult) return
    processMut.mutate(uploadResult.job.id, {
      onSuccess: () => {
        setStep('done')
      },
    })
  }

  const handleReset = () => {
    setStep('upload')
    setUploadResult(null)
    setMappings([])
    setPreview(null)
    setTargetEntity('')
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Data Import</h1>
          <p className="text-muted-foreground">Import data from Excel or CSV files</p>
        </div>
      </div>

      {/* Import Wizard */}
      <div className="rounded-lg border bg-card p-6">
        {/* Step indicators */}
        <div className="mb-6 flex items-center gap-2 text-sm">
          {(['upload', 'mapping', 'preview', 'done'] as WizardStep[]).map((s, i) => (
            <div key={s} className="flex items-center gap-2">
              {i > 0 && <ArrowRight className="h-4 w-4 text-muted-foreground" />}
              <span className={`rounded-full px-3 py-1 ${step === s ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}>
                {i + 1}. {s === 'upload' ? 'Upload' : s === 'mapping' ? 'Map Columns' : s === 'preview' ? 'Preview' : 'Complete'}
              </span>
            </div>
          ))}
        </div>

        {/* Step 1: Upload */}
        {step === 'upload' && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium">Target Entity</label>
              <select value={targetEntity} onChange={(e) => setTargetEntity(e.target.value)} className="mt-1 block w-full rounded-lg border bg-background px-3 py-2 text-sm">
                {TARGET_ENTITIES.map((e) => <option key={e.value} value={e.value}>{e.label}</option>)}
              </select>
            </div>
            <div className="flex items-center justify-center rounded-lg border-2 border-dashed p-12">
              <label className="flex cursor-pointer flex-col items-center gap-3">
                <Upload className="h-10 w-10 text-muted-foreground" />
                <span className="text-sm font-medium">Click to upload .xlsx, .xls, or .csv</span>
                <span className="text-xs text-muted-foreground">Max 50MB</span>
                <input type="file" accept=".xlsx,.xls,.csv" onChange={handleFileUpload} className="hidden" />
              </label>
            </div>
            {uploadMut.isPending && (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" /> Parsing file...
              </div>
            )}
            {uploadMut.isError && (
              <p className="text-sm text-red-600">{(uploadMut.error as Error).message}</p>
            )}
          </div>
        )}

        {/* Step 2: Column Mapping */}
        {step === 'mapping' && uploadResult && (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">{uploadResult.job.file_name}</p>
                <p className="text-sm text-muted-foreground">Entity: {uploadResult.job.target_entity} · {uploadResult.job.total_rows} rows detected</p>
              </div>
            </div>
            <div className="max-h-96 overflow-auto rounded border">
              <table className="w-full text-sm">
                <thead className="sticky top-0 border-b bg-muted/50">
                  <tr>
                    <th className="px-4 py-2 text-left font-medium">Source Column</th>
                    <th className="px-4 py-2 text-left font-medium">Maps To</th>
                  </tr>
                </thead>
                <tbody>
                  {uploadResult.headers.map((header) => {
                    const mapping = mappings.find((m) => m.source_column === header)
                    return (
                      <tr key={header} className="border-b last:border-0">
                        <td className="px-4 py-2 font-mono text-xs">{header}</td>
                        <td className="px-4 py-2">
                          <input
                            type="text"
                            value={mapping?.target_field || ''}
                            onChange={(e) => {
                              const val = e.target.value
                              setMappings((prev) => {
                                const existing = prev.find((m) => m.source_column === header)
                                if (existing) {
                                  return prev.map((m) => m.source_column === header ? { ...m, target_field: val } : m)
                                }
                                return [...prev, { source_column: header, target_field: val }]
                              })
                            }}
                            placeholder="(skip)"
                            className="w-full rounded border bg-background px-2 py-1 text-sm"
                          />
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
            <div className="flex gap-2">
              <button onClick={handleValidate} disabled={validateMut.isPending} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
                {validateMut.isPending ? 'Validating...' : 'Validate & Preview'}
              </button>
              <button onClick={handleReset} className="rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted">Cancel</button>
            </div>
          </div>
        )}

        {/* Step 3: Preview */}
        {step === 'preview' && preview && (
          <div className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-4">
              <div className="rounded-lg border p-4 text-center">
                <FileSpreadsheet className="mx-auto h-6 w-6 text-muted-foreground" />
                <p className="mt-1 text-2xl font-bold">{preview.total_rows}</p>
                <p className="text-xs text-muted-foreground">Total Rows</p>
              </div>
              <div className="rounded-lg border p-4 text-center">
                <CheckCircle2 className="mx-auto h-6 w-6 text-green-600" />
                <p className="mt-1 text-2xl font-bold text-green-600">{preview.valid_rows}</p>
                <p className="text-xs text-muted-foreground">Valid</p>
              </div>
              <div className="rounded-lg border p-4 text-center">
                <XCircle className="mx-auto h-6 w-6 text-red-600" />
                <p className="mt-1 text-2xl font-bold text-red-600">{preview.invalid_rows}</p>
                <p className="text-xs text-muted-foreground">Invalid</p>
              </div>
              <div className="rounded-lg border p-4 text-center">
                <AlertTriangle className="mx-auto h-6 w-6 text-yellow-600" />
                <p className="mt-1 text-2xl font-bold text-yellow-600">{preview.warnings}</p>
                <p className="text-xs text-muted-foreground">Warnings</p>
              </div>
            </div>

            {preview.errors.length > 0 && (
              <div className="max-h-48 overflow-auto rounded border">
                <table className="w-full text-sm">
                  <thead className="sticky top-0 border-b bg-muted/50">
                    <tr>
                      <th className="px-3 py-2 text-left font-medium">Row</th>
                      <th className="px-3 py-2 text-left font-medium">Error</th>
                    </tr>
                  </thead>
                  <tbody>
                    {preview.errors.map((err, i) => (
                      <tr key={i} className="border-b last:border-0">
                        <td className="px-3 py-2">{err.row_number}</td>
                        <td className="px-3 py-2 text-red-600">{err.message}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            <div className="flex gap-2">
              <button
                onClick={handleProcess}
                disabled={processMut.isPending || preview.valid_rows === 0}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
              >
                {processMut.isPending ? 'Importing...' : `Import ${preview.valid_rows} Valid Rows`}
              </button>
              <button onClick={handleReset} className="rounded-lg border px-4 py-2 text-sm font-medium hover:bg-muted">Cancel</button>
            </div>
          </div>
        )}

        {/* Step 4: Done */}
        {step === 'done' && (
          <div className="py-8 text-center">
            <CheckCircle2 className="mx-auto h-12 w-12 text-green-600" />
            <h3 className="mt-3 text-lg font-semibold">Import Complete</h3>
            <p className="text-muted-foreground">Data has been imported successfully.</p>
            <button onClick={handleReset} className="mt-4 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
              Import Another File
            </button>
          </div>
        )}
      </div>

      {/* Import History */}
      <div>
        <h3 className="mb-3 text-lg font-semibold">Import History</h3>
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">File</th>
                <th className="px-4 py-3 text-left font-medium">Entity</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Rows</th>
                <th className="px-4 py-3 text-left font-medium">Imported</th>
                <th className="px-4 py-3 text-left font-medium">Date</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">Loading...</td></tr>
              ) : !jobsData?.jobs?.length ? (
                <tr><td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">No imports yet.</td></tr>
              ) : (
                jobsData.jobs.map((job: any) => (
                  <tr key={job.id} className="border-b last:border-0">
                    <td className="px-4 py-3 font-medium">{job.file_name}</td>
                    <td className="px-4 py-3 capitalize">{job.target_entity}</td>
                    <td className="px-4 py-3"><StatusBadge status={job.status} /></td>
                    <td className="px-4 py-3">{job.total_rows}</td>
                    <td className="px-4 py-3">{job.imported_rows}/{job.valid_rows}</td>
                    <td className="px-4 py-3">{formatDate(job.created_at)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
          {jobsData && jobsData.meta.total_pages > 1 && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <span className="text-sm text-muted-foreground">Page {jobsData.meta.page} of {jobsData.meta.total_pages}</span>
              <div className="flex gap-2">
                <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Prev</button>
                <button onClick={() => setPage((p) => p + 1)} disabled={page >= jobsData.meta.total_pages} className="rounded border px-3 py-1 text-sm disabled:opacity-50">Next</button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
