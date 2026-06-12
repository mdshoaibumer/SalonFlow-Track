import type {
  LoginInput,
  LoginOutput,
  SessionInfo,
  User,
  Role,
  Permission,
  AuditFilter,
  AuditLogListOutput,
  CreateUserInput,
  UpdateUserInput,
  DiagnosticsInfo,
} from '@/types/auth'

// Auth service - calls Wails backend bindings
export async function login(input: LoginInput): Promise<LoginOutput> {
  return window.go.main.AuthService.Login(input)
}

export async function logout(): Promise<void> {
  return window.go.main.AuthService.Logout()
}

export async function getCurrentSession(): Promise<SessionInfo> {
  return window.go.main.AuthService.GetCurrentSession()
}

export async function isAuthenticated(): Promise<boolean> {
  return window.go.main.AuthService.IsAuthenticated()
}

export async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
  return window.go.main.AuthService.ChangePassword(oldPassword, newPassword)
}

export async function hasPermission(permission: string): Promise<boolean> {
  return window.go.main.AuthService.HasPermission(permission)
}

export async function hasAnyPermission(permissions: string[]): Promise<boolean> {
  return window.go.main.AuthService.HasAnyPermission(permissions)
}

// Token management for Remember Me
export async function getToken(): Promise<string> {
  return window.go.main.AuthService.GetToken()
}

export async function setToken(token: string): Promise<void> {
  return window.go.main.AuthService.SetToken(token)
}

// User management
export async function getUsers(): Promise<User[]> {
  return window.go.main.AuthService.GetUsers()
}

export async function getUser(id: string): Promise<User> {
  return window.go.main.AuthService.GetUser(id)
}

export async function createUser(input: CreateUserInput): Promise<User> {
  return window.go.main.AuthService.CreateUser(input)
}

export async function updateUser(input: UpdateUserInput): Promise<User> {
  return window.go.main.AuthService.UpdateUser(input)
}

export async function deleteUser(id: string): Promise<void> {
  return window.go.main.AuthService.DeleteUser(id)
}

export async function resetUserPassword(userId: string, newPassword: string): Promise<void> {
  return window.go.main.AuthService.ResetUserPassword(userId, newPassword)
}

export async function assignUserRole(userId: string, roleId: string): Promise<void> {
  return window.go.main.AuthService.AssignUserRole(userId, roleId)
}

export async function removeUserRole(userId: string, roleId: string): Promise<void> {
  return window.go.main.AuthService.RemoveUserRole(userId, roleId)
}

// Roles & Permissions
export async function getRoles(): Promise<Role[]> {
  return window.go.main.AuthService.GetRoles()
}

export async function getPermissions(): Promise<Permission[]> {
  return window.go.main.AuthService.GetPermissions()
}

export async function getRolePermissions(roleId: string): Promise<Permission[]> {
  return window.go.main.AuthService.GetRolePermissions(roleId)
}

// Audit logs
export async function getAuditLogs(filter: AuditFilter): Promise<AuditLogListOutput> {
  return window.go.main.AuthService.GetAuditLogs(filter)
}

// Diagnostics
export async function getDiagnostics(): Promise<DiagnosticsInfo> {
  return window.go.main.DiagnosticsService.GetDiagnostics()
}

export async function exportDiagnosticsBundle(): Promise<string> {
  return window.go.main.DiagnosticsService.ExportDiagnosticsBundle()
}
