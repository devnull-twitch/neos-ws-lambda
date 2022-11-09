/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.tpl", "assets/**/*.js"],
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      lavender: '#D8DCFF',
      maximumblue: '#AEADF0',
      rosybrown: '#C38D94',
      rosedust: '#A76571',
      raisinblack: '#292938'
    },
    extend: {},
  },
  plugins: [],
}
