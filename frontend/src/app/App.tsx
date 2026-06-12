import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from 'sonner'
import { ThemeProvider } from './providers/ThemeProvider'
import { AuthProvider } from './providers/AuthProvider'
import { AppRouter } from './router/AppRouter'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="light" storageKey="salonflow-theme">
        <BrowserRouter>
          <AuthProvider>
            <AppRouter />
          </AuthProvider>
          <Toaster
            position="bottom-right"
            toastOptions={{
              className: 'surface-raised text-[13px] border-border/60 shadow-elevation-3',
              duration: 3500,
            }}
            gap={8}
            offset={16}
            richColors
            closeButton
          />
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}
