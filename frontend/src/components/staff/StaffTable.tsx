import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Pencil, Trash2 } from 'lucide-react'
import type { Staff } from '@/types'

interface StaffTableProps {
  staff: Staff[]
  onEdit: (staff: Staff) => void
  onDelete: (staff: Staff) => void
}

export function StaffTable({ staff, onEdit, onDelete }: StaffTableProps) {
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0,
    }).format(amount)
  }

  const designationLabel = (d: string) => {
    const map: Record<string, string> = {
      stylist: 'Stylist',
      assistant: 'Assistant',
      receptionist: 'Receptionist',
      manager: 'Manager',
    }
    return map[d] || d
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Staff Code</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Phone</TableHead>
            <TableHead>Designation</TableHead>
            <TableHead className="text-right">Base Salary</TableHead>
            <TableHead className="text-right">Commission %</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {staff.length === 0 ? (
            <TableRow>
              <TableCell colSpan={8} className="h-24 text-center text-muted-foreground">
                No staff members found.
              </TableCell>
            </TableRow>
          ) : (
            staff.map((member) => (
              <TableRow key={member.id}>
                <TableCell className="font-mono text-sm">{member.staff_code}</TableCell>
                <TableCell className="font-medium">{member.full_name}</TableCell>
                <TableCell>{member.phone}</TableCell>
                <TableCell>{designationLabel(member.designation)}</TableCell>
                <TableCell className="text-right">{formatCurrency(member.base_salary)}</TableCell>
                <TableCell className="text-right">{member.commission_percentage}%</TableCell>
                <TableCell>
                  <Badge variant={member.status === 'active' ? 'default' : 'secondary'}>
                    {member.status}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => onEdit(member)}
                      aria-label={`Edit ${member.full_name}`}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => onDelete(member)}
                      aria-label={`Delete ${member.full_name}`}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
