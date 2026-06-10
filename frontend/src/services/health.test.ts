import { describe, it, expect } from 'vitest'
import { getHealthStatus } from './health'

describe('Health Service', () => {
  it('gets health status', async () => {
    const status = await getHealthStatus()
    expect(status.status).toBe('healthy')
    expect(status.version).toBe('0.2.0')
    expect(status.environment).toBe('production')
    expect(status.database).toBe('SQLite')
  })
})
