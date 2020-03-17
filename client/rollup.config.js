import typescript from 'rollup-plugin-typescript';
import { terser } from "rollup-plugin-terser";

export default {
    input: "./script/main.ts",
    plugins: [
        typescript(),
        terser(),
    ],
    output: {
        file: "../build/production/main.js",
        format: "iife",
        sourcemap: true,
    }
};