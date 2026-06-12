/**
 * SalonFlow Motion System
 * ========================
 * Premium animation primitives inspired by Linear, Notion, Arc Browser.
 * All durations tuned for 60fps desktop feel — snappy, not sluggish.
 *
 * Principles:
 * - Micro-interactions: 100-150ms (hover, press, toggle)
 * - State transitions: 200-300ms (open/close, navigate)
 * - Emphasis: 400-600ms (celebrate, draw attention)
 * - Never exceed 500ms for UI feedback
 * - Use spring physics for natural feel
 * - Reduce motion for accessibility
 */

// ─── Duration Tokens ─────────────────────────────────────────────────────────

export const duration = {
  instant: 0.1,
  fast: 0.15,
  normal: 0.2,
  moderate: 0.3,
  slow: 0.4,
  emphasis: 0.6,
} as const

// ─── Easing Curves (CSS & Motion-compatible) ─────────────────────────────────

export const easing = {
  // Default for most UI transitions - slightly bouncy exit
  default: [0.2, 0.0, 0.0, 1.0] as [number, number, number, number],
  // Snappy entrance (for dropdowns, popovers)
  enter: [0.0, 0.0, 0.2, 1.0] as [number, number, number, number],
  // Quick exit (for closing, dismissing)
  exit: [0.4, 0.0, 1.0, 1.0] as [number, number, number, number],
  // Bounce for celebrations/emphasis
  bounce: [0.34, 1.56, 0.64, 1.0] as [number, number, number, number],
  // Linear for progress bars
  linear: [0, 0, 1, 1] as [number, number, number, number],
} as const

// ─── Spring Presets ──────────────────────────────────────────────────────────

export const spring = {
  // Snappy, no overshoot — best for most UI
  snappy: { type: 'spring' as const, stiffness: 500, damping: 30, mass: 1 },
  // Gentle bounce — for popovers, cards
  gentle: { type: 'spring' as const, stiffness: 300, damping: 25, mass: 1 },
  // Bouncy — for success states, celebrations
  bouncy: { type: 'spring' as const, stiffness: 400, damping: 15, mass: 1 },
  // Stiff — for sidebar collapse, instant feedback
  stiff: { type: 'spring' as const, stiffness: 700, damping: 35, mass: 1 },
} as const

// ─── Transition Presets ──────────────────────────────────────────────────────

export const transition = {
  fast: { duration: duration.fast, ease: easing.default },
  normal: { duration: duration.normal, ease: easing.default },
  enter: { duration: duration.normal, ease: easing.enter },
  exit: { duration: duration.fast, ease: easing.exit },
  spring: spring.snappy,
  springGentle: spring.gentle,
} as const

// ─── Animation Variants (Framer Motion) ──────────────────────────────────────

/** Fade in/out */
export const fadeVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: transition.enter },
  exit: { opacity: 0, transition: transition.exit },
}

/** Slide up and fade (for page content, cards) */
export const slideUpVariants = {
  hidden: { opacity: 0, y: 8 },
  visible: { opacity: 1, y: 0, transition: transition.enter },
  exit: { opacity: 0, y: -4, transition: transition.exit },
}

/** Slide down (for dropdowns, menus) */
export const slideDownVariants = {
  hidden: { opacity: 0, y: -4, scale: 0.98 },
  visible: { opacity: 1, y: 0, scale: 1, transition: { ...transition.enter, duration: duration.fast } },
  exit: { opacity: 0, y: -4, scale: 0.98, transition: transition.exit },
}

/** Scale in (for dialogs, modals) */
export const scaleVariants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: { opacity: 1, scale: 1, transition: spring.gentle },
  exit: { opacity: 0, scale: 0.97, transition: transition.exit },
}

/** Stagger children container */
export const staggerContainer = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.04,
      delayChildren: 0.02,
    },
  },
}

/** Stagger child item */
export const staggerItem = {
  hidden: { opacity: 0, y: 6 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: duration.normal, ease: easing.enter },
  },
}

/** KPI card entrance (with counter-like effect) */
export const kpiVariants = {
  hidden: { opacity: 0, y: 12, scale: 0.96 },
  visible: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: spring.gentle,
  },
}

/** Sidebar collapse/expand */
export const sidebarVariants = {
  expanded: { width: 240, transition: spring.stiff },
  collapsed: { width: 64, transition: spring.stiff },
}

/** List item (for nav items, table rows) */
export const listItemVariants = {
  hidden: { opacity: 0, x: -8 },
  visible: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 8 },
}

// ─── Hover/Tap Micro-interactions ────────────────────────────────────────────

export const hoverScale = {
  whileHover: { scale: 1.02 },
  whileTap: { scale: 0.98 },
  transition: { duration: duration.instant },
}

export const hoverLift = {
  whileHover: { y: -2, boxShadow: '0 4px 12px rgba(0,0,0,0.08)' },
  transition: { duration: duration.fast },
}

export const pressScale = {
  whileTap: { scale: 0.97 },
  transition: { duration: duration.instant },
}

// ─── Page Transition ─────────────────────────────────────────────────────────

export const pageTransition = {
  initial: { opacity: 0, y: 6 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -4 },
  transition: { duration: duration.normal, ease: easing.enter },
}

// ─── Reduced Motion ──────────────────────────────────────────────────────────

export const reducedMotionTransition = {
  duration: 0,
  ease: easing.linear,
}
