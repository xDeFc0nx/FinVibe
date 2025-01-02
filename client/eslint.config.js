import globals from "globals";
import solidPlugin from "eslint-plugin-solid";
import tseslint from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";

export default [
  {
    ignores: ["tailwind.config.cjs", "postcss.config.js"],
  },
  {
    ignores: ["dist", "node_modules"],
  },
  {
    files: ["**/*.{ts,tsx,js,jsx}"],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parser: tsParser,
    },
    plugins: {
      solid: solidPlugin,
      "@typescript-eslint": tseslint,
    },
    rules: {
      ...solidPlugin.configs.recommended.rules,

      "no-unused-vars": ["warn", { argsIgnorePattern: "^_" }],
    },
  },
];
