import { Outlet, useLocation } from 'react-router-dom'
import { motion, AnimatePresence } from 'motion/react'
import { Sidebar } from '@/components/layout/Sidebar'
import { Header } from '@/components/layout/Header'
import { ErrorBoundary } from '@/components/shared/ErrorBoundary'
import { useNavigationShortcuts } from '@/lib/desktop-ux'

export function MainLayout() {
  const location = useLocation()
  useNavigationShortcuts()

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header />
        <main className="relative flex-1 overflow-y-auto scrollbar-hidden">
          {/* Subtle ambient gradient */}
          <div className="pointer-events-none absolute inset-0 gradient-mesh opacity-60" />
          <AnimatePresence mode="wait">
            <motion.div
              key={location.pathname}
              initial={{ opacity: 0, y: 6, filter: 'blur(2px)' }}
              animate={{ opacity: 1, y: 0, filter: 'blur(0px)' }}
              exit={{ opacity: 0, y: -3, filter: 'blur(1px)' }}
              transition={{ duration: 0.18, ease: [0.2, 0, 0, 1] }}
              className="relative p-6"
            >
              <ErrorBoundary key={location.pathname}>
                <Outlet />
              </ErrorBoundary>
            </motion.div>
          </AnimatePresence>
        </main>
      </div>
    </div>
  )
}
