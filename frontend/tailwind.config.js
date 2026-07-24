/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: '#6E6C6F',       // Smoked Pearl
          hover: '#31302F',         // Primary Hover (Off Black)
          emphasis: '#31302F',      // Off Black
          bg: '#F1F0F3',           // Background (updated for luxury palette)
          surface: '#FFFFFF',      // Surface (pure white for premium feel)
          secondary: '#E3DEDE',    // Secondary
          border: '#CACDDC',       // Border
          accent: '#A5A4A8',       // Muted Text
          'text-primary': '#31302F',
          'text-secondary': '#6E6C6F',
          'text-muted': '#A5A4A8',
          'text-light': '#F1F0F3',
        },
      },
      fontFamily: {
        serif: ['Merriweather', 'Georgia', 'serif'],
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'sans-serif'],
      },
      borderRadius: {
        'xl': '18px',
        '2xl': '20px',
      },
      boxShadow: {
        'soft': '0 2px 8px rgba(0, 0, 0, 0.04)',
        'soft-lg': '0 4px 16px rgba(0, 0, 0, 0.08)',
        'soft-xl': '0 8px 24px rgba(0, 0, 0, 0.06)',
        'glow': '0 0 20px rgba(110, 108, 111, 0.08)',
      },
      backgroundImage: {
        'luxury-gradient': 'linear-gradient(135deg, #F1F0F3 0%, #FFFFFF 45%, #E3DEDE 100%)',
        'warm-gradient': 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
        'card-glow': 'radial-gradient(circle at top right, rgba(241, 240, 243, 0.8), rgba(255, 255, 255, 0.2))',
      },
      spacing: {
        'sidebar-left': '320px',
        'sidebar-right': '360px',
      },
    },
  },
  plugins: [],
};
