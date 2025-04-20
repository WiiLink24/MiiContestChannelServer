const { addDynamicIconSelectors } = require('@iconify/tailwind');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.html"],
  theme: {
    fontFamily: {
      sans: ["Inter", "sans-serif"],
      serif: ["Rubik", "sans-serif"],
    },
    extend: {},
  },
  plugins: [
    addDynamicIconSelectors()
  ],
};

