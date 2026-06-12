/**
 * SalonFlow Motion Components
 * ============================
 * Reusable animated wrappers that apply the motion system consistently.
 * Drop-in replacements for static containers.
 */

import { motion, AnimatePresence, useReducedMotion } from 'motion/react'
import { forwardRef, type ReactNode } from 'react'
import {
  fadeVariants,
  slideUpVariants,
  scaleVariants,
  staggerContainer,
  staggerItem,
  kpiVariants,
  pageTransition,
  hoverLift,
  pressScale,
  spring,
  duration,
  easing,
} from '@/lib/motion'
import { cn } from '@/lib/utils'

// ─── Page Wrapper (wraps every page for enter/exit) ──────────────────────────

interface MotionPageProps {
  children: ReactNode
  className?: string
}

export const MotionPage = forwardRef<HTMLDivElement, MotionPageProps>(
  ({ children, className }, ref) => {
    const shouldReduceMotion = useReducedMotion()

    return (
      <motion.div
        ref={ref}
        initial={shouldReduceMotion ? false : pageTransition.initial}
        animate={pageTransition.animate}
        transition={pageTransition.transition}
        className={cn('space-y-6', className)}
      >
        {children}
      </motion.div>
    )
  }
)
MotionPage.displayName = 'MotionPage'

// ─── Staggered List (for KPI cards, grid items) ──────────────────────────────

interface MotionStaggerProps {
  children: ReactNode
  className?: string
  delay?: number
}

export function MotionStagger({ children, className, delay = 0 }: MotionStaggerProps) {
  const shouldReduceMotion = useReducedMotion()

  if (shouldReduceMotion) {
    return <div className={className}>{children}</div>
  }

  return (
    <motion.div
      variants={staggerContainer}
      initial="hidden"
      animate="visible"
      transition={{ delayChildren: delay }}
      className={className}
    >
      {children}
    </motion.div>
  )
}

// ─── Stagger Item ────────────────────────────────────────────────────────────

interface MotionItemProps {
  children: ReactNode
  className?: string
}

export function MotionItem({ children, className }: MotionItemProps) {
  return (
    <motion.div variants={staggerItem} className={className}>
      {children}
    </motion.div>
  )
}

// ─── KPI Card Wrapper ────────────────────────────────────────────────────────

export function MotionKPI({ children, className }: MotionItemProps) {
  return (
    <motion.div
      variants={kpiVariants}
      className={className}
      whileHover={hoverLift.whileHover}
    >
      {children}
    </motion.div>
  )
}

// ─── Fade In ─────────────────────────────────────────────────────────────────

interface MotionFadeProps {
  children: ReactNode
  show?: boolean
  className?: string
}

export function MotionFade({ children, show = true, className }: MotionFadeProps) {
  return (
    <AnimatePresence mode="wait">
      {show && (
        <motion.div
          variants={fadeVariants}
          initial="hidden"
          animate="visible"
          exit="exit"
          className={className}
        >
          {children}
        </motion.div>
      )}
    </AnimatePresence>
  )
}

// ─── Slide Up (for content sections) ─────────────────────────────────────────

export function MotionSlideUp({ children, className }: MotionItemProps) {
  return (
    <motion.div
      variants={slideUpVariants}
      initial="hidden"
      animate="visible"
      className={className}
    >
      {children}
    </motion.div>
  )
}

// ─── Scale In (for dialogs, popovers) ────────────────────────────────────────

export function MotionScale({ children, show = true, className }: MotionFadeProps) {
  return (
    <AnimatePresence mode="wait">
      {show && (
        <motion.div
          variants={scaleVariants}
          initial="hidden"
          animate="visible"
          exit="exit"
          className={className}
        >
          {children}
        </motion.div>
      )}
    </AnimatePresence>
  )
}

// ─── Pressable (button-like press feedback) ──────────────────────────────────

interface MotionPressableProps {
  children: ReactNode
  className?: string
}

export function MotionPressable({ children, className }: MotionPressableProps) {
  return (
    <motion.div
      whileTap={pressScale.whileTap}
      className={cn('cursor-pointer', className)}
    >
      {children}
    </motion.div>
  )
}

// ─── Presence Wrapper (for conditional content) ──────────────────────────────

interface MotionPresenceProps {
  children: ReactNode
  show: boolean
  mode?: 'wait' | 'sync' | 'popLayout'
}

export function MotionPresence({ children, show, mode = 'wait' }: MotionPresenceProps) {
  return (
    <AnimatePresence mode={mode}>
      {show && (
        <motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          exit={{ opacity: 0, height: 0 }}
          transition={{ duration: duration.normal, ease: easing.default }}
          style={{ overflow: 'hidden' }}
        >
          {children}
        </motion.div>
      )}
    </AnimatePresence>
  )
}

// ─── Re-exports ──────────────────────────────────────────────────────────────

export { motion, AnimatePresence, useReducedMotion }
export { spring, duration, easing }
