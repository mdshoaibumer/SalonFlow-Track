import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useCustomerList } from '@/hooks/useCustomers'
import { useServiceList } from '@/hooks/useServices'
import { useCreateInvoice } from '@/hooks/useInvoices'
import { Search, Plus, Trash2, Receipt } from 'lucide-react'
import type { Customer, Service, CreateInvoiceInput, PaymentMethod } from '@/types'

interface BillingItem {
  service: Service
  quantity: number
  discount: number
}

export function BillingPage() {
  const [customerSearch, setCustomerSearch] = useState('')
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null)
  const [selectedStaffId, setSelectedStaffId] = useState('')
  const [items, setItems] = useState<BillingItem[]>([])
  const [discount, setDiscount] = useState(0)
  const [tax, setTax] = useState(0)
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('cash')
  const [notes, setNotes] = useState('')

  const { data: customerData } = useCustomerList({ search: customerSearch || undefined, per_page: 5 })
  const { data: serviceData } = useServiceList({ status: 'active', per_page: 100 })
  const createInvoice = useCreateInvoice()

  const subtotal = items.reduce((sum, item) => sum + (item.service.price * item.quantity - item.discount), 0)
  const grandTotal = subtotal - discount + tax

  const addService = (service: Service) => {
    const existing = items.find((i) => i.service.id === service.id)
    if (existing) {
      setItems(items.map((i) => i.service.id === service.id ? { ...i, quantity: i.quantity + 1 } : i))
    } else {
      setItems([...items, { service, quantity: 1, discount: 0 }])
    }
  }

  const removeItem = (serviceId: string) => {
    setItems(items.filter((i) => i.service.id !== serviceId))
  }

  const handleGenerateInvoice = () => {
    if (!selectedCustomer || !selectedStaffId || items.length === 0) return

    const input: CreateInvoiceInput = {
      customer_id: selectedCustomer.id,
      staff_id: selectedStaffId,
      items: items.map((item) => ({
        service_id: item.service.id,
        quantity: item.quantity,
        discount: item.discount,
      })),
      discount,
      tax,
      payment_method: paymentMethod,
      notes,
    }

    createInvoice.mutate(input, {
      onSuccess: () => {
        setSelectedCustomer(null)
        setItems([])
        setDiscount(0)
        setTax(0)
        setNotes('')
        setCustomerSearch('')
      },
    })
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Billing</h1>
        <p className="text-muted-foreground">Create a new invoice</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left: Customer & Services */}
        <div className="lg:col-span-2 space-y-4">
          {/* Customer Search */}
          <div className="rounded-lg border bg-card p-4 space-y-3">
            <h3 className="font-semibold">Customer</h3>
            {selectedCustomer ? (
              <div className="flex items-center justify-between p-3 rounded border bg-muted/50">
                <div>
                  <p className="font-medium">{selectedCustomer.full_name}</p>
                  <p className="text-sm text-muted-foreground">{selectedCustomer.phone}</p>
                </div>
                <Button variant="ghost" size="sm" onClick={() => setSelectedCustomer(null)}>
                  Change
                </Button>
              </div>
            ) : (
              <>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search customer by name or phone..."
                    className="pl-9"
                    value={customerSearch}
                    onChange={(e) => setCustomerSearch(e.target.value)}
                  />
                </div>
                {customerData && customerData.customers.length > 0 && customerSearch && (
                  <div className="border rounded-md divide-y max-h-40 overflow-y-auto">
                    {customerData.customers.map((c) => (
                      <button
                        key={c.id}
                        className="w-full p-2 text-left hover:bg-muted/50 text-sm"
                        onClick={() => { setSelectedCustomer(c); setCustomerSearch('') }}
                      >
                        <span className="font-medium">{c.full_name}</span>
                        <span className="text-muted-foreground ml-2">{c.phone}</span>
                      </button>
                    ))}
                  </div>
                )}
              </>
            )}
          </div>

          {/* Staff ID */}
          <div className="rounded-lg border bg-card p-4 space-y-3">
            <h3 className="font-semibold">Staff</h3>
            <Input
              placeholder="Enter staff ID"
              value={selectedStaffId}
              onChange={(e) => setSelectedStaffId(e.target.value)}
            />
          </div>

          {/* Services */}
          <div className="rounded-lg border bg-card p-4 space-y-3">
            <h3 className="font-semibold">Services</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              {serviceData?.services.map((svc) => (
                <Button
                  key={svc.id}
                  variant="outline"
                  size="sm"
                  className="justify-start"
                  onClick={() => addService(svc)}
                >
                  <Plus className="h-3 w-3 mr-1" />
                  {svc.name} - ₹{svc.price}
                </Button>
              ))}
            </div>
          </div>
        </div>

        {/* Right: Invoice Summary */}
        <div className="space-y-4">
          <div className="rounded-lg border bg-card p-4 space-y-4">
            <h3 className="font-semibold">Invoice Summary</h3>

            {items.length === 0 ? (
              <p className="text-sm text-muted-foreground">Add services to begin</p>
            ) : (
              <div className="space-y-2">
                {items.map((item) => (
                  <div key={item.service.id} className="flex items-center justify-between text-sm">
                    <div>
                      <span>{item.service.name}</span>
                      <span className="text-muted-foreground ml-1">x{item.quantity}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span>₹{(item.service.price * item.quantity).toLocaleString('en-IN')}</span>
                      <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => removeItem(item.service.id)}>
                        <Trash2 className="h-3 w-3 text-destructive" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <div className="border-t pt-3 space-y-2 text-sm">
              <div className="flex justify-between">
                <span>Subtotal</span>
                <span>₹{subtotal.toLocaleString('en-IN')}</span>
              </div>
              <div className="flex items-center justify-between">
                <span>Discount</span>
                <Input
                  type="number"
                  className="w-24 h-7 text-right"
                  value={discount}
                  onChange={(e) => setDiscount(Number(e.target.value))}
                  min={0}
                />
              </div>
              <div className="flex items-center justify-between">
                <span>Tax</span>
                <Input
                  type="number"
                  className="w-24 h-7 text-right"
                  value={tax}
                  onChange={(e) => setTax(Number(e.target.value))}
                  min={0}
                />
              </div>
              <div className="flex justify-between font-bold text-base border-t pt-2">
                <span>Grand Total</span>
                <span>₹{grandTotal.toLocaleString('en-IN')}</span>
              </div>
            </div>

            <div className="space-y-3 pt-2">
              <Select value={paymentMethod} onValueChange={(v) => setPaymentMethod(v as PaymentMethod)}>
                <SelectTrigger>
                  <SelectValue placeholder="Payment Method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cash">Cash</SelectItem>
                  <SelectItem value="upi">UPI</SelectItem>
                  <SelectItem value="card">Card</SelectItem>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                </SelectContent>
              </Select>

              <Input
                placeholder="Notes (optional)"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
              />

              <Button
                className="w-full"
                size="lg"
                disabled={!selectedCustomer || !selectedStaffId || items.length === 0 || createInvoice.isPending}
                onClick={handleGenerateInvoice}
              >
                <Receipt className="mr-2 h-4 w-4" />
                {createInvoice.isPending ? 'Generating...' : 'Generate Invoice'}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
