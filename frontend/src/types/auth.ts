// Auth types
export interface LoginInput {
  username: string
  password: string
  remember_me: boolean
  device_id: string
}

export interface LoginOutput {
  token: string
  user: SessionInfo
  expires_at: string
}

export interface SessionInfo {
  user_id: string
  username: string
  display_name: string
  email: string
  roles: string[]
  permissions: string[]
}

export interface User {
  id: string
  username: string
  email: string
  display_name: string
  phone: string
  is_active: boolean
  is_locked: boolean
  failed_attempts: number
  locked_until?: string
  last_login_at?: string
  password_changed_at: string
  must_change_password: boolean
  created_at: string
  updated_at: string
  roles?: Role[]
  permissions?: string[]
}

export interface Role {
  id: string
  name: string
  description: string
  is_system: boolean
  created_at: string
  updated_at: string
}

export interface Permission {
  id: string
  code: string
  module: string
  action: string
  description: string
  created_at: string
}

export interface AuditLog {
  id: string
  timestamp: string
  user_id: string
  username: string
  action: string
  module: string
  entity_type: string
  entity_id: string
  description: string
  old_value: string
  new_value: string
  device_id: string
  ip_address: string
  user_agent: string
  app_version: string
  severity: string
  created_at: string
}

export interface AuditFilter {
  user_id?: string
  module?: string
  action?: string
  entity_type?: string
  entity_id?: string
  severity?: string
  from_date?: string
  to_date?: string
  page?: number
  per_page?: number
}

export interface AuditLogListOutput {
  logs: AuditLog[]
  total: number
  page: number
}

export interface DiagnosticsInfo {
  app_version: string
  go_version: string
  os: string
  arch: string
  database_path: string
  database_size_bytes: number
  db_version: string
  log_directory: string
  num_cpu: number
  num_goroutine: number
  mem_alloc_mb: number
  mem_total_alloc_mb: number
  uptime: string
  last_backup: string
  total_users: number
  total_invoices: number
  total_customers: number
}

export interface CreateUserInput {
  username: string
  email: string
  password: string
  display_name: string
  phone: string
  role_id: string
}

export interface UpdateUserInput {
  id: string
  email: string
  display_name: string
  phone: string
  is_active: boolean
}
