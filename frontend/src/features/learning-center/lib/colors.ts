// AI Finance Tutor Color Palette
// Inspired by premium fintech platforms: Apple, Linear, Stripe, Vercel

export const colors = {
  // Primary brand colors
  primary: '#6E6C6F', // Smoked Pearl
  primaryHover: '#31302F', // Primary Hover for interactive states
  emphasis: '#31302F', // Off Black
  
  // Background colors (updated for soft neutral luxury)
  background: '#F1F0F3', // Light background (premium base)
  surface: '#FFFFFF', // Pure white surface (luxury feel)
  secondary: '#E3DEDE', // Secondary surface
  
  // Border and accent colors
  border: '#CACDDC', // Pale Shale
  accent: '#A5A4A8', // Mysterious Mauve (muted text/icons)
  
  // Semantic colors
  text: {
    primary: '#31302F', // Off Black
    secondary: '#6E6C6F', // Smoked Pearl
    muted: '#A5A4A8', // Mysterious Mauve
    light: '#F1F0F3', // Aragonite White
  },
  
  // State colors (subtle)
  success: '#7B9E7A', // Muted green
  warning: '#B8A66B', // Muted gold
  error: '#A87C7C', // Muted red
  info: '#7A9BB8', // Muted blue
};

// Gradient system for premium aesthetic
export const gradients = {
  luxury: 'linear-gradient(135deg, #F1F0F3 0%, #FFFFFF 45%, #E3DEDE 100%)',
  cardGlow: 'radial-gradient(circle at top right, rgba(241, 240, 243, 0.8), rgba(255, 255, 255, 0.2))',
  subtle: 'linear-gradient(180deg, rgba(255, 255, 255, 0.6) 0%, rgba(241, 240, 243, 0.3) 100%)',
};

// Shadow system for soft ambient lighting
export const shadows = {
  soft: '0 2px 8px rgba(0, 0, 0, 0.04)',
  'soft-lg': '0 4px 16px rgba(0, 0, 0, 0.08)',
  'soft-xl': '0 8px 24px rgba(0, 0, 0, 0.06)',
  glow: '0 0 20px rgba(110, 108, 111, 0.08)',
};

// Tailwind color mappings for custom config
export const tailwindColors = {
  'brand-primary': colors.primary,
  'brand-emphasis': colors.emphasis,
  'brand-bg': colors.background,
  'brand-surface': colors.surface,
  'brand-secondary': colors.secondary,
  'brand-border': colors.border,
  'brand-accent': colors.accent,
  'brand-text-primary': colors.text.primary,
  'brand-text-secondary': colors.text.secondary,
  'brand-text-muted': colors.text.muted,
  'brand-text-light': colors.text.light,
};
