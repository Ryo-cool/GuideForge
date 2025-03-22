/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        background: '#e6e7ee',
        'text-primary': '#44476a',
        'text-secondary': '#31344b',
        'text-muted': '#7e8299',
        primary: '#5e72e4',
        'primary-hover': '#4e62d2',
        secondary: '#8392ab',
        success: '#2dce89',
        info: '#11cdef',
        warning: '#fb6340',
        danger: '#f5365c',
        light: '#f0f3fc',
        dark: '#273444',
      },
      boxShadow: {
        'soft': '6px 6px 12px #d1d9e6, -6px -6px 12px #ffffff',
        'inset': 'inset 2px 2px 5px #d1d9e6, inset -2px -2px 5px #ffffff',
        'focus': '0 0 0 0.2rem rgba(94, 114, 228, 0.25)',
      },
      borderRadius: {
        DEFAULT: '0.75rem',
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
    },
  },
  plugins: [],
} 