import type { SalaryRecord, SalaryCycle, Advance, SalaryStats, GenerateSalaryInput, CreateAdvanceInput } from '@/types'

export interface ListAdvancesParams {
  staff_id?: string
  status?: string
  page?: number
  per_page?: number
}

export async function generateSalary(input: GenerateSalaryInput) {
  return window.go.main.SalaryService.GenerateSalary(input)
}

export async function getSalaryById(id: string): Promise<SalaryRecord> {
  return window.go.main.SalaryService.GetSalary(id)
}

export async function listSalaries(month: number, year: number): Promise<SalaryRecord[]> {
  return window.go.main.SalaryService.ListSalaries({ month, year })
}

export async function paySalary(id: string): Promise<void> {
  await window.go.main.SalaryService.PaySalary(id)
}

export async function listSalaryCycles(year?: number): Promise<SalaryCycle[]> {
  return window.go.main.SalaryService.ListCycles(year || new Date().getFullYear())
}

export async function createAdvance(input: CreateAdvanceInput): Promise<Advance> {
  return window.go.main.SalaryService.CreateAdvance(input)
}

export async function approveAdvance(id: string): Promise<Advance> {
  return window.go.main.SalaryService.ApproveAdvance(id)
}

export async function rejectAdvance(id: string): Promise<Advance> {
  return window.go.main.SalaryService.RejectAdvance(id)
}

export async function listAdvances(params: ListAdvancesParams = {}) {
  return window.go.main.SalaryService.ListAdvances({
    staff_id: params.staff_id || '',
    status: params.status || '',
    page: params.page || 1,
    per_page: params.per_page || 20,
  })
}

export async function getSalaryStats(): Promise<SalaryStats> {
  return window.go.main.SalaryService.GetSalaryStats()
}
