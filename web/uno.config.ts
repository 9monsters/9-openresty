import { defineConfig, presetAttributify, presetIcons, presetTypography, presetUno, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: [
    ['flex-center', 'flex items-center justify-center'],
    ['flex-between', 'flex items-center justify-between'],
    ['card-hover', 'hover:shadow-md transition-shadow duration-200'],
  ],
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      collections: {
        tabler: () => import('@iconify-json/tabler/icons.json').then(i => i.default),
      },
      extraProperties: {
        display: 'inline-block',
        height: '1.2em',
        width: '1.2em',
        'vertical-align': 'text-bottom',
      },
    }),
    presetTypography(),
  ],
  transformers: [transformerDirectives(), transformerVariantGroup()],
})
