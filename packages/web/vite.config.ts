import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { AntdvNextResolver } from '@antdv-next/auto-import-resolver'
import { defineConfig } from 'vite'
import { resolve } from 'node:path'

// Build mode: STANDALONE=true builds to packages/web/dist for independent deployment.
// Default builds directly into the Go embed directory.
const isStandalone = process.env.STANDALONE === 'true'
const embedOutDir = resolve(__dirname, '../../internal/embedweb/dist')
const standaloneOutDir = resolve(__dirname, 'dist')

export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  plugins: [
    vue(),
    UnoCSS(),
    Components({
      resolvers: [AntdvNextResolver()],
      directoryAsNamespace: true,
    }),
    AutoImport({
      imports: [
        'vue',
        'vue-router',
        'pinia',
        '@vueuse/core',
        { 'antdv-next': ['message', 'notification', 'Modal'] },
      ],
      vueTemplate: true,
      dts: 'src/auto-imports.d.ts',
    }),
  ],
  base: isStandalone ? '/' : undefined,
  build: {
    outDir: isStandalone ? standaloneOutDir : embedOutDir,
    emptyOutDir: true,
    rollupOptions: {
      output: {
        // Go's embed package excludes files starting with _ or .
        // Vue plugin generates helper chunks like _plugin-vue_export-helper-*.js,
        // so strip the leading underscore from chunk names.
        chunkFileNames(chunkInfo) {
          const name = chunkInfo.name.replace(/^_/, '')
          return `assets/${name}-[hash].js`
        },
      },
    },
  },
  server: {
    port: 3001,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
})
