/** @type {import('tailwindcss').Config} */

module.exports = {
  content: [
    "./app/**/*.{html,js,jsx}",
    "./components/**/*.{html,js,jsx}",
    "./sections/**/*.{html,js,jsx}",
    "./styles/**/*.{js,jsx}",
  ],
  mode: "jit",
  theme: {
    extend: {
      colors: {
        "primary-black": "#000000",
        "secondary-gray": "#18181B",
        "primary-pink": "#E32CF7",
        "primary-blue": "#0070EF",
        "secondary-blue": "2D89F2",
        "primary-green": "#3ED976",
      },
      transitionTimingFunction: {
        "out-flex": "cubic-bezier(0.05, 0.6, 0.4, 0.9)",
      },
    },
  },

  plugins: [],
};
