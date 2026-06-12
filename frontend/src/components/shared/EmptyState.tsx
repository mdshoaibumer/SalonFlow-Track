import { motion } from 'motion/react'
import type { LucideIcon } from 'lucide-react'
import { FileX } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description: string
  action?: { label: string; onClick: () => void }
}

export function EmptyState({ icon: Icon = FileX, title, description, action }: EmptyStateProps) {
  return (
    <motion.div
      className="flex flex-col items-center justify-center py-16 text-center"
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: [0.43, 0.13, 0.23, 0.96] }}
    >
      <motion.div
        className="relative rounded-2xl bg-gradient-to-br from-muted/80 to-muted/40 p-5 mb-5"
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.1, type: 'spring', stiffness: 200, damping: 15 }}
      >
        {/* Decorative rings */}
        <div className="absolute inset-0 rounded-2xl border border-border/30 animate-pulse" />
        <div className="absolute -inset-2 rounded-3xl border border-border/10" />
        <Icon className="h-8 w-8 text-muted-foreground/70" />
      </motion.div>
      <motion.h3
        className="text-lg font-semibold"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2 }}
      >
        {title}
      </motion.h3>
      <motion.p
        className="text-sm text-muted-foreground mt-1.5 max-w-sm"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
      >
        {description}
      </motion.p>
      {action && (
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          <Button onClick={action.onClick} className="mt-5 rounded-xl" size="sm">
            {action.label}
          </Button>
        </motion.div>
      )}
    </motion.div>
  )
}
