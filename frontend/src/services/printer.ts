import type { PrinterSettings, PrintJob, ReceiptData } from '@/types'

export async function getPrinterSettings(): Promise<PrinterSettings> {
  return window.go.main.PrinterService.GetSettings()
}

export async function savePrinterSettings(settings: Partial<PrinterSettings>): Promise<void> {
  await window.go.main.PrinterService.SaveSettings(settings as PrinterSettings)
}

export async function printInvoice(data: ReceiptData): Promise<{ job: PrintJob; html: string }> {
  const [job, html] = await window.go.main.PrinterService.PrintInvoice(data)
  return { job, html }
}

export async function printReceipt(data: ReceiptData): Promise<{ job: PrintJob; content: any }> {
  const [job, content] = await window.go.main.PrinterService.PrintReceipt(data)
  return { job, content }
}

export async function printTest(): Promise<{ job: PrintJob; content: any }> {
  const [job, content] = await window.go.main.PrinterService.PrintTest()
  return { job, content }
}

export async function listPrintJobs(page = 1, perPage = 20) {
  const offset = (page - 1) * perPage
  const [jobs, total] = await window.go.main.PrinterService.ListPrintJobs(perPage, offset)
  return { jobs: jobs || [], meta: { page, per_page: perPage, total, total_pages: Math.ceil(total / perPage) } }
}

export async function getPrintJob(id: string): Promise<PrintJob> {
  return window.go.main.PrinterService.GetPrintJob(id)
}
