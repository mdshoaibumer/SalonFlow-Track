import { z } from 'zod'

export const staffFormSchema = z.object({
  full_name: z.string().min(1, 'Name is required').max(100, 'Name is too long'),
  phone: z
    .string()
    .min(1, 'Phone is required')
    .regex(/^[6-9]\d{9}$/, 'Must be a valid 10-digit Indian mobile number'),
  email: z.string().email('Invalid email').or(z.literal('')).optional(),
  gender: z.enum(['male', 'female', 'other']),
  designation: z.enum(['stylist', 'assistant', 'receptionist', 'manager']),
  joining_date: z.string().min(1, 'Joining date is required'),
  base_salary: z.number().min(0, 'Salary must be 0 or more'),
  commission_percentage: z
    .number()
    .min(0, 'Commission must be 0 or more')
    .max(100, 'Commission cannot exceed 100'),
})

export const updateStaffFormSchema = staffFormSchema.extend({
  status: z.enum(['active', 'inactive']),
})

export type StaffFormValues = z.infer<typeof staffFormSchema>
export type UpdateStaffFormValues = z.infer<typeof updateStaffFormSchema>
