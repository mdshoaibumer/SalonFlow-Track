import { useState, useEffect } from 'react'
import { toastSuccess } from '@/lib/toast'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/services/api-client'
import { getHealthStatus } from '@/services/health'
import { PageHeader } from '@/components/shared/PageHeader'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LoadingState } from '@/components/shared/LoadingState'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Settings, Store, Bell, Database, Shield, Check } from 'lucide-react'
import { useTheme } from '@/app/providers/ThemeProvider'
import type { Setting } from '@/types'

const SETTINGS_STORAGE_KEY = 'salonflow-settings'

function loadSavedSettings(): Record<string, string> {
  try {
    const raw = localStorage.getItem(SETTINGS_STORAGE_KEY)
    return raw ? JSON.parse(raw) : {}
  } catch {
    return {}
  }
}

export function SettingsPage() {
  const [activeTab, setActiveTab] = useState('general')
  const { theme, setTheme } = useTheme()
  const [saved, setSaved] = useState(false)
  const [salonSettings, setSalonSettings] = useState<Record<string, string>>(loadSavedSettings)

  const { data: health } = useQuery({
    queryKey: ['health'],
    queryFn: getHealthStatus,
    retry: 1,
  })

  const { data: settings, isLoading } = useQuery({
    queryKey: ['settings'],
    queryFn: async () => {
      const res = await apiClient.get<Setting[]>('/settings')
      return res.data || []
    },
  })

  // Initialize form from fetched settings (only once)
  useEffect(() => {
    if (settings && settings.length > 0) {
      const saved = loadSavedSettings()
      if (Object.keys(saved).length === 0) {
        const initial: Record<string, string> = {}
        for (const s of settings) {
          initial[s.key] = s.value
        }
        setSalonSettings(initial)
      }
    }
  }, [settings])

  const updateSetting = (key: string, value: string) => {
    setSalonSettings(prev => ({ ...prev, [key]: value }))
  }

  const handleSave = () => {
    localStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(salonSettings))
    toastSuccess('Settings saved')
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <PageHeader title="Settings" description="Application and salon configuration" />
        <LoadingState variant="page" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Settings"
        description="Application and salon configuration"
        actions={
          <Button onClick={handleSave} size="sm">
            {saved ? <><Check className="mr-2 h-4 w-4" /> Saved</> : 'Save Changes'}
          </Button>
        }
      />

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="salon">Salon</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
          <TabsTrigger value="system">System</TabsTrigger>
        </TabsList>

        <TabsContent value="general">
          <div className="grid gap-6 mt-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-base flex items-center gap-2">
                  <Settings className="h-4 w-4" />
                  Appearance
                </CardTitle>
                <CardDescription>Customize the look and feel</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Theme</label>
                    <Select value={theme} onValueChange={(v) => setTheme(v as 'light' | 'dark' | 'system')}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="light">Light</SelectItem>
                        <SelectItem value="dark">Dark</SelectItem>
                        <SelectItem value="system">System</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Language</label>
                    <Select defaultValue="en">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="en">English</SelectItem>
                        <SelectItem value="hi">Hindi</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Currency Format</label>
                  <Input defaultValue="₹" disabled className="max-w-[100px]" />
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="salon">
          <div className="grid gap-6 mt-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-base flex items-center gap-2">
                  <Store className="h-4 w-4" />
                  Salon Details
                </CardTitle>
                <CardDescription>Your business information</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Salon Name</label>
                    <Input
                      value={salonSettings['salon_name'] || ''}
                      onChange={(e) => updateSetting('salon_name', e.target.value)}
                      placeholder="Enter salon name"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Phone</label>
                    <Input
                      value={salonSettings['salon_phone'] || ''}
                      onChange={(e) => updateSetting('salon_phone', e.target.value)}
                      placeholder="Phone number"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Email</label>
                    <Input
                      value={salonSettings['salon_email'] || ''}
                      onChange={(e) => updateSetting('salon_email', e.target.value)}
                      placeholder="Email"
                      type="email"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">GST Number</label>
                    <Input
                      value={salonSettings['gst_number'] || ''}
                      onChange={(e) => updateSetting('gst_number', e.target.value)}
                      placeholder="GSTIN"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Address</label>
                  <Input
                    value={salonSettings['salon_address'] || ''}
                    onChange={(e) => updateSetting('salon_address', e.target.value)}
                    placeholder="Full address"
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-base">Invoice Settings</CardTitle>
                <CardDescription>Configure invoice defaults</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Invoice Prefix</label>
                    <Input
                      value={salonSettings['invoice_prefix'] || 'INV'}
                      onChange={(e) => updateSetting('invoice_prefix', e.target.value)}
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Tax Rate (%)</label>
                    <Input
                      value={salonSettings['tax_rate'] || '18'}
                      onChange={(e) => updateSetting('tax_rate', e.target.value)}
                      type="number"
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="notifications">
          <div className="grid gap-6 mt-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-base flex items-center gap-2">
                  <Bell className="h-4 w-4" />
                  Notification Preferences
                </CardTitle>
                <CardDescription>Control when you receive notifications</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between py-2">
                  <div>
                    <p className="text-sm font-medium">Birthday Reminders</p>
                    <p className="text-xs text-muted-foreground">Get notified about customer birthdays</p>
                  </div>
                  <Badge variant="default">Enabled</Badge>
                </div>
                <div className="flex items-center justify-between py-2">
                  <div>
                    <p className="text-sm font-medium">Low Stock Alerts</p>
                    <p className="text-xs text-muted-foreground">Alert when product stock is low</p>
                  </div>
                  <Badge variant="default">Enabled</Badge>
                </div>
                <div className="flex items-center justify-between py-2">
                  <div>
                    <p className="text-sm font-medium">Salary Due Reminders</p>
                    <p className="text-xs text-muted-foreground">Reminder before salary processing date</p>
                  </div>
                  <Badge variant="secondary">Disabled</Badge>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="system">
          <div className="grid gap-6 mt-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-base flex items-center gap-2">
                  <Database className="h-4 w-4" />
                  System Information
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Application Version</span>
                    <span className="text-sm font-mono">{health?.version || 'v0.1.0'}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Backend Status</span>
                    <Badge variant={health?.status === 'healthy' ? 'default' : 'destructive'}>
                      {health?.status || 'unknown'}
                    </Badge>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Database</span>
                    <span className="text-sm font-mono">{health?.database || 'SQLite'}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Uptime</span>
                    <span className="text-sm">{health?.uptime || '-'}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-base flex items-center gap-2">
                  <Shield className="h-4 w-4" />
                  Security
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium">Auto-Lock Timeout</p>
                    <p className="text-xs text-muted-foreground">Lock the app after inactivity</p>
                  </div>
                  <Select defaultValue="15">
                    <SelectTrigger className="w-[120px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="5">5 minutes</SelectItem>
                      <SelectItem value="15">15 minutes</SelectItem>
                      <SelectItem value="30">30 minutes</SelectItem>
                      <SelectItem value="never">Never</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
