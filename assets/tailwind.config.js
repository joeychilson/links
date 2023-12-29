/** @type {import('tailwindcss').Config} */
export default {
    content: ["./components/**/*.templ", "./pages/**/*.templ", 'node_modules/preline/dist/*.js'],
    theme: {
        extend: {},
    },
    plugins: [
        require('preline/plugin'),
    ],
}