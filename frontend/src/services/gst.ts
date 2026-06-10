import type { GSTSettings, TaxRate, GSTReport } from '@/types'

export async function getGSTSettings(): Promise<GSTSettings> {
  return window.go.main.GSTService.GetSettings()
}

export async function saveGSTSettings(settings: Partial<GSTSettings>): Promise<void> {
  await window.go.main.GSTService.SaveSettings(settings as GSTSettings)
}

export async function listTaxRates(category = ''): Promise<TaxRate[]> {
  return window.go.main.GSTService.ListTaxRates(category)
}

export async function createTaxRate(rate: Partial<TaxRate>): Promise<void> {
  await window.go.main.GSTService.CreateTaxRate(rate as TaxRate)
}

export async function updateTaxRate(id: string, rate: Partial<TaxRate>): Promise<void> {
  await window.go.main.GSTService.UpdateTaxRate({ ...rate, id } as TaxRate)
}

export async function deleteTaxRate(id: string): Promise<void> {
  await window.go.main.GSTService.DeleteTaxRate(id)
}

export async function getGSTReport(startDate: string, endDate: string, period = 'daily'): Promise<GSTReport> {
  return window.go.main.GSTService.GetReport({ date_from: startDate, date_to: endDate, report_type: period })
}
