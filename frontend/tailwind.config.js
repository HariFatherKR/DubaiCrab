/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './index.html'],
  theme: {
    extend: {
      colors: {
        // Dubai Crab dark chocolate theme
        'crab-dark': '#3E2723',
        'crab-medium': '#4E342E',
        'crab-light': '#5D4037',
        'crab-accent': '#8D6E63',
        'crab-text': '#EFEBE9',
        'crab-muted': '#BCAAA4',
        'crab-orange': '#FF8A65',
        'crab-gold': '#FFD54F'
      }
    }
  },
  plugins: []
};
