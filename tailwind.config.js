/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    '**/*.templ',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}

