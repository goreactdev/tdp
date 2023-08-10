import react from '@vitejs/plugin-react'
import reactRefresh from '@vitejs/plugin-react-refresh'
import { babel } from '@rollup/plugin-babel';
import replace from '@rollup/plugin-replace';
import { defineConfig } from 'vite'


// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    
    esbuildOptions: {
      target: 'es2020',
    },
  },
  esbuild: {
    // https://github.com/vitejs/vite/issues/8644#issuecomment-1159308803
    logOverride: { 'this-is-undefined-in-esm': 'silent' },
  },
  plugins: [
    react({
      babel: {
        plugins: [
          "babel-plugin-twin",
          "babel-plugin-macros",
          "babel-plugin-styled-components",
        ],
        ignore: ["\x00commonjsHelpers.js"], // Weird babel-macro bug fix
      },
    }),
  ],
  server: {
    host: true,
    port: 3000,
  },
})
