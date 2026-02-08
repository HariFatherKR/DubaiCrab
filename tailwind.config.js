/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// CoachPro style colors (Light Theme)
				primary: {
					DEFAULT: '#14b8a6',
					hover: '#0d9488',
					light: '#5eead4',
					50: '#f0fdfa',
					100: '#ccfbf1',
					500: '#14b8a6',
					600: '#0d9488',
					700: '#0f766e'
				},
				secondary: '#06b6d4',
				accent: '#10b981',
				
				// Light theme backgrounds
				bg: {
					DEFAULT: '#e6f5f3',
					light: '#f0fdfa',
					card: 'rgba(255, 255, 255, 0.6)'
				},
				
				// Glass morphism
				glass: {
					bg: 'rgba(255, 255, 255, 0.6)',
					border: 'rgba(255, 255, 255, 0.8)',
					hover: 'rgba(255, 255, 255, 0.8)'
				},
				
				// Text colors (dark for light theme)
				text: {
					DEFAULT: '#1e293b',
					secondary: '#64748b',
					muted: '#94a3b8'
				},
				
				// Status colors
				success: '#059669',
				warning: '#d97706',
				error: '#dc2626',
				info: '#2563eb'
			},
			borderRadius: {
				'2xl': '1rem',
				'3xl': '1.5rem',
				'4xl': '2rem'
			},
			boxShadow: {
				'glass': '0 8px 32px rgba(0, 0, 0, 0.08)',
				'glass-lg': '0 16px 48px rgba(0, 0, 0, 0.1)',
				'glass-hover': '0 12px 40px rgba(0, 0, 0, 0.12)',
				'glow': '0 0 20px rgba(20, 184, 166, 0.2)',
				'glow-lg': '0 0 40px rgba(20, 184, 166, 0.3)',
				'soft': '0 2px 8px rgba(0, 0, 0, 0.04)'
			},
			backdropBlur: {
				'xs': '2px',
				'3xl': '64px'
			},
			fontFamily: {
				sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif']
			},
			backgroundImage: {
				'gradient-mint': 'linear-gradient(135deg, #e6f5f3 0%, #d9ece9 25%, #cce4e7 50%, #c5e1ec 75%, #d4eaf5 100%)',
				'gradient-teal': 'linear-gradient(135deg, #14b8a6, #06b6d4)'
			}
		}
	},
	plugins: []
};
