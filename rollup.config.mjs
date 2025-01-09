import { defineConfig } from 'rollup';
import esbuild from 'rollup-plugin-esbuild';

export default defineConfig({
  input: './src/main.ts',
  output: {
    file: 'bin/bin.js',
    format: 'cjs'
  },
  plugins: [esbuild()]
});
