/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./*.{html,templ}", "./templates/**/*.{html,templ}"],
    theme: {
        extend: {
            fontFamily: {},
            colors: {
                'background': "var(--blogo-background)",
                'primary': "var(--blogo-primary)",
                'primary-emphasis': "var(--blogo-primary-emphasis)",
                'primary-dark': "var(--blogo-primary-dark)",
                'primary-light': "var(--blogo-primary-light)",
                'blogo-text': "var(--blogo-text)",
                'blogo-text-fade': "var(--blogo-text-fade)",
                'blogo-text-select': "var(--blogo-text-select)",
            },
        },
    },
    plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
};
