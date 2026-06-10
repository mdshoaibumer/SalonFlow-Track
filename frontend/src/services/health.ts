export interface HealthStatus {
  status: string
  version: string
  environment: string
  database: string
  uptime: string
  go_version: string
}

export async function getHealthStatus(): Promise<HealthStatus> {
  const version = await window.go.main.App.GetVersion()
  const environment = await window.go.main.App.GetEnvironment()
  return {
    status: 'healthy',
    version,
    environment,
    database: 'SQLite',
    uptime: '-',
    go_version: '-',
  }
}
