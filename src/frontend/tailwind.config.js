/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./*.{html,templ}", "./templates/*.{html,templ}"],
    theme: {
        extend: {
            fontFamily: {},
        },
        colors: {
            'background': "var(--blogo-background)",
            'primary': "var(--blogo-primary)",
            'primary-emphasis': "var(--blogo-primary-emphasis)",
            'secondary': "var(--blogo-secondary)",
            'secondary-emphasis': "var(--blogo-secondary-emphasis)",
            'blogo-text': "var(--blogo-text)",
            'blogo-text-fade': "var(--blogo-text-fade)",
            'blogo-title': "var(--blogo-text-title)",
            'blogo-bold': "var(--blogo-text-bold)",
        },
    },
    plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
};
