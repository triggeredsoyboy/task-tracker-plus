/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js}", "node_modules/preline/dist/*.js"],
  theme: {
    container: {
      center: true,
      padding: "24px",
    },
    fontFamily: {
      playfair: ["Playfair Display", "serif"],
      poppins: ["Poppins", "sans-serif"],
    },
    extend: {
      colors: {
        primary: "#3b82f6",
        secondary: "#60a5fa",
        tertiary: "rgb(17 24 39 / 0.1)", // gray-900/10
        dark: "#0f172a", // slate-900
        light: "#64748b", // slate-500
      },
      screens: {
        sm: "640px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
        "2xl": "1360px",
      },
    },
  },
  plugins: [require("@tailwindcss/forms"), require("preline/plugin")],
};
