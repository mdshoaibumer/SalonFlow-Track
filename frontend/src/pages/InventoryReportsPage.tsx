import { useInventoryAnalyticsReport } from '@/hooks/useAnalytics'
import { Package, AlertTriangle, TrendingUp, TrendingDown } from 'lucide-react'
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer,
} from 'recharts'

export function InventoryReportsPage() {
  const { data, isLoading } = useInventoryAnalyticsReport()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Inventory Reports</h1>
        <p className="text-muted-foreground">Stock movement analysis and inventory insights</p>
      </div>

      {isLoading && <p className="text-muted-foreground">Loading...</p>}

      {data && (
        <>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="rounded-lg border bg-card p-6 flex items-center gap-4">
              <Package className="h-8 w-8 text-blue-600" />
              <div>
                <p className="text-sm text-muted-foreground">Total Inventory Value</p>
                <p className="text-2xl font-bold">₹{Math.round(data.total_value).toLocaleString()}</p>
              </div>
            </div>
            <div className="rounded-lg border bg-card p-6 flex items-center gap-4">
              <AlertTriangle className="h-8 w-8 text-red-500" />
              <div>
                <p className="text-sm text-muted-foreground">Low Stock Products</p>
                <p className="text-2xl font-bold">{data.low_stock_count}</p>
              </div>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-card p-6">
              <div className="flex items-center gap-2 mb-4">
                <TrendingUp className="h-4 w-4 text-green-600" />
                <h3 className="text-lg font-semibold">Fast Moving Products</h3>
              </div>
              {data.fast_moving && data.fast_moving.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={data.fast_moving} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" tick={{ fontSize: 10 }} />
                    <YAxis type="category" dataKey="name" tick={{ fontSize: 10 }} width={120} />
                    <Tooltip />
                    <Bar dataKey="value" fill="#10b981" name="Units Consumed" />
                  </BarChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <div className="flex items-center gap-2 mb-4">
                <TrendingDown className="h-4 w-4 text-orange-500" />
                <h3 className="text-lg font-semibold">Slow Moving Products</h3>
              </div>
              {data.slow_moving && data.slow_moving.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={data.slow_moving} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" tick={{ fontSize: 10 }} />
                    <YAxis type="category" dataKey="name" tick={{ fontSize: 10 }} width={120} />
                    <Tooltip />
                    <Bar dataKey="value" fill="#f59e0b" name="Units Consumed" />
                  </BarChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Purchase Trend (Monthly)</h3>
              {data.purchase_trend && data.purchase_trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={data.purchase_trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`₹${Number(value).toLocaleString()}`]} />
                    <Line type="monotone" dataKey="value" stroke="#3b82f6" strokeWidth={2} name="Purchases" />
                  </LineChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>

            <div className="rounded-lg border bg-card p-6">
              <h3 className="text-lg font-semibold mb-4">Consumption Trend (Monthly)</h3>
              {data.consumption_trend && data.consumption_trend.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={data.consumption_trend}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" tick={{ fontSize: 11 }} />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip />
                    <Line type="monotone" dataKey="value" stroke="#ef4444" strokeWidth={2} name="Units Consumed" />
                  </LineChart>
                </ResponsiveContainer>
              ) : <p className="text-center text-muted-foreground py-8">No data</p>}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
