// tailwind.config.js
const { fontFamily } = require('tailwindcss/defaultTheme')

module.exports = {
  content: [
    './src/app/**/*.{js,ts,jsx,tsx}',
    './src/components/**/*.{js,ts,jsx,tsx}',
    './src/lib/**/*.{js,ts,jsx,tsx}',
  ],
  darkMode: 'media', // or 'class' if you want manual dark mode
  theme: {
    extend: {
      colors: {
        background: 'rgb(var(--background))',
        foreground: 'rgb(var(--foreground))',
        primary: {
          DEFAULT: '#6366f1',
          light: '#818cf8',
          dark: '#4f46e5',
        },
        secondary: {
          DEFAULT: '#10b981',
          light: '#34d399',
          dark: '#059669',
        },
        message: {
          sent: '#6366f1',
          received: '#3f3f46',
        }
      },
      fontFamily: {
        sans: ['var(--font-sans)', ...fontFamily.sans],
        mono: ['var(--font-mono)', ...fontFamily.mono],
      },
      animation: {
        'message-in': 'messageIn 0.3s ease-out',
        'typing-indicator': 'pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        messageIn: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        }
      },
      boxShadow: {
        chat: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
        message: '0 2px 4px -1px rgba(0, 0, 0, 0.1)',
      }
    },
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('tailwindcss-animate'),
    // Custom plugin for chat-specific utilities
    function({ addUtilities }) {
      const newUtilities = {
        '.chat-container': {
          height: 'calc(100vh - 160px)',
          '@screen md': {
            height: 'calc(100vh - 128px)',
          },
        },
        '.message-bubble': {
          borderRadius: '1.125rem',
          maxWidth: '85%',
          wordBreak: 'break-word',
        },
        '.typing-indicator': {
          width: '48px',
          height: '24px',
          position: 'relative',
        },
      }
      addUtilities(newUtilities)
    }
  ],
}
