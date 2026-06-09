import { useState } from 'react'
import { useAppointments, useCreateAppointment, useUpdateAppointmentStatus, useDeleteAppointment } from '@/hooks/useAppointments'
import type { Appointment, AppointmentFilter, AppointmentStatus } from '@/types'

const STATUS_COLORS: Record<string, string> = {
  booked: 'bg-blue-100 text-blue-800',
  confirmed: 'bg-indigo-100 text-indigo-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
  no_show: 'bg-gray-100 text-gray-800',
}

export function AppointmentsPage() {
  const [tab, setTab] = useState<'calendar' | 'list' | 'create'>('list')
  const [filter, setFilter] = useState<AppointmentFilter>({})

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Appointments</h1>
        <p className="text-muted-foreground">Manage bookings and schedules</p>
      </div>

      <div className="flex gap-2 border-b pb-2">
        <button onClick={() => setTab('list')} className={`px-4 py-2 rounded-t ${tab === 'list' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>List View</button>
        <button onClick={() => setTab('calendar')} className={`px-4 py-2 rounded-t ${tab === 'calendar' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>Calendar</button>
        <button onClick={() => setTab('create')} className={`px-4 py-2 rounded-t ${tab === 'create' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>New Appointment</button>
      </div>

      {tab === 'list' && <AppointmentListTab filter={filter} setFilter={setFilter} />}
      {tab === 'calendar' && <CalendarTab />}
      {tab === 'create' && <CreateAppointmentTab onDone={() => setTab('list')} />}
    </div>
  )
}

function AppointmentListTab({ filter, setFilter }: { filter: AppointmentFilter; setFilter: (f: AppointmentFilter) => void }) {
  const { data: appointments, isLoading } = useAppointments(filter)
  const statusMutation = useUpdateAppointmentStatus()
  const deleteMutation = useDeleteAppointment()

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <select className="border rounded px-3 py-1 text-sm" value={filter.status || ''} onChange={e => setFilter({ ...filter, status: e.target.value as AppointmentStatus || undefined })}>
          <option value="">All Statuses</option>
          <option value="booked">Booked</option>
          <option value="confirmed">Confirmed</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
          <option value="no_show">No Show</option>
        </select>
        <input type="date" className="border rounded px-3 py-1 text-sm" value={filter.start_date || ''} onChange={e => setFilter({ ...filter, start_date: e.target.value })} />
        <input type="date" className="border rounded px-3 py-1 text-sm" value={filter.end_date || ''} onChange={e => setFilter({ ...filter, end_date: e.target.value })} />
      </div>

      {isLoading && <p className="text-muted-foreground">Loading appointments...</p>}

      <div className="border rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-muted">
            <tr>
              <th className="px-4 py-2 text-left">Customer</th>
              <th className="px-4 py-2 text-left">Staff</th>
              <th className="px-4 py-2 text-left">Start</th>
              <th className="px-4 py-2 text-left">End</th>
              <th className="px-4 py-2 text-left">Status</th>
              <th className="px-4 py-2 text-left">Amount</th>
              <th className="px-4 py-2 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {(appointments || []).map((appt: Appointment) => (
              <tr key={appt.id} className="border-t">
                <td className="px-4 py-2">{appt.customer_name || appt.customer_id}</td>
                <td className="px-4 py-2">{appt.staff_name || appt.staff_id}</td>
                <td className="px-4 py-2">{new Date(appt.start_time).toLocaleString()}</td>
                <td className="px-4 py-2">{new Date(appt.end_time).toLocaleString()}</td>
                <td className="px-4 py-2"><span className={`px-2 py-0.5 rounded text-xs ${STATUS_COLORS[appt.status] || ''}`}>{appt.status}</span></td>
                <td className="px-4 py-2">₹{appt.total_amount}</td>
                <td className="px-4 py-2 flex gap-1">
                  {appt.status === 'booked' && <button className="text-xs px-2 py-1 bg-green-100 rounded" onClick={() => statusMutation.mutate({ id: appt.id, status: 'confirmed' })}>Confirm</button>}
                  {appt.status === 'confirmed' && <button className="text-xs px-2 py-1 bg-yellow-100 rounded" onClick={() => statusMutation.mutate({ id: appt.id, status: 'in_progress' })}>Start</button>}
                  {appt.status === 'in_progress' && <button className="text-xs px-2 py-1 bg-green-100 rounded" onClick={() => statusMutation.mutate({ id: appt.id, status: 'completed' })}>Complete</button>}
                  <button className="text-xs px-2 py-1 bg-red-100 rounded" onClick={() => deleteMutation.mutate(appt.id)}>Delete</button>
                </td>
              </tr>
            ))}
            {(!appointments || appointments.length === 0) && <tr><td colSpan={7} className="px-4 py-8 text-center text-muted-foreground">No appointments found</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function CalendarTab() {
  // FullCalendar integration point - for now show placeholder with integration note
  return (
    <div className="border rounded-lg p-8 text-center text-muted-foreground">
      <p className="text-lg font-medium">Calendar View</p>
      <p className="mt-2">FullCalendar integration ready. Install @fullcalendar/react for interactive calendar.</p>
      <p className="mt-1 text-sm">npm install @fullcalendar/react @fullcalendar/daygrid @fullcalendar/timegrid @fullcalendar/interaction</p>
    </div>
  )
}

function CreateAppointmentTab({ onDone }: { onDone: () => void }) {
  const createMutation = useCreateAppointment()
  const [form, setForm] = useState({
    customer_id: '',
    staff_id: '',
    start_time: '',
    end_time: '',
    notes: '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(form, { onSuccess: onDone })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-lg">
      <div>
        <label className="text-sm font-medium">Customer ID</label>
        <input className="w-full border rounded px-3 py-2 mt-1" value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} required />
      </div>
      <div>
        <label className="text-sm font-medium">Staff ID</label>
        <input className="w-full border rounded px-3 py-2 mt-1" value={form.staff_id} onChange={e => setForm({ ...form, staff_id: e.target.value })} required />
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-sm font-medium">Start Time</label>
          <input type="datetime-local" className="w-full border rounded px-3 py-2 mt-1" value={form.start_time} onChange={e => setForm({ ...form, start_time: e.target.value })} required />
        </div>
        <div>
          <label className="text-sm font-medium">End Time</label>
          <input type="datetime-local" className="w-full border rounded px-3 py-2 mt-1" value={form.end_time} onChange={e => setForm({ ...form, end_time: e.target.value })} required />
        </div>
      </div>
      <div>
        <label className="text-sm font-medium">Notes</label>
        <textarea className="w-full border rounded px-3 py-2 mt-1" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} rows={3} />
      </div>
      <button type="submit" disabled={createMutation.isPending} className="px-4 py-2 bg-primary text-primary-foreground rounded">
        {createMutation.isPending ? 'Creating...' : 'Create Appointment'}
      </button>
    </form>
  )
}
