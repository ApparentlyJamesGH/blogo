/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./*.{html,templ}", "./templates/*.{html,templ}"],
    theme: {
        extend: {
            fontFamily: {},
        },
    },
    plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
};
