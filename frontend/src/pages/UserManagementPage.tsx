import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Users, Plus, Shield, Lock, Unlock, Trash2, KeyRound } from 'lucide-react'
import { toast } from 'sonner'
import { useAuth } from '@/app/providers/AuthProvider'
import * as authService from '@/services/auth'
import type { User, Role, CreateUserInput } from '@/types/auth'

export function UserManagementPage() {
  const { hasPermission } = useAuth()
  const queryClient = useQueryClient()
  const [showCreateForm, setShowCreateForm] = useState(false)

  const { data: users = [], isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: authService.getUsers,
  })

  const { data: roles = [] } = useQuery({
    queryKey: ['roles'],
    queryFn: authService.getRoles,
  })

  const deleteMutation = useMutation({
    mutationFn: authService.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast.success('User deleted')
    },
    onError: () => toast.error('Failed to delete user'),
  })

  if (!hasPermission('users.view')) {
    return <div className="p-6 text-muted-foreground">You do not have permission to view this page.</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Users className="h-6 w-6 text-violet-600" />
          <h1 className="text-xl font-bold">User Management</h1>
        </div>
        {hasPermission('users.manage') && (
          <button
            onClick={() => setShowCreateForm(true)}
            className="flex items-center gap-2 rounded-lg bg-violet-600 px-4 py-2 text-sm font-medium text-white hover:bg-violet-700 transition-colors"
          >
            <Plus className="h-4 w-4" />
            Add User
          </button>
        )}
      </div>

      {showCreateForm && (
        <CreateUserForm
          roles={roles}
          onClose={() => setShowCreateForm(false)}
          onSuccess={() => {
            setShowCreateForm(false)
            queryClient.invalidateQueries({ queryKey: ['users'] })
          }}
        />
      )}

      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading users...</div>
      ) : (
        <div className="rounded-lg border border-border overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">User</th>
                <th className="px-4 py-3 text-left font-medium">Roles</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Last Login</th>
                <th className="px-4 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {users.map((user) => (
                <UserRow
                  key={user.id}
                  user={user}
                  canManage={hasPermission('users.manage')}
                  onDelete={() => deleteMutation.mutate(user.id)}
                />
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function UserRow({ user, canManage, onDelete }: { user: User; canManage: boolean; onDelete: () => void }) {
  const [showResetPassword, setShowResetPassword] = useState(false)

  return (
    <>
      <tr className="hover:bg-muted/30">
        <td className="px-4 py-3">
          <div>
            <p className="font-medium">{user.display_name}</p>
            <p className="text-xs text-muted-foreground">@{user.username} • {user.email}</p>
          </div>
        </td>
        <td className="px-4 py-3">
          <div className="flex flex-wrap gap-1">
            {user.roles?.map((role) => (
              <span key={role.id} className="inline-flex items-center gap-1 rounded-full bg-violet-100 px-2 py-0.5 text-xs font-medium text-violet-700 dark:bg-violet-900/30 dark:text-violet-300">
                <Shield className="h-3 w-3" />
                {role.name}
              </span>
            ))}
          </div>
        </td>
        <td className="px-4 py-3">
          {user.is_locked ? (
            <span className="inline-flex items-center gap-1 text-xs text-red-600">
              <Lock className="h-3 w-3" /> Locked
            </span>
          ) : user.is_active ? (
            <span className="inline-flex items-center gap-1 text-xs text-emerald-600">
              <Unlock className="h-3 w-3" /> Active
            </span>
          ) : (
            <span className="text-xs text-muted-foreground">Inactive</span>
          )}
        </td>
        <td className="px-4 py-3 text-xs text-muted-foreground">
          {user.last_login_at ? new Date(user.last_login_at).toLocaleDateString() : 'Never'}
        </td>
        <td className="px-4 py-3 text-right">
          {canManage && (
            <div className="flex items-center justify-end gap-1">
              <button
                onClick={() => setShowResetPassword(true)}
                className="rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                title="Reset Password"
              >
                <KeyRound className="h-4 w-4" />
              </button>
              <button
                onClick={onDelete}
                className="rounded p-1.5 text-muted-foreground hover:bg-red-100 hover:text-red-600"
                title="Delete User"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          )}
        </td>
      </tr>
      {showResetPassword && (
        <tr>
          <td colSpan={5} className="px-4 py-3 bg-muted/30">
            <ResetPasswordForm userId={user.id} username={user.username} onClose={() => setShowResetPassword(false)} />
          </td>
        </tr>
      )}
    </>
  )
}

function CreateUserForm({ roles, onClose, onSuccess }: { roles: Role[]; onClose: () => void; onSuccess: () => void }) {
  const [form, setForm] = useState<CreateUserInput>({
    username: '', email: '', password: '', display_name: '', phone: '', role_id: 'role-staff',
  })

  const mutation = useMutation({
    mutationFn: () => authService.createUser(form),
    onSuccess: () => {
      toast.success('User created successfully')
      onSuccess()
    },
    onError: (err: any) => toast.error(err?.message || 'Failed to create user'),
  })

  return (
    <div className="rounded-lg border border-border bg-background p-4 space-y-4">
      <h3 className="font-medium">Create New User</h3>
      <div className="grid grid-cols-2 gap-3">
        <input
          placeholder="Username"
          value={form.username}
          onChange={(e) => setForm({ ...form, username: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        />
        <input
          placeholder="Display Name"
          value={form.display_name}
          onChange={(e) => setForm({ ...form, display_name: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        />
        <input
          placeholder="Email"
          type="email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        />
        <input
          placeholder="Phone"
          value={form.phone}
          onChange={(e) => setForm({ ...form, phone: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        />
        <input
          placeholder="Password (min 8, upper+lower+digit+special)"
          type="password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        />
        <select
          value={form.role_id}
          onChange={(e) => setForm({ ...form, role_id: e.target.value })}
          className="rounded-lg border border-border px-3 py-2 text-sm"
        >
          {roles.map((r) => (
            <option key={r.id} value={r.id}>{r.name}</option>
          ))}
        </select>
      </div>
      <div className="flex gap-2 justify-end">
        <button onClick={onClose} className="px-3 py-1.5 text-sm rounded-lg border border-border hover:bg-muted">
          Cancel
        </button>
        <button
          onClick={() => mutation.mutate()}
          disabled={mutation.isPending}
          className="px-3 py-1.5 text-sm rounded-lg bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-50"
        >
          {mutation.isPending ? 'Creating...' : 'Create User'}
        </button>
      </div>
    </div>
  )
}

function ResetPasswordForm({ userId, username, onClose }: { userId: string; username: string; onClose: () => void }) {
  const [password, setPassword] = useState('')

  const mutation = useMutation({
    mutationFn: () => authService.resetUserPassword(userId, password),
    onSuccess: () => {
      toast.success(`Password reset for ${username}`)
      onClose()
    },
    onError: (err: any) => toast.error(err?.message || 'Failed to reset password'),
  })

  return (
    <div className="flex items-center gap-3">
      <span className="text-sm">Reset password for <strong>{username}</strong>:</span>
      <input
        type="password"
        placeholder="New password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        className="rounded-lg border border-border px-3 py-1.5 text-sm w-48"
      />
      <button
        onClick={() => mutation.mutate()}
        disabled={mutation.isPending || !password}
        className="px-3 py-1.5 text-sm rounded-lg bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-50"
      >
        Reset
      </button>
      <button onClick={onClose} className="text-sm text-muted-foreground hover:text-foreground">Cancel</button>
    </div>
  )
}
