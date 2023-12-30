/** @type {import('tailwindcss').Config} */
export default {
    content: ["./components/**/*.templ", "./layouts/**/*.templ", "./pages/**/*.templ", "node_modules/preline/dist/*.js"],
    safelist: [
        {
            pattern: /ml-+/,
        },
    ],
    theme: {
        extend: {},
    },
    plugins: [require("@tailwindcss/forms"), require("preline/plugin")],
};
